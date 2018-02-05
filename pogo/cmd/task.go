package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mt-inside/pogo/pogo/tasks"
	"github.com/spf13/cobra"

	pb "github.com/mt-inside/pogo/proto"
)

func init() {
	rootCmd.AddCommand(taskCommand)
	taskCommand.AddCommand(addCommand)
	taskCommand.AddCommand(listCommand)
	taskCommand.AddCommand(startCommand)
	taskCommand.AddCommand(completeCommand)
}

/* This layer should: unmarshall user input, checks & sanitises, send on
* - next layer is responsible for turning into internal types and e.g.
* finding existing object
* There are no interal types atm, so this turns straight into PBs
 */

var taskCommand = &cobra.Command{
	Use:   "task",
	Short: "Manages tasks",
}

var addCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a Task",
	Long:  "TODO",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tasks.AddTask(&pb.ProtoTask{Title: strings.Join(args, " ")})
	},
}

// TODO: add --completed
var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List Tasks",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		for _, t := range tasks.ListTasks() {
			fmt.Printf("%d: %s\n", t.Id.Idx, t.Title)
		}
	},
}

func validTaskId(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("requires command ID")
	}
	if _, err := strconv.Atoi(args[0]); err != nil {
		return fmt.Errorf("invalid id: %s", args[0])
	}
	return nil
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Start a pomodoro, working on Task",
	Long:  "TODO",
	Args:  validTaskId,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.ParseInt(args[0], 10, 64)
		tasks.StartTask(id)
	},
}

var completeCommand = &cobra.Command{
	Use:   "complete",
	Short: "Mark a Task as complete",
	Long:  "TODO",
	Args:  validTaskId,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.ParseInt(args[0], 10, 64)
		tasks.CompleteTask(id)
	},
}
