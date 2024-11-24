/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/rbieldt/CLI-Todo/go/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	port int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the gRPC server",
	Long: `Start the gRPC server for the Todo application.
This server will handle task management operations including:
- Adding new tasks
- Listing existing tasks
- Marking tasks as completed`,
	Run: func(cmd *cobra.Command, args []string) {
		startGRPCServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Add flags for server configuration
	serverCmd.Flags().IntVarP(&port, "port", "p", 50051, "Port to run the server on")
}

type Task struct {
	Id          int64
	Description string
	Completed   bool
}

type server struct {
	pb.UnimplementedTodoServiceServer
	tasks      []Task
	nextTaskId int64
}

func newServer() *server {
	return &server{
		tasks:      make([]Task, 0),
		nextTaskId: 1,
	}
}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	task := Task{
		Id:          s.nextTaskId,
		Description: in.GetDescription(),
		Completed:   false,
	}
	s.tasks = append(s.tasks, task)
	s.nextTaskId++

	log.Printf("Added task: %v", task)
	return &pb.AddTaskResponse{
		Task: &pb.Task{
			Id:          task.Id,
			Description: task.Description,
			Completed:   task.Completed,
		},
	}, nil
}

func (s *server) GetTasks(ctx context.Context, in *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	var pbTasks []*pb.Task
	for _, task := range s.tasks {
		if in.GetIncludeCompleted() || !task.Completed {
			pbTasks = append(pbTasks, &pb.Task{
				Id:          task.Id,
				Description: task.Description,
				Completed:   task.Completed,
			})
		}
	}

	log.Printf("Returning %d tasks", len(pbTasks))
	return &pb.GetTasksResponse{
		Tasks: pbTasks,
	}, nil
}

func (s *server) CompleteTask(ctx context.Context, in *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	taskId := in.GetTaskId()
	for i := range s.tasks {
		if s.tasks[i].Id == taskId {
			s.tasks[i].Completed = true
			log.Printf("Marked task %d as completed", taskId)
			return &pb.CompleteTaskResponse{
				Task: &pb.Task{
					Id:          s.tasks[i].Id,
					Description: s.tasks[i].Description,
					Completed:   s.tasks[i].Completed,
				},
			}, nil
		}
	}

	log.Printf("Task %d not found", taskId)
	return nil, fmt.Errorf("task with ID %d not found", taskId)
}

func startGRPCServer() {
	address := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	srv := newServer()
	pb.RegisterTodoServiceServer(s, srv)

	log.Printf("Server starting on port %d...", port)
	log.Printf("Ready to handle requests!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
