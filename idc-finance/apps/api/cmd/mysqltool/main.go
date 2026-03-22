package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"idc-finance/internal/platform/config"
	mysqlStore "idc-finance/internal/platform/store/mysql"
)

func main() {
	mode := flag.String("mode", "all", "run mode: migrate|seed|all")
	root := flag.String("root", "", "project root, default current working directory")
	flag.Parse()

	cfg := config.Load()
	if cfg.MySQLDSN == "" {
		log.Fatal("MYSQL_DSN is required")
	}

	projectRoot := *root
	if projectRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		projectRoot = cwd
	}

	db, err := mysqlStore.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	switch *mode {
	case "migrate":
		err = mysqlStore.RunDirectory(db, filepath.Join(projectRoot, "migrations", "mysql"))
	case "seed":
		err = mysqlStore.RunDirectory(db, filepath.Join(projectRoot, "seed", "mysql"))
	case "all":
		if err = mysqlStore.RunDirectory(db, filepath.Join(projectRoot, "migrations", "mysql")); err == nil {
			err = mysqlStore.RunDirectory(db, filepath.Join(projectRoot, "seed", "mysql"))
		}
	default:
		log.Fatalf("unsupported mode: %s", *mode)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("mysql tool completed: mode=%s", *mode)
}
