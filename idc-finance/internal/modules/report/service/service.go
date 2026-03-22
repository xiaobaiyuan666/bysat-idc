package service

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	customerDomain "idc-finance/internal/modules/customer/domain"
	customerService "idc-finance/internal/modules/customer/service"
	orderDomain "idc-finance/internal/modules/order/domain"
	orderService "idc-finance/internal/modules/order/service"
	reportDTO "idc-finance/internal/modules/report/dto"
	"idc-finance/internal/platform/audit"
)

type Service struct {
	db        *sql.DB
	customers *customerService.Service
	orders    *orderService.Service
	audit     *audit.Service
}

func New(db *sql.DB, customerSvc *customerService.Service, orderSvc *orderService.Service, auditSvc *audit.Service) *Service {
	return &Service{
		db:        db,
		customers: customerSvc,
		orders:    orderSvc,
		audit:     auditSvc,
	}
}

func (service *Service) GetWorkbench() reportDTO.WorkbenchResponse {
	if service.db != nil {
		return service.getWorkbenchFromMySQL()
	}
	return service.getWorkbenchFromMemory()
}

func (service *Service) GetReportOverview() reportDTO.ReportOverviewResponse {
	if service.db != nil {
		return service.getReportOverviewFromMySQL()
	}
	return service.getReportOverviewFromMemory()
}

func (service *Service) getWorkbenchFromMySQL() reportDTO.WorkbenchResponse {
	customerTotal := service.queryCount(`SELECT COUNT(*) FROM customers`)
	activeServiceTotal := service.queryCount(`SELECT COUNT(*) FROM services WHERE status = 'ACTIVE'`)
	unpaidInvoiceTotal := service.queryCount(`SELECT COUNT(*) FROM invoices WHERE status = 'UNPAID'`)
	pendingIdentityTotal := service.queryCount(`SELECT COUNT(*) FROM customer_identities WHERE verify_status = 'PENDING'`)
	openTicketTotal := service.queryCount(`SELECT COUNT(*) FROM tickets WHERE status IN ('OPEN', 'PROCESSING')`)
	todayIncome := service.queryFloat(`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE DATE(paid_at) = CURDATE()`)
	monthIncome := service.queryFloat(`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE DATE_FORMAT(paid_at, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')`)
	overdueInvoiceTotal := service.queryCount(`SELECT COUNT(*) FROM invoices WHERE status = 'UNPAID' AND due_at IS NOT NULL AND due_at < NOW()`)
	expiringServiceTotal := service.queryCount(`SELECT COUNT(*) FROM services WHERE status IN ('ACTIVE', 'SUSPENDED') AND next_due_at IS NOT NULL AND next_due_at BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 7 DAY)`)
	suspendedServiceTotal := service.queryCount(`SELECT COUNT(*) FROM services WHERE status = 'SUSPENDED'`)

	return reportDTO.WorkbenchResponse{
		SummaryCards: []reportDTO.MetricCard{
			{Key: "customers", Label: "客户总数", Value: fmt.Sprintf("%d", customerTotal), Hint: "当前系统客户档案", Tone: "primary"},
			{Key: "services", Label: "运行中服务", Value: fmt.Sprintf("%d", activeServiceTotal), Hint: "状态为 ACTIVE 的实例", Tone: "success"},
			{Key: "unpaidInvoices", Label: "待支付账单", Value: fmt.Sprintf("%d", unpaidInvoiceTotal), Hint: "待催收与待跟进账单", Tone: "warning"},
			{Key: "pendingIdentity", Label: "待审核实名", Value: fmt.Sprintf("%d", pendingIdentityTotal), Hint: "需要后台处理的实名申请", Tone: "warning"},
			{Key: "openTickets", Label: "处理中工单", Value: fmt.Sprintf("%d", openTicketTotal), Hint: "OPEN / PROCESSING 工单", Tone: "danger"},
			{Key: "todayIncome", Label: "今日收款", Value: formatCurrency(todayIncome), Hint: fmt.Sprintf("本月累计 %s", formatCurrency(monthIncome)), Tone: "primary"},
		},
		RiskCards: []reportDTO.MetricCard{
			{Key: "overdueInvoices", Label: "逾期账单", Value: fmt.Sprintf("%d", overdueInvoiceTotal), Hint: "需要催收或财务跟进", Tone: "danger"},
			{Key: "expiringServices", Label: "7 日内到期服务", Value: fmt.Sprintf("%d", expiringServiceTotal), Hint: "建议生成续费待办", Tone: "warning"},
			{Key: "suspendedServices", Label: "已暂停服务", Value: fmt.Sprintf("%d", suspendedServiceTotal), Hint: "核对欠费与运维动作", Tone: "danger"},
		},
		RevenueTrend:      service.queryTrend(`SELECT DATE(paid_at) AS point_day, COALESCE(SUM(amount), 0) AS total_amount, COUNT(*) AS total_count FROM payments WHERE paid_at >= ? GROUP BY DATE(paid_at) ORDER BY point_day`, 7),
		ServiceStatus:     service.queryStatusBuckets(`SELECT status, COUNT(*) AS total_count, 0 AS total_amount FROM services GROUP BY status`, serviceLabel),
		PendingIdentities: service.queryPendingIdentities(),
		OverdueInvoices:   service.queryOverdueInvoices(),
		ExpiringServices:  service.queryExpiringServices(),
		OpenTickets:       service.queryOpenTickets(),
		RecentAudits:      service.queryRecentAudits(),
	}
}

