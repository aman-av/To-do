package main

import (
	"context"
	"log"
	"time"

	pb "github.com/aman-av/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTodoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//Create Task
	task, _ := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       "Learn gRPC1",
		Description: "Build a ToDo App",
	})
	log.Printf("Created Task: %v", task)

	task, _ = client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       "Learn gRPC2",
		Description: "Build a ToDo App",
	})
	log.Printf("Created Task: %v", task)

	task, _ = client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       "Learn gRPC3",
		Description: "Build a ToDo App",
	})
	log.Printf("Created Task: %v", task)

	// List Tasks
	stream, err := client.ListTasks(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Error listing tasks: %v", err)
	}
	for {
		task, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Task: %v", task)
	}
}
