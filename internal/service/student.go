package service

import (
	"context"

	pb "student/api/student/v1"
	"student/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type StudentService struct {
	pb.UnimplementedStudentServer

	student *biz.StudentUsecase
	log     *log.Helper
}

func NewStudentService(stu *biz.StudentUsecase, logger log.Logger) *StudentService {
	return &StudentService{
		student: stu,
		log:     log.NewHelper(logger),
	}
}

// func (s *StudentService) Createstudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
// 	return &pb.CreateStudentReply{}, nil
// }
// func (s *StudentService) Updatestudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
// 	return &pb.UpdateStudentReply{}, nil
// }
// func (s *StudentService) Deletestudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
// 	return &pb.DeleteStudentReply{}, nil
// }
// func (s *StudentService) Liststudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentReply, error) {
// 	return &pb.ListStudentReply{}, nil
// }
// 获取学生信息
func (s *StudentService) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	stu, err := s.student.Get(ctx, req.Id)

	if err != nil {
		return nil, err
	}
	return &pb.GetStudentReply{
		Id:     stu.ID,
		Status: stu.Status,
		Name:   stu.Name,
	}, nil
}