func (service *Service) getReportOverviewFromMySQL() reportDTO.ReportOverviewResponse {
	monthIncome := service.queryFloat(`SELECT COALESCE(SUM(amount), 0) FROM payments WHERE DATE_FORMAT(paid_at, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')`)
	monthRefund := service.queryFloat(`SELECT COALESCE(SUM(amount), 0) FROM refunds WHERE DATE_FORMAT(created_at, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')`)
	paidInvoiceCount := service.queryCount(`SELECT COUNT(*) FROM invoices WHERE status = 'PAID'`)
	unpaidAmount := service.queryFloat(`SELECT COALESCE(SUM(total_amount), 0) FROM invoices WHERE status = 'UNPAID'`)
	activeServiceCount := service.queryCount(`SELECT COUNT(*) FROM services WHERE status = 'ACTIVE'`)
	suspendedServiceCount := service.queryCount(`SELECT COUNT(*) FROM services WHERE status = 'SUSPENDED'`)

	return reportDTO.ReportOverviewResponse{
		Headline: []reportDTO.MetricCard{
			{Key: "monthIncome", Label: "本月收款", Value: formatCurrency(monthIncome), Hint: "已完成支付流水", Tone: "primary"},
			{Key: "monthRefund", Label: "本月退款", Value: formatCurrency(monthRefund), Hint: "退款单已完成金额", Tone: "warning"},
			{Key: "netIncome", Label: "本月净收入", Value: formatCurrency(monthIncome - monthRefund), Hint: "收款减去退款", Tone: "success"},
			{Key: "paidInvoices", Label: "已支付账单", Value: fmt.Sprintf("%d", paidInvoiceCount), Hint: "已完成回款单据", Tone: "primary"},
			{Key: "unpaidAmount", Label: "待收金额", Value: formatCurrency(unpaidAmount), Hint: "UNPAID 账单应收", Tone: "danger"},
			{Key: "activeServices", Label: "运行中服务", Value: fmt.Sprintf("%d", activeServiceCount), Hint: fmt.Sprintf("暂停中 %d", suspendedServiceCount), Tone: "success"},
		},
		RevenueTrend:    service.queryTrend(`SELECT DATE(paid_at) AS point_day, COALESCE(SUM(amount), 0) AS total_amount, COUNT(*) AS total_count FROM payments WHERE paid_at >= ? GROUP BY DATE(paid_at) ORDER BY point_day`, 30),
		RefundTrend:     service.queryTrend(`SELECT DATE(created_at) AS point_day, COALESCE(SUM(amount), 0) AS total_amount, COUNT(*) AS total_count FROM refunds WHERE created_at >= ? GROUP BY DATE(created_at) ORDER BY point_day`, 30),
		InvoiceStatus:   service.queryStatusBuckets(`SELECT status, COUNT(*) AS total_count, COALESCE(SUM(total_amount), 0) AS total_amount FROM invoices GROUP BY status`, invoiceLabel),
		ServiceStatus:   service.queryStatusBuckets(`SELECT status, COUNT(*) AS total_count, 0 AS total_amount FROM services GROUP BY status`, serviceLabel),
		PaymentChannels: service.queryStatusBuckets(`SELECT channel, COUNT(*) AS total_count, COALESCE(SUM(amount), 0) AS total_amount FROM payments GROUP BY channel`, channelLabel),
		BillingCycles:   service.queryStatusBuckets(`SELECT billing_cycle, COUNT(*) AS total_count, COALESCE(SUM(amount), 0) AS total_amount FROM orders GROUP BY billing_cycle`, billingCycleLabel),
		CustomerGroups:  service.queryStatusBuckets(`SELECT COALESCE(g.name, '未分组') AS group_name, COUNT(*) AS total_count, 0 AS total_amount FROM customers c LEFT JOIN customer_groups g ON g.id = c.group_id GROUP BY COALESCE(g.name, '未分组')`, passthroughLabel),
		TopProducts:     service.queryRankItems(`SELECT product_name, COUNT(*) AS total_count, COALESCE(SUM(amount), 0) AS total_amount, '' AS extra FROM orders GROUP BY product_name ORDER BY total_amount DESC LIMIT 5`),
		TopReceivables:  service.queryRankItems(`SELECT c.name, COUNT(*) AS total_count, COALESCE(SUM(i.total_amount), 0) AS total_amount, IFNULL(DATE_FORMAT(MAX(i.due_at), '%Y-%m-%d %H:%i:%s'), '') AS extra FROM invoices i INNER JOIN customers c ON c.id = i.customer_id WHERE i.status = 'UNPAID' GROUP BY c.id, c.name ORDER BY total_amount DESC LIMIT 5`),
	}
}

