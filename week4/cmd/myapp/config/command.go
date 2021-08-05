package config

import "github.com/spf13/pflag"

// ParseCommandLine 解析命令行参数
func ParseCommandLine() *pflag.FlagSet {
	// server
	pflag.String("server_address", "0.0.0.0:8080", "service listen port")
	pflag.Int("port", 8080, "server port")
	pflag.Int("prom_port", 8399, "prometheus port")
	pflag.String("app_name", "", "app_name")
	pflag.String("access_control_allow_origin", "*", "access_control_allow_origin")
	pflag.Int("access_token_expired_seconds", 300, "access_token_expired_seconds")
	pflag.Int("refresh_token_expired_seconds", 1296000, "refresh_token_expired_seconds") //60 * 60 * 24 * 15
	pflag.String("url_secret", "", "url_secret")
	pflag.Bool("url_secret_on", false, "url_secret_on")
	pflag.Bool("debug", false, "debug")

	// mysql
	pflag.String("mysql_address", "", "mysql_address")
	pflag.String("mysql_user", "", "mysql_user")
	pflag.String("mysql_password", "", "mysql_password")
	pflag.String("mysql_database", "", "mysql_database")
	pflag.Int("mysql_access_timeout", 500, "mysql_access_timeout")
	pflag.Int("conn_max_life_time", 200, "conn_max_life_time")
	pflag.Int("max_idle_conns", 0, "max_idle_conns")

	// redis
	pflag.String("redis_address", "", "redis_address")
	pflag.String("redis_password", "", "redis_password")
	pflag.Int("redis_access_timeout", 500, "redis_access_timeout")

	pflag.Parse()
	return pflag.CommandLine
}
