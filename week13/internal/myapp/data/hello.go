package data

import (
	"fmt"
	"github.com/pkg/errors"
	"week13/pkg/zerror"
)

type HelloRepo struct {
	data Data
}

func NewHelloRepo(data Data) HelloRepo {
	return HelloRepo{data: data}
}

func (repo *HelloRepo) SayHello() (string, error) {
	repo.data.Redis.Set("hello", "halo")
	res, err := repo.data.Redis.Get("hello")
	if err != nil {
		return "", errors.Wrap(zerror.ErrInnerServer, fmt.Sprintf("redis get失败, key=%s", "hello"))
	}
	fmt.Println(res)

	return repo.data.Orm.Name(), nil
}