func (service *Service) getWorkbenchFromMemory() reportDTO.WorkbenchResponse {
	customers := service.customers.List()
	invoices, _ := service.orders.ListInvoices(orderDomain.InvoiceListFilter{
		Page:  1,
		Limit: 0,
		Sort:  "created_at",
		Order: "desc",
	})
	services, _ := service.orders.ListServices(orderDomain.ServiceListFilter{
		Page:  1,
		Limit: 0,
		Sort:  "created_at",
		Order: "desc",
	})

	var pendingIdentityTotal, openTicketTotal, unpaidInvoiceTotal, activeServiceTotal, suspendedServiceTotal, overdueInvoiceTotal, expiringServiceTotal int64
	pendingIdentities := make([]reportDTO.PendingIdentityItem, 0)
	openTickets := make([]reportDTO.TicketTodoItem, 0)
	overdueInvoices := make([]reportDTO.OverdueInvoiceItem, 0)
	expiringServices := make([]reportDTO.ExpiringServiceItem, 0)
	now := time.Now()

	for _, customer := range customers {
		if customer.Identity != nil && string(customer.Identity.VerifyStatus) == "PENDING" {
			pendingIdentityTotal++
			pendingIdentities = append(pendingIdentities, reportDTO.PendingIdentityItem{
				CustomerID:   customer.ID,
				CustomerNo:   customer.CustomerNo,
				CustomerName: customer.Name,
				SubjectName:  firstNonEmpty(customer.Identity.SubjectName, customer.Name),
				SubmittedAt:  customer.Identity.ReviewedAt,
			})
		}

		items, ok := service.customers.ListTickets(customer.ID)
		if ok {
			for _, ticket := range items {
				if ticket.Status == "OPEN" || ticket.Status == "PROCESSING" {
					openTicketTotal++
					openTickets = append(openTickets, reportDTO.TicketTodoItem{
						TicketNo:     ticket.No,
						CustomerName: customer.Name,
						Title:        ticket.Name,
						Status:       ticketLabel(ticket.Status),
						UpdatedAt:    ticket.UpdatedAt,
					})
				}
			}
		}
	}

	for _, invoice := range invoices {
		if string(invoice.Status) == "UNPAID" {
			unpaidInvoiceTotal++
			if dueTime, err := parseDateTime(invoice.DueAt); err == nil && dueTime.Before(now) {
				overdueInvoiceTotal++
				overdueInvoices = append(overdueInvoices, reportDTO.OverdueInvoiceItem{
					InvoiceID:    invoice.ID,
					InvoiceNo:    invoice.InvoiceNo,
					CustomerName: lookupCustomerName(customers, invoice.CustomerID),
					Amount:       invoice.TotalAmount,
					DueAt:        invoice.DueAt,
					DaysOverdue:  int64(now.Sub(dueTime).Hours() / 24),
				})
			}
		}
	}

	for _, item := range services {
		switch string(item.Status) {
		case "ACTIVE":
			activeServiceTotal++
		case "SUSPENDED":
			suspendedServiceTotal++
		}
		if dueTime, err := parseDateTime(item.NextDueAt); err == nil {
			daysRemaining := int64(dueTime.Sub(now).Hours() / 24)
			if daysRemaining >= 0 && daysRemaining <= 7 {
				expiringServiceTotal++
				expiringServices = append(expiringServices, reportDTO.ExpiringServiceItem{
					ServiceID:     item.ID,
					ServiceNo:     item.ServiceNo,
					CustomerName:  lookupCustomerName(customers, item.CustomerID),
					ProductName:   item.ProductName,
					Status:        serviceLabel(string(item.Status)),
					NextDueAt:     item.NextDueAt,
					DaysRemaining: daysRemaining,
				})
			}
		}
	}

	return reportDTO.WorkbenchResponse{
		SummaryCards: []reportDTO.MetricCard{
			{Key: "customers", Label: "客户总数", Value: fmt.Sprintf("%d", len(customers)), Hint: "内存模式演示数据", Tone: "primary"},
			{Key: "services", Label: "运行中服务", Value: fmt.Sprintf("%d", activeServiceTotal), Hint: "当前服务状态分布", Tone: "success"},
			{Key: "unpaidInvoices", Label: "待支付账单", Value: fmt.Sprintf("%d", unpaidInvoiceTotal), Hint: "待催收账单数量", Tone: "warning"},
			{Key: "pendingIdentity", Label: "待审核实名", Value: fmt.Sprintf("%d", pendingIdentityTotal), Hint: "等待人工审核", Tone: "warning"},
			{Key: "openTickets", Label: "处理中工单", Value: fmt.Sprintf("%d", openTicketTotal), Hint: "需要客服继续跟进", Tone: "danger"},
			{Key: "todayIncome", Label: "今日收款", Value: formatCurrency(0), Hint: "内存模式不统计支付日维度", Tone: "primary"},
		},
		RiskCards: []reportDTO.MetricCard{
			{Key: "overdueInvoices", Label: "逾期账单", Value: fmt.Sprintf("%d", overdueInvoiceTotal), Hint: "账单已超过到期时间", Tone: "danger"},
			{Key: "expiringServices", Label: "7 日内到期服务", Value: fmt.Sprintf("%d", expiringServiceTotal), Hint: "建议提前续费触达", Tone: "warning"},
			{Key: "suspendedServices", Label: "已暂停服务", Value: fmt.Sprintf("%d", suspendedServiceTotal), Hint: "需要核对恢复或终止", Tone: "danger"},
		},
		ServiceStatus: []reportDTO.StatusBucket{
			{Name: "运行中", Count: activeServiceTotal},
			{Name: "已暂停", Count: suspendedServiceTotal},
		},
		PendingIdentities: pendingIdentities,
		OverdueInvoices:   overdueInvoices,
		ExpiringServices:  expiringServices,
		OpenTickets:       openTickets,
		RecentAudits:      mapAuditEntries(service.audit.List()),
	}
}

