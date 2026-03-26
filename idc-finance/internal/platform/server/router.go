package server

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strconv"

	automationHandler "idc-finance/internal/modules/automation/handler"
	automationRepository "idc-finance/internal/modules/automation/repository"
	automationService "idc-finance/internal/modules/automation/service"
	catalogHandler "idc-finance/internal/modules/catalog/handler"
	catalogRepository "idc-finance/internal/modules/catalog/repository"
	catalogService "idc-finance/internal/modules/catalog/service"
	customerHandler "idc-finance/internal/modules/customer/handler"
	customerRepository "idc-finance/internal/modules/customer/repository"
	customerService "idc-finance/internal/modules/customer/service"
	orderHandler "idc-finance/internal/modules/order/handler"
	orderRepository "idc-finance/internal/modules/order/repository"
	orderService "idc-finance/internal/modules/order/service"
	portalHandler "idc-finance/internal/modules/portal/handler"
	portalService "idc-finance/internal/modules/portal/service"
	providerHandler "idc-finance/internal/modules/provider/handler"
	providerRepository "idc-finance/internal/modules/provider/repository"
	providerService "idc-finance/internal/modules/provider/service"
	reportHandler "idc-finance/internal/modules/report/handler"
	reportService "idc-finance/internal/modules/report/service"
	ticketHandler "idc-finance/internal/modules/ticket/handler"
	ticketRepository "idc-finance/internal/modules/ticket/repository"
	ticketService "idc-finance/internal/modules/ticket/service"
	"idc-finance/internal/platform/audit"
	"idc-finance/internal/platform/auth"
	"idc-finance/internal/platform/config"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"
	"idc-finance/internal/platform/rbac"
	systemStore "idc-finance/internal/platform/system"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.AppConfig, db *sql.DB) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestID())
	registerCompiledWeb(engine, cfg.AdminWebDir, cfg.PortalWebDir)

	auditService := buildAuditService(cfg, db)
	taskService := buildAutomationService(cfg, db)
	systemService := systemStore.NewService(filepath.Join("data", "system-state.json"), auditService)

	customerRepo := buildCustomerRepository(cfg, db)
	customerSvc := customerService.New(customerRepo, auditService)
	customerHTTP := customerHandler.New(customerSvc)

	catalogRepo := buildCatalogRepository(cfg, db)
	providerRepo := buildProviderRepository(cfg, db)
	providerSvc := providerService.New(cfg.MofangCloud, cfg.FinanceUpstream, providerRepo, db, auditService, taskService)
	providerHTTP := providerHandler.New(providerSvc)
	catalogSvc := catalogService.New(catalogRepo, auditService, providerSvc, filepath.Join("data", "catalog-upstream-sync-history.json"))
	catalogHTTP := catalogHandler.New(catalogSvc)

	orderRepo := buildOrderRepository(cfg, db)
	orderSvc := orderService.New(orderRepo, catalogSvc, auditService, providerSvc, taskService)
	orderHTTP := orderHandler.New(orderSvc)
	automationHTTP := automationHandler.New(taskService, newAutomationRetryExecutor(providerSvc, orderSvc).Execute)

	reportSvc := reportService.New(db, customerSvc, orderSvc, auditService)
	reportHTTP := reportHandler.New(reportSvc)

	ticketRepo := buildTicketRepository(cfg, db)
	ticketSvc := ticketService.New(ticketRepo, customerSvc, orderSvc, auditService, taskService)
	ticketHTTP := ticketHandler.New(ticketSvc)

	portalSvc := portalService.New(customerSvc, orderSvc)
	portalHTTP := portalHandler.New(portalSvc)

	api := engine.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      "OK",
			"message":   "healthy",
			"requestId": middleware.GetRequestID(c),
			"data": gin.H{
				"service":       cfg.AppName,
				"version":       "0.1.0",
				"storageDriver": cfg.StorageDriver,
				"mysqlReady":    db != nil,
			},
		})
	})

	admin := api.Group("/admin")
	authGroup := admin.Group("/auth")
	auth.RegisterRoutes(authGroup)

	protected := admin.Group("")
	protected.Use(auth.AdminGuard())
	protected.GET("/menus", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(rbac.PhaseOneMenus(), middleware.GetRequestID(c)))
	})
	protected.GET("/permissions", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(systemStore.PermissionCodes(), middleware.GetRequestID(c)))
	})
	protected.GET("/audit-logs", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(auditService.List(), middleware.GetRequestID(c)))
	})
	protected.GET("/admin-members", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(systemService.ListAdmins(), middleware.GetRequestID(c)))
	})
	protected.POST("/admin-members", func(c *gin.Context) {
		var payload systemStore.SaveAdminInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest("无效的管理员参数", middleware.GetRequestID(c)))
			return
		}
		item, err := systemService.CreateAdmin(payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest(err.Error(), middleware.GetRequestID(c)))
			return
		}
		c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
	})
	protected.PATCH("/admin-members/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var payload systemStore.SaveAdminInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest("无效的管理员参数", middleware.GetRequestID(c)))
			return
		}
		item, err := systemService.UpdateAdmin(id, payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest(err.Error(), middleware.GetRequestID(c)))
			return
		}
		c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
	})
	protected.GET("/roles", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(systemService.ListRoles(), middleware.GetRequestID(c)))
	})
	protected.POST("/roles", func(c *gin.Context) {
		var payload systemStore.SaveRoleInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest("无效的角色参数", middleware.GetRequestID(c)))
			return
		}
		item, err := systemService.CreateRole(payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest(err.Error(), middleware.GetRequestID(c)))
			return
		}
		c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
	})
	protected.PATCH("/roles/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		var payload systemStore.SaveRoleInput
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest("无效的角色参数", middleware.GetRequestID(c)))
			return
		}
		item, err := systemService.UpdateRole(id, payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, appErrors.BadRequest(err.Error(), middleware.GetRequestID(c)))
			return
		}
		c.JSON(http.StatusOK, appErrors.Ok(item, middleware.GetRequestID(c)))
	})
	protected.GET("/admins", func(c *gin.Context) {
		c.JSON(http.StatusOK, appErrors.Ok(systemService.ListAdmins(), middleware.GetRequestID(c)))
	})
	protected.GET("/menu-debug/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		c.JSON(http.StatusOK, appErrors.Ok(gin.H{"requestedId": id}, middleware.GetRequestID(c)))
	})

	customerHTTP.RegisterRoutes(protected)
	catalogHTTP.RegisterAdminRoutes(protected)
	orderHTTP.RegisterAdminRoutes(protected)
	ticketHTTP.RegisterAdminRoutes(protected)
	providerHTTP.RegisterAdminRoutes(protected)
	automationHTTP.RegisterRoutes(protected)
	reportHTTP.RegisterRoutes(protected)

	portal := api.Group("/portal")
	portal.Use(auth.PortalGuard())
	catalogHTTP.RegisterPortalRoutes(portal)
	orderHTTP.RegisterPortalRoutes(portal)
	ticketHTTP.RegisterPortalRoutes(portal)
	portalHTTP.RegisterRoutes(portal)

	return engine
}

