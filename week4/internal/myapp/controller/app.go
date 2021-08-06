package controller

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
)

func NewApp(port int, stop chan struct{}) error {
	app := iris.New()
	app.Use(recover.New())

	app.Get("hello", GetControllerManager().HelloController.SayHello)

	go func() {
		<-stop
		_ = app.Shutdown(context.Background())
	}()
	return app.Run(iris.Addr(fmt.Sprintf(":%v", port)), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
