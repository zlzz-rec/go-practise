package data

import "fmt"

type HelloRepo struct {
	data Data
}

func NewHelloRepo(data Data) HelloRepo {
	return HelloRepo{data: data}
}

func (repo *HelloRepo) SayHello() string {
	repo.data.Redis.Set("hello", "halo")
	res, _ := repo.data.Redis.Get("hello")
	fmt.Println(res)
	return repo.data.Orm.Name()
}
