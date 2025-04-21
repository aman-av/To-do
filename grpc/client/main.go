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
		Title:       "Learn gRPC",
		Description: "Build a ToDo App",
	})
	log.Printf("Created Task: %v", task)

	// List Tasks
	taskList, _ := client.ListTasks(ctx, &pb.Empty{})
	log.Println("Task List:")
	for _, t := range taskList.Tasks {
		log.Printf("- %v", t)
	}

	// Update Task
	task, _ = client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:          "dcb111e2-ec03-4fcf-a024-8bf2eed22536",
		Title:       "Update gRPC",
		Description: "Update a ToDo App",
	})
	log.Printf("Update Task: %v", task)

	// List Tasks
	taskList, _ = client.ListTasks(ctx, &pb.Empty{})
	log.Println("Task List:")
	for _, t := range taskList.Tasks {
		log.Printf("- %v", t)
	}

	// Delete Task
	client.DeleteTask(ctx, &pb.DeleteTaskRequest{
		Id: "d4793116-5f33-419c-bce3-ef6c37f5e520",
	})

	// List Tasks
	taskList, _ = client.ListTasks(ctx, &pb.Empty{})
	log.Println("Task List:")
	for _, t := range taskList.Tasks {
		log.Printf("- %v", t)
	}

}
