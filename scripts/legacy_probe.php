<?php

declare(strict_types=1);

mysqli_report(MYSQLI_REPORT_ERROR | MYSQLI_REPORT_STRICT);

$host = '127.0.0.1';
$port = 3308;
$user = 'zjmf_legacy';
$password = 'ZjmfLegacy!2026';
$database = 'zjmf_legacy374';
$authCode = 'l5TKhHZiSrITLTx5dn';

function db(): mysqli
{
    global $host, $port, $user, $password, $database;
    static $conn = null;
    if ($conn instanceof mysqli) {
        return $conn;
    }
    $conn = new mysqli($host, $user, $password, $database, $port);
    $conn->set_charset('utf8mb4');
    return $conn;
}

function cmf_password_local(string $plain): string
{
    global $authCode;
    return '###' . md5(md5($authCode . $plain));
}

function fetchAllAssoc(mysqli_result $result): array
{
    $rows = [];
    while ($row = $result->fetch_assoc()) {
        $rows[] = $row;
    }
    return $rows;
}

function queryValue(string $sql)
{
    $result = db()->query($sql);
    $row = $result->fetch_row();
    return $row[0] ?? null;
}

function showState(): void
{
    $queries = [
        'clients' => 'select id,username,email,status,phonenumber from shd_clients order by id asc limit 20',
        'products' => 'select id,name,type,api_type,upstream_pid,zjmf_api_id from shd_products order by id asc limit 20',
        'orders' => 'select id,uid,ordernum,status,invoiceid,amount from shd_orders order by id desc limit 20',
        'invoices' => 'select id,uid,status,payment_status,total,type,url from shd_invoices order by id desc limit 20',
        'hosts' => 'select id,uid,orderid,productid,domain,domainstatus,dcimid,dedicatedip,nextduedate from shd_host order by id desc limit 20',
    ];
    foreach ($queries as $name => $sql) {
        echo "## {$name}\n";
        $rows = fetchAllAssoc(db()->query($sql));
        echo json_encode($rows, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
    }
}

function ensureDemoClient(): void
{
    $email = 'legacy-demo@example.com';
    $plainPassword = 'Legacy123!';
    $conn = db();
    $stmt = $conn->prepare('select id,password from shd_clients where email = ? limit 1');
    $stmt->bind_param('s', $email);
    $stmt->execute();
    $row = $stmt->get_result()->fetch_assoc();
    if ($row) {
        echo json_encode([
            'created' => false,
            'client_id' => (int) $row['id'],
            'email' => $email,
            'plain_password' => $plainPassword,
        ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        return;
    }

    $now = time();
    $username = 'legacydemo';
    $avatar = '用户头像2-12.jpg';
    $phoneCode = 86;
    $phone = '13800138000';
    $currency = (int) (queryValue('select id from shd_currencies where `default` = 1 limit 1') ?: 1);
    $defaultGateway = (string) (queryValue('select `gateway` from shd_payment_gateways where setting = "show" and value = "on" order by `order` asc limit 1') ?: '');
    $apiPassword = base64_encode('legacy-demo-api-password');
    $passwordHash = cmf_password_local($plainPassword);

    $uuid = md5($email . 'legacy-demo');
    $status = 1;
    $saleId = 0;
    $marketingEmails = 0;
    $apiOpen = 1;
    $clientSql = sprintf(
        "insert into shd_clients (
            uuid,username,email,avatar,phone_code,phonenumber,currency,password,lastloginip,create_time,lastlogin,
            status,defaultgateway,sale_id,qq,companyname,address1,api_password,marketing_emails_opt_in,api_open,api_create_time
        ) values (
            '%s','%s','%s','%s',%d,'%s',%d,'%s','',%d,%d,%d,'%s',%d,'','','','%s',%d,%d,%d
        )",
        $conn->real_escape_string($uuid),
        $conn->real_escape_string($username),
        $conn->real_escape_string($email),
        $conn->real_escape_string($avatar),
        $phoneCode,
        $conn->real_escape_string($phone),
        $currency,
        $conn->real_escape_string($passwordHash),
        $now,
        $now,
        $status,
        $conn->real_escape_string($defaultGateway),
        $saleId,
        $conn->real_escape_string($apiPassword),
        $marketingEmails,
        $apiOpen,
        $now
    );
    $conn->query($clientSql);
    $clientId = $conn->insert_id;

    echo json_encode([
        'created' => true,
        'client_id' => $clientId,
        'email' => $email,
        'plain_password' => $plainPassword,
    ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
}

function ensureDemoUpstreamHost(): void
{
    $clientEmail = 'legacy-demo@example.com';
    $conn = db();
    $clientStmt = $conn->prepare('select id from shd_clients where email = ? limit 1');
    $clientStmt->bind_param('s', $clientEmail);
    $clientStmt->execute();
    $client = $clientStmt->get_result()->fetch_assoc();
    if (!$client) {
        throw new RuntimeException('demo client not found');
    }
    $uid = (int) $client['id'];

    $productId = (int) (queryValue('select id from shd_products where upstream_pid = 96 order by id asc limit 1') ?: 2);
    $checkStmt = $conn->prepare('select id,orderid,dcimid from shd_host where uid = ? and productid = ? and dcimid = 8346 limit 1');
    $dcimId = 8346;
    $checkStmt->bind_param('ii', $uid, $productId);
    $checkStmt->execute();
    $existing = $checkStmt->get_result()->fetch_assoc();
    if ($existing) {
        echo json_encode([
            'created' => false,
            'host_id' => (int) $existing['id'],
            'order_id' => (int) $existing['orderid'],
            'dcimid' => (int) $existing['dcimid'],
            'service_url' => 'http://127.0.0.1:8091/servicedetail?id=' . (int) $existing['id'],
        ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        return;
    }

    $now = time();
    $due = 1776321001;
    $invoiceNum = 'INV-LEGACY-DEMO-' . date('YmdHis', $now);
    $orderNum = 'ORD-LEGACY-DEMO-' . date('YmdHis', $now);

    $subtotal = '304.00';
    $zeroDec = '0.00';
    $status = 'Paid';
    $payment = 'EpayAlipay';
    $type = 'product';
    $paymentStatus = 'Paid';
    $url = '';
    $empty = '';
    $zero = 0;
    $invoiceSql = sprintf(
        "insert into shd_invoices (
            uid,invoice_num,create_time,update_time,due_time,paid_time,last_capture_attempt,
            subtotal,credit,tax,tax2,total,taxrate,taxrate2,status,payment,notes,type,payment_status,
            due_email_times,is_cron,suffix,use_credit_limit,invoice_id,url,paymt,is_delete,credit_limit_prepayment
        ) values (
            %d,'%s',%d,%d,%d,%d,%d,%s,%s,%s,%s,%s,%s,%s,'%s','%s','%s','%s','%s',
            %d,%d,%d,%d,%d,'%s','%s',%d,%d
        )",
        $uid,
        $conn->real_escape_string($invoiceNum),
        $now,
        $now,
        $due,
        $now,
        $zero,
        $subtotal,
        $zeroDec,
        $zeroDec,
        $zeroDec,
        $subtotal,
        $zeroDec,
        $zeroDec,
        $conn->real_escape_string($status),
        $conn->real_escape_string($payment),
        $empty,
        $conn->real_escape_string($type),
        $conn->real_escape_string($paymentStatus),
        $zero,
        $zero,
        $zero,
        $zero,
        $zero,
        $url,
        $empty,
        $zero,
        $zero
    );
    $conn->query($invoiceSql);
    $invoiceId = $conn->insert_id;

    $orderStatus = 'Active';
    $promoCode = '';
    $promoType = '';
    $promoValue = '0.00';
    $notes = 'legacy local study demo';
    $orderSql = sprintf(
        "insert into shd_orders (
            uid,ordernum,status,pay_time,create_time,update_time,amount,payment,promo_code,promo_type,promo_value,invoiceid,delete_time,notes
        ) values (
            %d,'%s','%s',%d,%d,%d,%s,'%s','%s','%s',%s,%d,%d,'%s'
        )",
        $uid,
        $conn->real_escape_string($orderNum),
        $conn->real_escape_string($orderStatus),
        $now,
        $now,
        $now,
        $subtotal,
        $conn->real_escape_string($payment),
        $promoCode,
        $promoType,
        $promoValue,
        $invoiceId,
        $zero,
        $conn->real_escape_string($notes)
    );
    $conn->query($orderSql);
    $orderId = $conn->insert_id;

    $domain = 'Xxun298067534754';
    $username = 'root';
    $password = '';
    $domainStatus = 'Active';
    $assignedIps = '221.236.21.196';
    $os = 'Ubuntu-24.04.1-x64';
    $osUrl = '';
    $reinstallInfo = '';
    $remark = 'legacy upstream demo host';
    $showLastActMessage = 1;
    $port = 46094;
    $percentValue = '120.00';
    $upstreamCost = '';
    $upstreamConfig = json_encode([
        'upstream_host_id' => 8346,
        'upstream_provider_id' => 7469,
        'source' => 'local-study-seed',
    ], JSON_UNESCAPED_UNICODE);
    $billingCycle = 'monthly';
    $hostSql = sprintf(
        "insert into shd_host (
            uid,orderid,productid,serverid,regdate,domain,payment,firstpaymentamount,amount,billingcycle,last_settle,
            nextduedate,nextinvoicedate,termination_date,completed_date,domainstatus,username,password,notes,promoid,
            overideautosuspend,overidesuspenduntil,ns1,ns2,diskusage,disklimit,bwusage,bwlimit,user_cate_id,lastupdate,
            create_time,update_time,suspend_time,auto_terminate_end_cycle,auto_terminate_reason,dedicatedip,assignedips,
            dcimid,dcim_os,os,os_url,reinstall_info,remark,show_last_act_message,port,dcim_area,flag,flag_cycle,stream_info,
            initiative_renew,percent_value,agent_client,upstream_cost,upstream_configoption
        ) values (
            %d,%d,%d,%d,%d,'%s','%s',%s,%s,'%s',%d,
            %d,%d,%d,%d,'%s','%s','%s','%s',%d,
            %d,%d,'%s','%s',%d,%d,%s,%d,%d,%d,
            %d,%d,%d,%d,'%s','%s','%s',
            %d,%d,'%s','%s','%s','%s',%d,%d,%d,%d,'%s','%s',
            %d,%s,%d,'%s','%s'
        )",
        $uid,
        $orderId,
        $productId,
        $zero,
        $now,
        $conn->real_escape_string($domain),
        $conn->real_escape_string($payment),
        $subtotal,
        $subtotal,
        $billingCycle,
        $zero,
        $due,
        $due,
        $zero,
        $zero,
        $conn->real_escape_string($domainStatus),
        $conn->real_escape_string($username),
        $conn->real_escape_string($password),
        $conn->real_escape_string($notes),
        $zero,
        $zero,
        $zero,
        $empty,
        $empty,
        $zero,
        $zero,
        $zeroDec,
        $zero,
        $zero,
        $now,
        $now,
        $now,
        $zero,
        $zero,
        $empty,
        $conn->real_escape_string($assignedIps),
        $conn->real_escape_string($assignedIps),
        $dcimId,
        $zero,
        $conn->real_escape_string($os),
        $conn->real_escape_string($osUrl),
        $conn->real_escape_string($reinstallInfo),
        $conn->real_escape_string($remark),
        $showLastActMessage,
        $port,
        $zero,
        $zero,
        $empty,
        $conn->real_escape_string($empty),
        $zero,
        $percentValue,
        $zero,
        $conn->real_escape_string($upstreamCost),
        $conn->real_escape_string($upstreamConfig)
    );
    $conn->query($hostSql);
    $hostId = $conn->insert_id;

    $description = "ECS - 成都 · 西信 | 弹性云 (" . date('Y-m-d H', $now) . ' - ' . date('Y-m-d H', $due) . ')';
    $invoiceItemType = 'host';
    $taxed = 0;
    $invoiceItemSql = sprintf(
        "insert into shd_invoice_items (
            invoice_id,uid,type,rel_id,description,description2,amount,taxed,due_time,payment,notes,delete_time
        ) values (
            %d,%d,'%s',%d,'%s','%s',%s,%d,%d,'%s','%s',%d
        )",
        $invoiceId,
        $uid,
        $conn->real_escape_string($invoiceItemType),
        $hostId,
        $conn->real_escape_string($description),
        $conn->real_escape_string($description),
        $subtotal,
        $taxed,
        $due,
        $conn->real_escape_string($payment),
        $empty,
        $zero
    );
    $conn->query($invoiceItemSql);

    $conn->query("update shd_invoices set url = 'servicedetail?id={$hostId}' where id = {$invoiceId}");

    echo json_encode([
        'created' => true,
        'client_id' => $uid,
        'invoice_id' => $invoiceId,
        'order_id' => $orderId,
        'host_id' => $hostId,
        'dcimid' => $dcimId,
        'service_url' => 'http://127.0.0.1:8091/servicedetail?id=' . $hostId,
    ], JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
}

$action = $argv[1] ?? 'state';

switch ($action) {
    case 'state':
        showState();
        break;
    case 'ensure-demo-client':
        ensureDemoClient();
        break;
    case 'ensure-demo-host':
        ensureDemoUpstreamHost();
        break;
    default:
        fwrite(STDERR, "unknown action: {$action}\n");
        exit(1);
}
