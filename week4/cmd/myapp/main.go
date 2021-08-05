package main

import (
	"fmt"
	"week4/cmd/myapp/config"
	"week4/internal/myapp"
	"week4/pkg/prometheus"
)

func init() {
	// 解析启动参数
	config.Setup()

	// 初始化数据库
	//if err := mysql.CreateRepository(&mysql.Config{
	//	HostAddress:   config.Opts.MysqlAddress,
	//	Username:      config.Opts.MysqlUser,
	//	Password:      config.Opts.MysqlPassword,
	//	Database:      config.Opts.MysqlDatabase,
	//	EnableLogging: config.Opts.Debug,
	//}); err != nil {
	//	panic(err)
	//}
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- prometheus.NewPrometheus(config.Opts.PromPort, stop)
	}()
	go func() {
		appControllers, _,_ := initApp()
		done <- myapp.NewApp(config.Opts.Port, appControllers, stop)
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Printf("项目启动进程组出现错误, %+v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}