func (service *Service) getReportOverviewFromMemory() reportDTO.ReportOverviewResponse {
	invoices, _ := service.orders.ListInvoices(orderDomain.InvoiceListFilter{
		Page:  1,
		Limit: 0,
		Sort:  "created_at",
		Order: "desc",
	})
	services, _ := service.orders.ListServices(orderDomain.ServiceListFilter{
		Page:  1,
		Limit: 0,
		Sort:  "created_at",
		Order: "desc",
	})
	orders, _ := service.orders.ListOrders(orderDomain.OrderListFilter{
		Page:  1,
		Limit: 0,
		Sort:  "create_time",
		Order: "desc",
	})
	customers := service.customers.List()

	var unpaidAmount float64
	invoiceStatus := map[string]int64{}
	serviceStatus := map[string]int64{}
	cycleStats := map[string]reportDTO.StatusBucket{}
	groupStats := map[string]int64{}
	productStats := map[string]reportDTO.RankItem{}
	receivableStats := map[int64]reportDTO.RankItem{}

	for _, invoice := range invoices {
		invoiceStatus[string(invoice.Status)]++
		if string(invoice.Status) == "UNPAID" {
			unpaidAmount += invoice.TotalAmount
			item := receivableStats[invoice.CustomerID]
			item.Name = lookupCustomerName(customers, invoice.CustomerID)
			item.Count++
			item.Amount += invoice.TotalAmount
			item.Extra = invoice.DueAt
			receivableStats[invoice.CustomerID] = item
		}
	}

	for _, item := range services {
		serviceStatus[string(item.Status)]++
	}

	for _, item := range orders {
		bucket := cycleStats[item.BillingCycle]
		bucket.Name = billingCycleLabel(item.BillingCycle)
		bucket.Count++
		bucket.Amount += item.Amount
		cycleStats[item.BillingCycle] = bucket

		product := productStats[item.ProductName]
		product.Name = item.ProductName
		product.Count++
		product.Amount += item.Amount
		productStats[item.ProductName] = product
	}

	for _, customer := range customers {
		groupName := firstNonEmpty(customer.GroupName, "未分组")
		groupStats[groupName]++
	}

	return reportDTO.ReportOverviewResponse{
		Headline: []reportDTO.MetricCard{
			{Key: "monthIncome", Label: "本月收款", Value: formatCurrency(0), Hint: "内存模式不统计月度流水", Tone: "primary"},
			{Key: "monthRefund", Label: "本月退款", Value: formatCurrency(0), Hint: "内存模式不统计退款趋势", Tone: "warning"},
			{Key: "netIncome", Label: "本月净收入", Value: formatCurrency(0), Hint: "等待切换 MySQL 真数据", Tone: "success"},
			{Key: "paidInvoices", Label: "已支付账单", Value: fmt.Sprintf("%d", invoiceStatus["PAID"]), Hint: "内存模式统计", Tone: "primary"},
			{Key: "unpaidAmount", Label: "待收金额", Value: formatCurrency(unpaidAmount), Hint: "UNPAID 账单金额", Tone: "danger"},
			{Key: "activeServices", Label: "运行中服务", Value: fmt.Sprintf("%d", serviceStatus["ACTIVE"]), Hint: fmt.Sprintf("暂停中 %d", serviceStatus["SUSPENDED"]), Tone: "success"},
		},
		InvoiceStatus:  sortBucketsMap(invoiceStatus, invoiceLabel),
		ServiceStatus:  sortBucketsMap(serviceStatus, serviceLabel),
		BillingCycles:  sortCycleStats(cycleStats),
		CustomerGroups: sortGroupStats(groupStats),
		TopProducts:    sortRankStats(productStats),
		TopReceivables: sortReceivableStats(receivableStats),
	}
}

