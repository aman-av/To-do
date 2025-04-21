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

	stream, err := client.CreateTasks(ctx)
	if err != nil {
		log.Fatalf("Error creating tasks: %v", err)
	}

	tasksToSend := []pb.CreateTaskRequest{
		{Title: "Task A", Description: "Bulk A"},
		{Title: "Task B", Description: "Bulk B"},
		{Title: "Task C", Description: "Bulk C"},
	}

	for _, t := range tasksToSend {
		if err := stream.Send(&t); err != nil {
			log.Fatalf("error sending task: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error closing stream: %v", err)
	}
	log.Printf("Bulk Created Tasks: %d", res.CreatedCount)

	// List Tasks
	stream2, err := client.ListTasks(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Error listing tasks: %v", err)
	}
	for {
		task, err := stream2.Recv()
		if err != nil {
			break
		}
		log.Printf("Task: %v", task)
	}
}
