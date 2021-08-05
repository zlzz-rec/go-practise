package myapp

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"week4/internal/myapp/controller"
)

type AppControllers struct {
	HelloController controller.HelloController
}

func NewAppControllers(HelloController controller.HelloController) AppControllers {
	return AppControllers{HelloController: HelloController}
}

func NewApp(port int, appControllers AppControllers, stop chan struct{}) error {
	app := iris.New()
	app.Use(recover.New())

	app.Get("hello", appControllers.HelloController.SayHello)

	go func() {
		<-stop
		_ = app.Shutdown(context.Background())
	}()
	return app.Run(iris.Addr(fmt.Sprintf(":%v", port)), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