func (service *Service) queryCount(query string, args ...any) int64 {
	var value sql.NullInt64
	if err := service.db.QueryRowContext(context.Background(), query, args...).Scan(&value); err != nil {
		return 0
	}
	if !value.Valid {
		return 0
	}
	return value.Int64
}

func (service *Service) queryFloat(query string, args ...any) float64 {
	var value sql.NullFloat64
	if err := service.db.QueryRowContext(context.Background(), query, args...).Scan(&value); err != nil {
		return 0
	}
	if !value.Valid {
		return 0
	}
	return value.Float64
}

func (service *Service) queryTrend(query string, days int) []reportDTO.TrendPoint {
	start := time.Now().AddDate(0, 0, -(days - 1))
	rows, err := service.db.QueryContext(context.Background(), query, start)
	if err != nil {
		return buildEmptyTrend(days)
	}
	defer rows.Close()

	type point struct {
		amount float64
		count  int64
	}
	points := make(map[string]point)
	for rows.Next() {
		var (
			day    time.Time
			amount sql.NullFloat64
			count  sql.NullInt64
		)
		if err := rows.Scan(&day, &amount, &count); err != nil {
			continue
		}
		key := day.Format("2006-01-02")
		points[key] = point{amount: nullableFloat(amount), count: nullableInt(count)}
	}

	result := make([]reportDTO.TrendPoint, 0, days)
	for index := days - 1; index >= 0; index-- {
		day := time.Now().AddDate(0, 0, -index)
		key := day.Format("2006-01-02")
		current := points[key]
		result = append(result, reportDTO.TrendPoint{
			Label:  day.Format("01-02"),
			Amount: current.amount,
			Count:  current.count,
		})
	}
	return result
}

