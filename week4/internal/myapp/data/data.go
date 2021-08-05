package data

import (
	"fmt"
	"week4/cmd/myapp/config"
	"week4/pkg/mysql"
)

// Data data层持久化资源对象
type Data struct {
	Orm mysql.Repository
}

// NewData 创建data层持久化资源对象的方法
func NewData() (Data, func(), error) {
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

	data := Data{Orm: *orm}

	cleanup := func() {
		fmt.Println("closing mysql connection")
	}

	return data, cleanup, nil
}
