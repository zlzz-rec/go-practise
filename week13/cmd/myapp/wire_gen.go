// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"week13/internal/myapp/biz"
	"week13/internal/myapp/controller"
	"week13/internal/myapp/data"
	"week13/internal/myapp/service"
)

// Injectors from wire.go:

func Setup() (controller.AllControllers, func(), error) {
	dataData, cleanup, err := data.NewData()
	if err != nil {
		return controller.AllControllers{}, nil, err
	}
	helloRepo := data.NewHelloRepo(dataData)
	helloBiz := biz.NewHelloBiz(helloRepo)
	helloService := service.NewHelloService(helloBiz)
	helloController := controller.NewHelloController(helloService)
	allControllers := controller.NewControllerManager(helloController)
	return allControllers, func() {
		cleanup()
	}, nil
}
