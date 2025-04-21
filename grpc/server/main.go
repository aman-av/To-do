package main

import (
	"context"
	"fmt"
	pb "github.com/aman-av/grpc/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	tasks map[string]*pb.Task
}

func (s *server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	id := uuid.New().String()
	task := &pb.Task{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
	}
	s.tasks[id] = task
	return task, nil
}

func (s *server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	task, exists := s.tasks[req.Id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (s *server) ListTasks(_ *pb.Empty, stream pb.TodoService_ListTasksServer) error {
	for _, task := range s.tasks {
		// Simulate a delay
		time.Sleep(1 * time.Second)
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	task, exists := s.tasks[req.Id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	task.Title = req.Title
	task.Description = req.Description
	return task, nil
}

func (s *server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.Empty, error) {
	delete(s.tasks, req.Id)
	return &pb.Empty{}, nil
}

func (s *server) CreateTasks(stream pb.TodoService_CreateTasksServer) error {
	var createdCount int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Successfully received all tasks; now respond
			return stream.SendAndClose(&pb.CreateTasksResponse{
				CreatedCount: createdCount,
			})
		}
		if err != nil {
			break
		}
		id := uuid.New().String()
		task := &pb.Task{
			Id:          id,
			Title:       req.Title,
			Description: req.Description,
		}
		s.tasks[id] = task
		createdCount++
	}
	return nil
}

func (s *server) TaskChat(stream pb.TodoService_TaskChatServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		task, exists := s.tasks[req.Id]
		if !exists {
			return fmt.Errorf("task not found")
		}
		
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, &server{tasks: make(map[string]*pb.Task)})
	log.Println("Server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
