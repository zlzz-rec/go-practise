package main

import (
	"fmt"
	"week4/internal/myapp/config"
	"week4/internal/myapp/controller"
	"week4/pkg/prometheus"
)

func init() {
	// 解析启动参数
	config.Setup()

	// 初始化资源
	controllerManager, _, err := Setup();
	if err != nil {
		panic(err)
	}
	controller.ControllerManager = controllerManager
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- prometheus.NewPrometheus(config.Opts.PromPort, stop)
	}()
	go func() {
		done <- controller.NewApp(config.Opts.Port, stop)
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
