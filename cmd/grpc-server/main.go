package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "example.com/rest-grpc-todo-demo-go/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	mu     sync.Mutex
	tasks  []*pb.Task
	nextID int32 = 1
)

type taskServer struct {
	pb.UnimplementedTaskServiceServer
}

func (s *taskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	return &pb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *taskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, t := range tasks {
		if t.Id == req.Id {
			return &pb.GetTaskResponse{
				Task: t,
			}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "task %d not found", req.Id)
}

func (s *taskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range tasks {
		if t.Id == req.Id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return &pb.DeleteTaskResponse{}, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "task %d not found", req.Id)
}

func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	mu.Lock()
	defer mu.Unlock()

	t := &pb.Task{
		Id:    nextID,
		Title: req.Title,
		Done:  false,
	}
	nextID++
	tasks = append(tasks, t)

	return t, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &taskServer{})

	log.Println("gRPC server listening on localhost:50051 ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
