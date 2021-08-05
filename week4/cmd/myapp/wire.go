// +build wireinject
// 当前文件为独立文件, 最终不会被build到工程代码中
// 在make generate后wire组件会根据绑定关系自动生成wire_gen.go文件. 请不要编辑wire_gen.go文件

package main

import (
	"github.com/google/wire"
	"week4/internal/myapp"
	"week4/internal/myapp/biz"
	"week4/internal/myapp/controller"
	"week4/internal/myapp/data"
	"week4/internal/myapp/service"
)

func initApp() (myapp.AppControllers, func(), error) {
	panic(wire.Build(myapp.NewAppControllers, controller.ProviderSet, service.ProviderSet, biz.ProviderSet, data.ProviderSet))
}
