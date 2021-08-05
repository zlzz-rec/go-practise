package data

type HelloRepo struct {
	data Data
}

func NewHelloRepo(data Data) HelloRepo {
	return HelloRepo{data: data}
}

func (repo *HelloRepo) SayHello() string {
	return repo.data.Orm.Name()
}
