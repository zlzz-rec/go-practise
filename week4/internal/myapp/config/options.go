package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Opts *Options

func SetupOnce() {
	if Opts != nil {
		return
	} else {
		Opts = &Options{}
	}
	err := viper.BindPFlags(ParseCommandLine())
	if err != nil {
		panic(fmt.Sprintf("parse command args error, %s", err))
	}
	err = viper.Unmarshal(Opts)
	if err != nil {
		panic(fmt.Sprintf("marshal command args error, %s", err))
	}
}

// Options 启动参数
type Options struct {
	RedisAddress             string `mapstructure:"redis_address"`
	RedisPassword            string `mapstructure:"redis_password"`
	RedisDB                  int    `mapstructure:"redis_db"`
	MysqlAddress             string `mapstructure:"mysql_address"`
	MysqlUser                string `mapstructure:"mysql_user"`
	MysqlPassword            string `mapstructure:"mysql_password"`
	MysqlDatabase            string `mapstructure:"mysql_database"`
	MaxConnections           int    `mapstructure:"max_connections"`
	ConnMaxLifetime          int    `mapstructure:"conn_max_life_time"`
	MaxIdleConns             int    `mapstructure:"max_idle_conns"`
	EnableLogging            int    `mapstructure:"enable_logging"`
	MySQLAccessTimeout       int    `mapstructure:"mysql_access_timeout"`
	Port                     int    `mapstructure:"port"`
	PromPort                 int    `mapstructure:"prom_port"`
	RedisAccessTimeout       int    `mapstructure:"redis_access_timeout"`
	AccessControlAllowOrigin string `mapstructure:"access_control_allow_origin"`
	UrlSecret                string `mapstructure:"url_secret"`
	UrlSecretOn              bool   `mapstructure:"url_secret_on"`
	MD5Flag                  bool   `mapstructure:"md5_flag"`
	BaseURL                  string `mapstructure:"base_url"`
	Debug                    bool   `mapstructure:"debug"`
}
