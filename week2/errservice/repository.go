package errservice

import (
	"errors"
	errors2 "github.com/pkg/errors"
)

type Teacher struct {
	Id   int
	Name string
	Age  int
}

type TeacherRepo struct {
}

var ErrMysqlQuery = errors.New("mysql query error")
var ErrUnknown = errors.New("unknown error")


func (repo *TeacherRepo) Query(id int) (*Teacher, error) {
	if id < 0 {
		return nil, errors2.Wrapf(ErrMysqlQuery, "查询老师信息失败, 老师id:%d\n", id)
	}

	teacher := &Teacher{
		Id:   id,
		Name: "Jian Mao",
		Age:  33,
	}
	return teacher, nil
}
