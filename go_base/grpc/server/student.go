package main

import (
	"context"
	"fmt"
	grpc_model "go/base/grpc/idl/model"
	grpc_service "go/base/grpc/idl/service"
	"io"
)

type Student struct {
	grpc_service.UnimplementedStudentServer
}

func (s Student) QueryStudent(ctx context.Context, query *grpc_service.QueryStudentRequest) (resp *grpc_service.QueryStudentResponse, err error) {
	fmt.Printf("request: %+v\n", query)
	resp = &grpc_service.QueryStudentResponse{
		Students: []*grpc_model.Student{
			{Id: 123, Name: "昝高伟", Age: 18},
			{Id: 456, Name: "zgw", Age: 28},
		},
	}
	return
}

func (s Student) QueryStudents1(ctx context.Context, query *grpc_service.StudentIds) (resp *grpc_service.QueryStudentResponse, err error) {
	fmt.Printf("request: %+v\n", query)
	resp = &grpc_service.QueryStudentResponse{
		Students: []*grpc_model.Student{
			{Id: 123, Name: "zgw", Age: 18},
			{Id: 456, Name: "zzz", Age: 28},
		},
	}
	return
}

// Server streaming RPC
func (s Student) QueryStudents2(query *grpc_service.StudentIds, server grpc_service.Student_QueryStudents2Server) error {
	for i := 0; i < 2; i++ {
		id := int64(i) + 100
		stu := &grpc_model.Student{Id: id, Name: "zgw", Age: 18}
		err := server.Send(stu) //向流中发送一个结果
		if err != nil {
			fmt.Printf("send Student %d failed: %s\n", id, err)
			return err
		}
	}
	return nil
}

// Client streaming RPC
func (s Student) QueryStudents3(server grpc_service.Student_QueryStudents3Server) error {
	students := make([]*grpc_model.Student, 0, 10)
	for {
		sid, err := server.Recv() //从流中取出一个结果
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("recv request3 failed: %s\n", err)
			continue
		}
		stu := &grpc_model.Student{Id: sid.Id, Name: "zzz", Age: 18}
		students = append(students, stu)
	}
	return server.SendMsg(&grpc_service.QueryStudentResponse{Students: students})
}

// Bidirectional streaming RPC
func (s Student) QueryStudents4(server grpc_service.Student_QueryStudents4Server) error {
	for {
		sid, err := server.Recv() //从流中取出一个结果
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("recv request3 failed: %s\n", err)
			continue
		}
		stu := &grpc_model.Student{Id: sid.Id, Name: "zzz", Age: 18}
		err = server.Send(stu) //向流中发送一个结果
		if err != nil {
			fmt.Printf("send Student %d failed: %s\n", stu.Id, err)
			return err
		}
	}
	return nil
}