func (service *Service) queryStatusBuckets(query string, labeler func(string) string) []reportDTO.StatusBucket {
	rows, err := service.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.StatusBucket, 0)
	for rows.Next() {
		var (
			name   sql.NullString
			count  sql.NullInt64
			amount sql.NullFloat64
		)
		if err := rows.Scan(&name, &count, &amount); err != nil {
			continue
		}
		result = append(result, reportDTO.StatusBucket{
			Name:   labeler(name.String),
			Count:  nullableInt(count),
			Amount: nullableFloat(amount),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Amount == result[j].Amount {
			return result[i].Count > result[j].Count
		}
		return result[i].Amount > result[j].Amount
	})
	return result
}

func (service *Service) queryRankItems(query string) []reportDTO.RankItem {
	rows, err := service.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.RankItem, 0)
	for rows.Next() {
		var (
			name   sql.NullString
			count  sql.NullInt64
			amount sql.NullFloat64
			extra  sql.NullString
		)
		if err := rows.Scan(&name, &count, &amount, &extra); err != nil {
			continue
		}
		result = append(result, reportDTO.RankItem{
			Name:   name.String,
			Count:  nullableInt(count),
			Amount: nullableFloat(amount),
			Extra:  extra.String,
		})
	}
	return result
}

func (service *Service) queryPendingIdentities() []reportDTO.PendingIdentityItem {
	rows, err := service.db.QueryContext(context.Background(), `
SELECT c.id, c.customer_no, c.name, ci.subject_name,
       IFNULL(DATE_FORMAT(ci.submitted_at, '%Y-%m-%d %H:%i:%s'), '')
FROM customer_identities ci
INNER JOIN customers c ON c.id = ci.customer_id
WHERE ci.verify_status = 'PENDING'
ORDER BY ci.submitted_at ASC, ci.id ASC
LIMIT 6`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.PendingIdentityItem, 0)
	for rows.Next() {
		var item reportDTO.PendingIdentityItem
		if err := rows.Scan(&item.CustomerID, &item.CustomerNo, &item.CustomerName, &item.SubjectName, &item.SubmittedAt); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (service *Service) queryOverdueInvoices() []reportDTO.OverdueInvoiceItem {
	rows, err := service.db.QueryContext(context.Background(), `
SELECT i.id, i.invoice_no, c.name, i.total_amount,
       IFNULL(DATE_FORMAT(i.due_at, '%Y-%m-%d %H:%i:%s'), ''),
       GREATEST(TIMESTAMPDIFF(DAY, i.due_at, NOW()), 0)
FROM invoices i
INNER JOIN customers c ON c.id = i.customer_id
WHERE i.status = 'UNPAID' AND i.due_at IS NOT NULL AND i.due_at < NOW()
ORDER BY i.due_at ASC
LIMIT 6`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.OverdueInvoiceItem, 0)
	for rows.Next() {
		var item reportDTO.OverdueInvoiceItem
		if err := rows.Scan(&item.InvoiceID, &item.InvoiceNo, &item.CustomerName, &item.Amount, &item.DueAt, &item.DaysOverdue); err != nil {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (service *Service) queryExpiringServices() []reportDTO.ExpiringServiceItem {
	rows, err := service.db.QueryContext(context.Background(), `
SELECT s.id, s.service_no, c.name, s.product_name, s.status,
       IFNULL(DATE_FORMAT(s.next_due_at, '%Y-%m-%d %H:%i:%s'), ''),
       GREATEST(TIMESTAMPDIFF(DAY, NOW(), s.next_due_at), 0)
FROM services s
INNER JOIN customers c ON c.id = s.customer_id
WHERE s.status IN ('ACTIVE', 'SUSPENDED')
  AND s.next_due_at IS NOT NULL
  AND s.next_due_at BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 7 DAY)
ORDER BY s.next_due_at ASC
LIMIT 6`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.ExpiringServiceItem, 0)
	for rows.Next() {
		var (
			item       reportDTO.ExpiringServiceItem
			statusCode string
		)
		if err := rows.Scan(&item.ServiceID, &item.ServiceNo, &item.CustomerName, &item.ProductName, &statusCode, &item.NextDueAt, &item.DaysRemaining); err != nil {
			continue
		}
		item.Status = serviceLabel(statusCode)
		result = append(result, item)
	}
	return result
}

func (service *Service) queryOpenTickets() []reportDTO.TicketTodoItem {
	rows, err := service.db.QueryContext(context.Background(), `
SELECT t.id, t.ticket_no, c.name, t.title, t.status,
       IFNULL(DATE_FORMAT(t.updated_at, '%Y-%m-%d %H:%i:%s'), '')
FROM tickets t
INNER JOIN customers c ON c.id = t.customer_id
WHERE t.status IN ('OPEN', 'PROCESSING')
ORDER BY t.updated_at DESC
LIMIT 6`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]reportDTO.TicketTodoItem, 0)
	for rows.Next() {
		var (
			item       reportDTO.TicketTodoItem
			statusCode string
		)
		if err := rows.Scan(&item.TicketID, &item.TicketNo, &item.CustomerName, &item.Title, &statusCode, &item.UpdatedAt); err != nil {
			continue
		}
		item.Status = ticketLabel(statusCode)
		result = append(result, item)
	}
	return result
}

func (service *Service) queryRecentAudits() []reportDTO.AuditTimelineItem {
	rows, err := service.db.QueryContext(context.Background(), `
SELECT id, actor, action, target, COALESCE(description, ''),
       IFNULL(DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), '')
FROM audit_logs
ORDER BY created_at DESC
LIMIT 8`)
	if err != nil {
		return mapAuditEntries(service.audit.List())
	}
	defer rows.Close()

	result := make([]reportDTO.AuditTimelineItem, 0)
	for rows.Next() {
		var (
			item       reportDTO.AuditTimelineItem
			actionCode string
			target     string
			desc       string
		)
		if err := rows.Scan(&item.ID, &item.Actor, &actionCode, &target, &desc, &item.CreatedAt); err != nil {
			continue
		}
		item.Action = auditActionLabel(actionCode)
		item.Target = target
		item.Description = normalizeAuditDescription(actionCode, target, desc)
		result = append(result, item)
	}
	return result
}

func buildEmptyTrend(days int) []reportDTO.TrendPoint {
	result := make([]reportDTO.TrendPoint, 0, days)
	for index := days - 1; index >= 0; index-- {
		day := time.Now().AddDate(0, 0, -index)
		result = append(result, reportDTO.TrendPoint{
			Label: day.Format("01-02"),
		})
	}
	return result
}

func mapAuditEntries(items []audit.Entry) []reportDTO.AuditTimelineItem {
	limit := len(items)
	if limit > 8 {
		limit = 8
	}
	result := make([]reportDTO.AuditTimelineItem, 0, limit)
	for index := 0; index < limit; index++ {
		item := items[index]
		result = append(result, reportDTO.AuditTimelineItem{
			ID:          item.ID,
			Actor:       item.Actor,
			Action:      auditActionLabel(item.Action),
			Target:      item.Target,
			Description: normalizeAuditDescription(item.Action, item.Target, item.Description),
			CreatedAt:   item.CreatedAt,
		})
	}
	return result
}

func sortBucketsMap(items map[string]int64, labeler func(string) string) []reportDTO.StatusBucket {
	result := make([]reportDTO.StatusBucket, 0, len(items))
	for key, count := range items {
		result = append(result, reportDTO.StatusBucket{Name: labeler(key), Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func sortCycleStats(items map[string]reportDTO.StatusBucket) []reportDTO.StatusBucket {
	result := make([]reportDTO.StatusBucket, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Amount > result[j].Amount
	})
	return result
}

func sortGroupStats(items map[string]int64) []reportDTO.StatusBucket {
	result := make([]reportDTO.StatusBucket, 0, len(items))
	for key, count := range items {
		result = append(result, reportDTO.StatusBucket{Name: key, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func sortRankStats(items map[string]reportDTO.RankItem) []reportDTO.RankItem {
	result := make([]reportDTO.RankItem, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Amount > result[j].Amount
	})
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func sortReceivableStats(items map[int64]reportDTO.RankItem) []reportDTO.RankItem {
	result := make([]reportDTO.RankItem, 0, len(items))
	for _, item := range items {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Amount > result[j].Amount
	})
	if len(result) > 5 {
		result = result[:5]
	}
	return result
}

func formatCurrency(value float64) string {
	return fmt.Sprintf("¥%.2f", value)
}

func nullableFloat(value sql.NullFloat64) float64 {
	if !value.Valid {
		return 0
	}
	return value.Float64
}

func nullableInt(value sql.NullInt64) int64 {
	if !value.Valid {
		return 0
	}
	return value.Int64
}

func parseDateTime(value string) (time.Time, error) {
	layouts := []string{"2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime")
}

func lookupCustomerName(customers []customerDomain.Customer, customerID int64) string {
	for _, item := range customers {
		if item.ID == customerID {
			return item.Name
		}
	}
	return fmt.Sprintf("客户 #%d", customerID)
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func serviceLabel(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "PENDING":
		return "待开通"
	case "ACTIVE":
		return "运行中"
	case "SUSPENDED":
		return "已暂停"
	case "TERMINATED":
		return "已终止"
	default:
		return firstNonEmpty(value, "未知")
	}
}

func invoiceLabel(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "UNPAID":
		return "待支付"
	case "PAID":
		return "已支付"
	case "REFUNDED":
		return "已退款"
	default:
		return firstNonEmpty(value, "未知")
	}
}

func billingCycleLabel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "monthly":
		return "月付"
	case "quarterly":
		return "季付"
	case "annual":
		return "年付"
	default:
		return firstNonEmpty(value, "其他周期")
	}
}

func channelLabel(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "ONLINE":
		return "在线支付"
	case "OFFLINE":
		return "线下收款"
	case "ALIPAY":
		return "支付宝"
	case "WECHAT":
		return "微信支付"
	default:
		return firstNonEmpty(value, "其他渠道")
	}
}

func ticketLabel(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "OPEN":
		return "待处理"
	case "PROCESSING":
		return "处理中"
	case "CLOSED":
		return "已关闭"
	default:
		return firstNonEmpty(value, "未知")
	}
}

func passthroughLabel(value string) string {
	return firstNonEmpty(value, "未分类")
}

func auditActionLabel(action string) string {
	switch strings.TrimSpace(action) {
	case "customer.bootstrap":
		return "初始化客户"
	case "order.activate":
		return "订单激活"
	case "customer.create":
		return "创建客户"
	case "customer.update":
		return "更新客户"
	case "customer.identity.review":
		return "实名审核"
	case "customer.contact.create":
		return "新增联系人"
	case "customer.contact.update":
		return "更新联系人"
	case "customer.contact.delete":
		return "删除联系人"
	case "order.checkout":
		return "提交订单"
	case "invoice.pay":
		return "支付账单"
	case "invoice.receive_payment":
		return "登记收款"
	case "invoice.refund":
		return "执行退款"
	case "service.activate":
		return "开通服务"
	case "service.suspend":
		return "暂停服务"
	case "service.terminate":
		return "终止服务"
	case "service.reboot":
		return "重启实例"
	case "service.reset-password":
		return "重置密码"
	case "service.reinstall":
		return "重装系统"
	default:
		if action == "" {
			return "系统操作"
		}
		return action
	}
}

func normalizeAuditDescription(action, target, description string) string {
	if strings.TrimSpace(description) == "" {
		return fmt.Sprintf("%s：%s", auditActionLabel(action), firstNonEmpty(target, "系统对象"))
	}
	if strings.Contains(description, "鏃") || strings.Contains(description, "鍚") || strings.Contains(description, "瀹") {
		return fmt.Sprintf("%s：%s", auditActionLabel(action), firstNonEmpty(target, "系统对象"))
	}
	return description
}
