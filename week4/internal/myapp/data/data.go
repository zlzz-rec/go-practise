package data

import (
	"fmt"
	"time"
	"week4/internal/myapp/config"
	"week4/pkg/mysql"
	"week4/pkg/redis"
)

// Data data层持久化资源对象
type Data struct {
	Orm   *mysql.Repository
	Redis *redis.RedisServer
}

// NewData 创建data层持久化资源对象的方法
func NewData() (Data, func(), error) {
	// mysql
	orm, err := mysql.CreateRepository(&mysql.Config{
		HostAddress:   config.Opts.MysqlAddress,
		Username:      config.Opts.MysqlUser,
		Password:      config.Opts.MysqlPassword,
		Database:      config.Opts.MysqlDatabase,
		EnableLogging: config.Opts.Debug,
	})
	if err != nil {
		return Data{}, nil, err
	}

	// redis
	maxIdleConnection := 50
	idleTimeout := time.Duration(180) * time.Second
	timeout := time.Duration(10) * time.Second
	cfg := redis.NewConfig(maxIdleConnection, idleTimeout, timeout, timeout, timeout, config.Opts.RedisDB, config.Opts.RedisPassword)
	normalServerTemplate := redis.NewNormalTemplate(config.Opts.RedisAddress, cfg)
	normalRedis := redis.NewRedisServer()
	normalRedis.SetRedisTemplate(normalServerTemplate)

	data := Data{Orm: orm, Redis: normalRedis}

	cleanup := func() {
		fmt.Println("TODO: closing data connection")
	}

	return data, cleanup, nil
}
