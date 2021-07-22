package errservice

import (
	"fmt"
	"github.com/pkg/errors"
)

type TeacherController struct {
}

func (c *TeacherController) Query(id int) {
	service := TeacherService{}
	//teacher, err := service.Query(10)
	teacher, err := service.Query(id)
	if err != nil {
		switch errors.Cause(err) {
		case ErrMysqlQuery:
			fmt.Printf("错误类型:%s\n", ErrMysqlQuery.Error())
			fmt.Printf("错误根因:%+v\n", errors.Cause(err))
			fmt.Printf("错误堆栈:%+v\n", err)
		default:
			fmt.Printf("错误类型:%s\n", ErrUnknown.Error())
		}
		return
	}

	fmt.Printf("老师信息为:%+v\n", teacher)
}
