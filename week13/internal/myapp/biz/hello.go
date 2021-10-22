package biz

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.feedtoken.tech/ztouch/zglib/ztime"
	"week13/internal/myapp/data"
	"week13/pkg/kafka"
	"week13/pkg/zerror"
)

type HelloBiz struct {
	helloRepo data.HelloRepo
}

func NewHelloBiz(helloRepo data.HelloRepo) HelloBiz {
	return HelloBiz{helloRepo: helloRepo}
}

type HelloUsecase interface {
	SayHello() string
}

func (b *HelloBiz) SayHello() (string, error) {
	name, err := b.helloRepo.SayHello()
	if err != nil {
		return "", err
	}
	if err = b.sendHello(name); err != nil {
		return name, errors.Wrap(zerror.ErrInnerServer, fmt.Sprintf("kafka发送消息失败, key:%s", name))
	}
	return name, nil
}

func (b *HelloBiz) sendHello(name string) error {
	type HelloMsg struct {
		Name  string
		Extra string
	}
	message, err := json.Marshal(HelloMsg{
		Name: name,
	})
	if err != nil {
		return err
	}

	if err = kafka.Producer.Produce("SayHelloTopic", ztime.GetNowTimeUnixStr(), string(message)); err != nil {
		return err
	}
	return nil
}
