package controller

import (
	"github.com/kataras/iris"
	"week4/internal/myapp/service"
)

type HelloController struct {
	helloService service.HelloService
}

func NewHelloController(helloService service.HelloService) HelloController{
	return HelloController{helloService: helloService}
}

func (c *HelloController)SayHello(ctx iris.Context) {
	hello := c.helloService.SayHello()
	_, _ = ctx.JSON(iris.Map{
		"code":    0,
		"message": "success",
		"data":    hello,
	})
}