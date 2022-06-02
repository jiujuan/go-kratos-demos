package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Student is a Student model.
type Student struct {
	ID        int32
	Name      string
	Info      string
	Status    int32
	UpdatedAt time.Time
	CreatedAt time.Time
}

// 定义 Student 的操作接口
type StudentRepo interface {
	GetStudent(context.Context, int32) (*Student, error) // 根据 id 获取学生信息
}

type StudentUsecase struct {
	repo StudentRepo
	log  *log.Helper
}

// 初始化 StudentUsecase
func NewStudentUsecase(repo StudentRepo, logger log.Logger) *StudentUsecase {
	return &StudentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 通过 id 获取 student 信息
func (uc *StudentUsecase) Get(ctx context.Context, id int32) (*Student, error) {
	uc.log.WithContext(ctx).Infof("biz.Get: %d", id)
	return uc.repo.GetStudent(ctx, id)
}
