package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	serverDownSignal := make(chan struct{})
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/serverDown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("发送server down信号")
		serverDownSignal <- struct{}{}
	})
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8081",
	}

	// 启动httpserver监听
	group.Go(func() error {
		return server.ListenAndServe()
	})

	// 接受linux信号
	group.Go(func() error {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, os.Interrupt, os.Kill)

		select {
		case <-ctx.Done():
			fmt.Println("监听linux信号协程监听到context结束")
		case shutdownSignal := <-signalChannel:
			fmt.Printf("收到linux关闭信号, 关闭http server, 信号为:%v\n", shutdownSignal)
		}
		return errors.Errorf("监听linux信号协程结束")
	})

	// 收到退出信号, 退出服务
	group.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("监听http信号协程监听到context结束")
		case <-serverDownSignal:
			fmt.Println("收到http关闭信号, 关闭http server")
		}
		_ = server.Shutdown(ctx)
		return errors.Errorf("监听http信号协程结束")
	})

	// 因为waitgroup中的error是sync.once, 所以存储的为先赋值的error
	if err := group.Wait(); err != nil {
		fmt.Printf("errorgroup出现关闭信号, error:%+v\n", err)
	}
}
