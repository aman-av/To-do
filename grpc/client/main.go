package main

import (
	"context"
	"io"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
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

	taskIDs := make([]string, 0)
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
		taskIDs = append(taskIDs, task.Id)
		log.Printf("Task: %v", task)
	}

	stream3, err := client.TaskChat(ctx)
	if err != nil {
		log.Fatalf("could not open stream: %v", err)
	}

	// Use a goroutine to send IDs
	go func() {
		for _, id := range taskIDs {
			err := stream3.Send(&pb.GetTaskRequest{Id: id})
			if err != nil {
				log.Fatalf("error sending ID: %v", err)
			}
			time.Sleep(1 * time.Millisecond) // simulate staggered sending
		}
		stream3.CloseSend()
	}()

	// Receive responses
	for {
		task, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving task: %v", err)
		}
		log.Printf("Got Task: %v", task)
	}

}
