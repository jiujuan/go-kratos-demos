package service

import (
	"context"

	pb "quickstart/api/helloworld/v1"
)

type StudentService struct {
	pb.UnimplementedStudentServer
}

func NewStudentService() *StudentService {
	return &StudentService{}
}

func (s *StudentService) Createstudent(ctx context.Context, req *pb.CreateStudentRequest) (*pb.CreateStudentReply, error) {
	return &pb.CreateStudentReply{}, nil
}
func (s *StudentService) Updatestudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.UpdateStudentReply, error) {
	return &pb.UpdateStudentReply{}, nil
}
func (s *StudentService) Deletestudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentReply, error) {
	return &pb.DeleteStudentReply{}, nil
}
func (s *StudentService) Getstudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.GetStudentReply, error) {
	return &pb.GetStudentReply{}, nil
}
func (s *StudentService) Liststudent(ctx context.Context, req *pb.ListStudentRequest) (*pb.ListStudentReply, error) {
	return &pb.ListStudentReply{}, nil
}
func (s *StudentService) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{}, nil
}
