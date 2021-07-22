package errservice

import "github.com/pkg/errors"

type TeacherService struct {
}

func (s *TeacherService)Query(id int) (*Teacher, error){
	repo := TeacherRepo{}
	teacher, err := repo.Query(id)
	if err != nil {
		return nil, errors.WithMessagef(err, "额外信息, %d\n", 123456)
	}
	return teacher, nil
}