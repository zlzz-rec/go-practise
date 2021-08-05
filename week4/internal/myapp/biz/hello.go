package biz

import "week4/internal/myapp/data"

type HelloBiz struct {
	helloRepo data.HelloRepo
}

func NewHelloBiz(helloRepo data.HelloRepo) HelloBiz {
	return HelloBiz{helloRepo: helloRepo}
}

type HelloUsecase interface {
	SayHello() string
}

func (b *HelloBiz) SayHello() string {
	return b.helloRepo.SayHello()
}
