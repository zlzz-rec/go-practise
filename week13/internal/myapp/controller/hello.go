package controller

import (
	"github.com/kataras/iris"
	"week13/internal/myapp/service"
	"week13/pkg/Response"
)

type HelloController struct {
	helloService service.HelloService
}

func NewHelloController(helloService service.HelloService) HelloController{
	return HelloController{helloService: helloService}
}

func (c *HelloController)SayHello(ctx iris.Context) {
	name, err := c.helloService.SayHello()
	if err != nil {
		Response.ErrResponse(ctx, err)
	}
	Response.SuccResponse(ctx, name, "")
}