package app

import (
	"database/sql"
	"fmt"
	"log"

	"idc-finance/internal/platform/config"
	"idc-finance/internal/platform/server"
	mysqlStore "idc-finance/internal/platform/store/mysql"

	"github.com/gin-gonic/gin"
)

type Application struct {
	config config.AppConfig
	router *gin.Engine
	db     *sql.DB
	err    error
}

func New() *Application {
	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	app := &Application{
		config: cfg,
	}

	if cfg.StorageDriver == "mysql" {
		db, err := mysqlStore.Open(cfg.MySQLDSN)
		if err != nil {
			if cfg.StorageStrict {
				app.err = fmt.Errorf("mysql bootstrap failed in strict mode: %w", err)
			} else {
				log.Printf("mysql bootstrap failed, fallback to memory: %v", err)
			}
		} else {
			app.db = db
			log.Printf("mysql bootstrap connected")
		}
	}

	app.router = server.NewRouter(cfg, app.db)

	return app
}

func (a *Application) Run() error {
	if a.err != nil {
		return a.err
	}

	if a.db != nil {
		defer func() {
			if err := a.db.Close(); err != nil {
				log.Printf("mysql close failed: %v", err)
			}
		}()
	}

	if a.config.StorageDriver != "memory" && a.config.StorageDriver != "mysql" {
		return fmt.Errorf("unsupported storage driver: %s", a.config.StorageDriver)
	}

	log.Printf("starting api: storageDriver=%s mysqlReady=%t strict=%t addr=%s", a.config.StorageDriver, a.db != nil, a.config.StorageStrict, a.config.HTTPAddr)
	return a.router.Run(a.config.HTTPAddr)
}
