package config

import "os"

type AppConfig struct {
	AppName         string
	HTTPAddr        string
	GinMode         string
	SiteURL         string
	StorageDriver   string
	StorageStrict   bool
	MySQLDSN        string
	AdminWebDir     string
	PortalWebDir    string
	MofangCloud     MofangCloudConfig
	FinanceUpstream FinanceUpstreamConfig
}

type MofangCloudConfig struct {
	BaseURL            string
	Username           string
	Password           string
	Lang               string
	InsecureSkipVerify bool
	ListPath           string
	DetailPath         string
}

type FinanceUpstreamConfig struct {
	BaseURL            string
	Username           string
	Password           string
	SourceName         string
	InsecureSkipVerify bool
}

func Load() AppConfig {
	return AppConfig{
		AppName:       getEnv("APP_NAME", "wuqiongyun-idc-api"),
		HTTPAddr:      getEnv("HTTP_ADDR", ":18080"),
		GinMode:       getEnv("GIN_MODE", "debug"),
		SiteURL:       getEnv("SITE_URL", "http://127.0.0.1:18080"),
		StorageDriver: getEnv("STORAGE_DRIVER", "memory"),
		StorageStrict: getBoolEnv("STORAGE_STRICT", false),
		MySQLDSN:      getEnv("MYSQL_DSN", ""),
		AdminWebDir:   getEnv("ADMIN_WEB_DIR", "web/admin"),
		PortalWebDir:  getEnv("PORTAL_WEB_DIR", "web/portal"),
		MofangCloud: MofangCloudConfig{
			BaseURL:            getEnv("MOFANG_CLOUD_BASE_URL", ""),
			Username:           getEnv("MOFANG_CLOUD_USERNAME", ""),
			Password:           getEnv("MOFANG_CLOUD_PASSWORD", ""),
			Lang:               getEnv("MOFANG_CLOUD_LANG", "zh-cn"),
			InsecureSkipVerify: getBoolEnv("MOFANG_CLOUD_INSECURE_SKIP_VERIFY", true),
			ListPath:           getEnv("MOFANG_CLOUD_LIST_PATH", "/v1/clouds?page=1&per_page=100&sort=desc&orderby=id"),
			DetailPath:         getEnv("MOFANG_CLOUD_INSTANCE_DETAIL_PATH", "/v1/clouds/:id"),
		},
		FinanceUpstream: FinanceUpstreamConfig{
			BaseURL:            getEnv("FINANCE_UPSTREAM_BASE_URL", ""),
			Username:           getEnv("FINANCE_UPSTREAM_USERNAME", ""),
			Password:           getEnv("FINANCE_UPSTREAM_PASSWORD", ""),
			SourceName:         getEnv("FINANCE_UPSTREAM_SOURCE_NAME", "上游财务"),
			InsecureSkipVerify: getBoolEnv("FINANCE_UPSTREAM_INSECURE_SKIP_VERIFY", true),
		},
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	switch value {
	case "1", "true", "TRUE", "True", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "False", "no", "NO", "off", "OFF":
		return false
	default:
		return fallback
	}
}
