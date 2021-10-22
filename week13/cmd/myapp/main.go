package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
	"week13/internal/myapp/config"
	"week13/internal/myapp/controller"
	"week13/pkg/kafka"
	"week13/pkg/log"
	"week13/pkg/prometheus"
	"week13/pkg/zerror"
)

func init() {
	// 解析启动参数
	config.SetupOnce()

	// 初始化分层资源
	controllerManager, _, err := Setup()
	if err != nil {
		panic(err)
	}

	// 初始化日志和kafka
	log.Init()
	if err := kafka.InitProducer([]string{"10.248.4.234:9002"}); err != nil {
		panic(err)
	}

	controller.SetControllerManagerOnce(&controllerManager)
}

func main() {
	group, ctx := errgroup.WithContext(context.Background())

	group.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			close(signalChannel)
			return nil
		case <-signalChannel:
			return errors.WithStack(zerror.ErrSignalDone)
		}
	})
	group.Go(func() error {
		return errors.Wrap(zerror.ErrServerDown, fmt.Sprintf("iris app发生错误, error : %s", controller.NewApp(ctx)))
	})
	group.Go(func() error {
		return errors.Wrap(zerror.ErrServerDown, fmt.Sprintf("prometheus服务发生错误, error : %s", prometheus.NewPrometheus(ctx)))
	})

	if err := group.Wait(); err != nil {
		_ = zerror.HandleError(err)
	}

	time.Sleep(5 * time.Second)
}
