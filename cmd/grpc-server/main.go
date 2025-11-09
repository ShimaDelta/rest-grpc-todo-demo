package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "example.com/rest-grpc-todo-demo-go/pb/proto"
	"google.golang.org/grpc"
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