func buildCatalogRepository(cfg config.AppConfig, db *sql.DB) catalogRepository.Repository {
	if cfg.StorageDriver == "mysql" && db != nil {
		return catalogRepository.NewMySQLRepository(db)
	}
	return catalogRepository.NewMemoryRepository()
}

func buildAutomationService(cfg config.AppConfig, db *sql.DB) *automationService.Service {
	if cfg.StorageDriver == "mysql" && db != nil {
		return automationService.NewWithDB(automationRepository.NewMySQLRepository(db), db)
	}
	return automationService.New(automationRepository.NewMemoryRepository())
}

func buildAuditService(cfg config.AppConfig, db *sql.DB) *audit.Service {
	if cfg.StorageDriver == "mysql" && db != nil {
		return audit.New(audit.NewMySQLRepository(db))
	}
	return audit.New(audit.NewMemoryRepository())
}

func buildCustomerRepository(cfg config.AppConfig, db *sql.DB) customerRepository.Repository {
	if cfg.StorageDriver == "mysql" && db != nil {
		return customerRepository.NewMySQLRepository(db)
	}
	return customerRepository.NewMemoryRepository()
}

func buildOrderRepository(cfg config.AppConfig, db *sql.DB) orderRepository.Repository {
	if cfg.StorageDriver == "mysql" && db != nil {
		return orderRepository.NewMySQLRepository(db)
	}
	return orderRepository.NewMemoryRepository()
}

func buildProviderRepository(cfg config.AppConfig, db *sql.DB) providerRepository.Repository {
	if cfg.StorageDriver == "mysql" && db != nil {
		return providerRepository.NewMySQLRepository(db)
	}
	return providerRepository.NewMemoryRepository()
}

func buildTicketRepository(cfg config.AppConfig, db *sql.DB) ticketRepository.Repository {
	if cfg.StorageDriver == "mysql" && db != nil {
		return ticketRepository.NewMySQLRepository(db)
	}
	return ticketRepository.NewMemoryRepository()
}
