package studentdex

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"grpc/student-api/pb"
)

type Server struct {
	repository *Repository
	pb.UnimplementedStudentdexServer
}

func NewServer(repository *Repository) *Server {
	return &Server{
		repository: repository,
	}
}

func (s *Server) Create(
	ctx context.Context,
	req *pb.StudentRequest,
) (*pb.Student, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id cannot be 0")
	}

	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name cannot be empty")
	}

	if len(req.Grade) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "grade cannot be empty")
	}

	/*types := make([]string, len(req.Type))
	for i, t := range req.Type {
		types[i] = t.String()
	}*/

	now := time.Now()

	studentDetails := Student{
		ID:        req.Id,
		Name:      req.Name,
		Grade:     req.Grade,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Insert(ctx, studentDetails); err != nil {
		return nil, fmt.Errorf("failed to insert student: %w", err)
	}

	res := studentToResponse(studentDetails)

	return &res, nil
}

func (s *Server) Read(ctx context.Context, req *pb.StudentFilter) (*pb.StudentListResponse, error) {
	student, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all student: %w", err)
	}

	res := make([]*pb.Student, len(student))

	for i, studentDex := range student {
		p := studentToResponse(studentDex)
		res[i] = &p
	}

	return &pb.StudentListResponse{
		Student: res,
	}, nil
}

func (s *Server) ReadOne(ctx context.Context, req *pb.StudentID) (*pb.Student, error) {
	student, err := s.repository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find student")
	}

	res := studentToResponse(student)

	return &res, nil
}

func (s *Server) Update(ctx context.Context, req *pb.StudentUpdateRequest) (*pb.Student, error) {
	p, err := s.repository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find student")
	}

	/*strTypes := make([]string, len(req.Type))
	for i, t := range req.Type {
		strTypes[i] = t.String()
	}*/

	p.UpdatedAt = time.Now()
	p.Name = req.Name
	p.Grade = req.Grade

	err = s.repository.Update(ctx, p)
	if errors.Is(err, ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "failed to update student")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update student")
	}

	res := studentToResponse(p)

	return &res, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.StudentID) (*emptypb.Empty, error) {
	return nil, nil
}
