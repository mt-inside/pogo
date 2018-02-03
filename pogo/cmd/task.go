package cmd

import (
	"github.com/mt-inside/pogo/pogo/task"
	"github.com/spf13/cobra"

	pb "github.com/mt-inside/pogo/proto"
)

func init() {
	rootCmd.AddCommand(taskCommand)
}

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Manages tasks",
	Run: func(cmd *cobra.Command, args []string) {
		task.AddTask(&pb.Task{"get up"})
		task.AddTask(&pb.Task{"go climbing"})
		task.ListTasks()
	},
}
