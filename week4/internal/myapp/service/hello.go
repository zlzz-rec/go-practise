package service

import "week4/internal/myapp/biz"

type HelloService struct {
	helloBiz biz.HelloBiz
}

func NewHelloService(helloBiz biz.HelloBiz) HelloService {
	return HelloService{helloBiz: helloBiz}
}

func (s *HelloService) SayHello() string {
	return s.helloBiz.SayHello()
}
