package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mt-inside/pogo/pogo/tasks"
	pb "github.com/mt-inside/pogo/proto"
)

func init() {
	rootCmd.AddCommand(taskCommand)
	taskCommand.AddCommand(addCommand)
	taskCommand.AddCommand(listCommand)
	taskCommand.AddCommand(startCommand)
	taskCommand.AddCommand(stopCommand)
	taskCommand.AddCommand(completeCommand)

	taskCommand.PersistentFlags().StringP("category", "c", "default", "Category that the task is in")
	viper.BindPFlag("category", taskCommand.PersistentFlags().Lookup("category"))
	taskCommand.PersistentFlags().BoolP("all-categories", "", false, "Work with Tasks across all categories")
	viper.BindPFlag("all-categories", taskCommand.PersistentFlags().Lookup("all-categories"))
	taskCommand.PersistentFlags().BoolP("completed", "", false, "Show completed Tasks")
	viper.BindPFlag("completed", taskCommand.PersistentFlags().Lookup("completed"))
}

/* This layer should: unmarshall user input, checks & sanitises, send on
* - next layer is responsible for turning into internal types and e.g.
* finding existing object
* There are no interal types atm, so this turns straight into PBs
 */

func NewFilter() *pb.TaskFilter {
	return &pb.TaskFilter{
		Task:   &pb.Task{},
		Fields: 0,
	}
}
func SetCategory(f *pb.TaskFilter, category string) *pb.TaskFilter {
	f.Task.Category = category
	f.Fields = f.Fields | pb.TaskFields_category
	return f
}
func AddState(f *pb.TaskFilter, state pb.TaskState) *pb.TaskFilter {
	f.Task.State = f.Task.State | state
	f.Fields = f.Fields | pb.TaskFields_state
	return f
}
func AddType(f *pb.TaskFilter, typ pb.TaskType) *pb.TaskFilter {
	f.Task.Type = f.Task.Type | typ
	f.Fields = f.Fields | pb.TaskFields_type
	return f
}

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
		/* Id left to default and ignored at the other end */
		tasks.AddTask(
			&pb.Task{
				Title: strings.Join(args, " "),
				// TODO gets the command-line arg too? Should...
				Category: viper.GetString("category"),
			},
		)
	},
}

/* TODO: move me */
func RenderTask(t *pb.Task) string {
	return fmt.Sprintf("%d: [%s] %s - %s", t.Id.Idx, t.Category, t.Title, t.State)
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List Tasks",
	Long:  "TODO",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		/* TODO: make the empty filter, have everything be an addative
		 * WithFoo() */
		filter := NewFilter()

		if !viper.GetBool("all-categories") {
			filter = SetCategory(filter, viper.GetString("category"))
		}

		filter = AddState(filter, pb.TaskState_TODO)
		if viper.GetBool("completed") {
			filter = AddState(filter, pb.TaskState_DONE)
		}

		filter = AddType(filter, pb.TaskType_TASK)

		for _, t := range tasks.ListTasks(filter) {
			fmt.Println(RenderTask(t))
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
	Short: "Start a pomodoro, working on the specified task",
	Long:  "TODO",
	Args:  validTaskId,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := strconv.ParseInt(args[0], 10, 64)
		tasks.StartTask(id)
	},
}

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop working on current task",
	Long:  "TODO",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		tasks.StopTask()
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
