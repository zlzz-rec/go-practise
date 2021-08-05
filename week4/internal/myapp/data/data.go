package data

import (
	"fmt"
	config2 "week4/internal/myapp/config"
	"week4/pkg/mysql"
)

// Data data层持久化资源对象
type Data struct {
	Orm mysql.Repository
}

// NewData 创建data层持久化资源对象的方法
func NewData() (Data, func(), error) {
	orm, err := mysql.CreateRepository(&mysql.Config{
		HostAddress:   config2.Opts.MysqlAddress,
		Username:      config2.Opts.MysqlUser,
		Password:      config2.Opts.MysqlPassword,
		Database:      config2.Opts.MysqlDatabase,
		EnableLogging: config2.Opts.Debug,
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
