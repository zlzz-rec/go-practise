package data

import (
	"fmt"
	"week13/internal/myapp/config"
	"week13/pkg/mysql"
	"week13/pkg/redis"
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
	normalServerTemplate := redis.NewNormalTemplate(config.Opts.RedisAddress, config.Opts.RedisPassword, config.Opts.RedisDB)
	normalRedis := redis.NewRedisServer()
	normalRedis.SetRedisTemplate(normalServerTemplate)

	data := Data{Orm: orm, Redis: normalRedis}

	cleanup := func() {
		fmt.Println("TODO: closing data connection")
	}

	return data, cleanup, nil
}
