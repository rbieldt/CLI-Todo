/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/rbieldt/CLI-Todo/go/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddress string
	description   string
	taskID        int64
	showAll       bool
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interact with the Todo server",
	Long: `A client to interact with the Todo gRPC server.
Supports operations for managing tasks including:
- Adding new tasks
- Listing tasks
- Marking tasks as completed`,
}

// addCmd represents the add subcommand
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if description == "" && len(args) > 0 {
			description = strings.Join(args, " ")
		}
		if description == "" {
			fmt.Println("Error: Task description is required")
			os.Exit(1)
		}
		client := getClient()
		addTask(client, description)
	},
}

// listCmd represents the list subcommand
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		listTasks(client, showAll)
	},
}

// completeCmd represents the complete subcommand
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		completeTask(client, taskID)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Add subcommands
	clientCmd.AddCommand(addCmd)
	clientCmd.AddCommand(listCmd)
	clientCmd.AddCommand(completeCmd)

	// Add flags
	clientCmd.PersistentFlags().StringVarP(&serverAddress, "server", "s", "localhost:50051", "Server address")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Task description")
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show completed tasks")
	completeCmd.Flags().Int64VarP(&taskID, "id", "i", 0, "Task ID to complete")
}

func getClient() pb.TodoServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewTodoServiceClient(conn)
}

func addTask(client pb.TodoServiceClient, description string) {
	resp, err := client.AddTask(context.Background(), &pb.AddTaskRequest{
		Description: description,
	})
	if err != nil {
		log.Fatalf("could not add task: %v", err)
	}
	fmt.Printf("Added task with ID: %d\n", resp.Task.Id)
}

func listTasks(client pb.TodoServiceClient, showAll bool) {
	resp, err := client.GetTasks(context.Background(), &pb.GetTasksRequest{
		IncludeCompleted: showAll,
	})
	if err != nil {
		log.Fatalf("could not get tasks: %v", err)
	}

	if len(resp.Tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, task := range resp.Tasks {
		status := "[ ]"
		if task.Completed {
			status = "[✓]"
		}
		fmt.Printf("%s %d: %s\n", status, task.Id, task.Description)
	}
}

func completeTask(client pb.TodoServiceClient, taskID int64) {
	_, err := client.CompleteTask(context.Background(), &pb.CompleteTaskRequest{
		TaskId: taskID,
	})
	if err != nil {
		log.Fatalf("could not complete task: %v", err)
	}
	fmt.Printf("Marked task %d as completed\n", taskID)
}
