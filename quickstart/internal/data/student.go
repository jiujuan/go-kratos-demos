package data

import (
	"context"
	"quickstart/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type studentRepo struct {
	data *Data
	log  *log.Helper
}

func NewStudentRepo(data *Data, logger log.Logger) biz.StudentRepo {
	return &studentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *studentRepo) Save(ctx context.Context, stu *biz.Student) (*biz.Student, error) {
	return stu, nil
}

func (repo *studentRepo) Get(ctx context.Context, stu *biz.Student) (*biz.Student, error) {
	return stu, nil
}
