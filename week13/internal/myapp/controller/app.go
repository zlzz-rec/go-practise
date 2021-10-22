package controller

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

func NewApp(ctx context.Context) error {
	app := iris.New()
	app.Use(recover.New())

	app.Get("hello", GetControllerManager().HelloController.SayHello)

	go func() {
		select {
		case <-ctx.Done():
			_ = app.Shutdown(ctx)
		}
	}()

	return app.Run(iris.Addr(fmt.Sprintf(":%v", 8080)), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
