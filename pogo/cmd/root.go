package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/mt-inside/pogo/pogo/tasks"
	pb "github.com/mt-inside/pogo/proto"
)

/* Avoids exporting rootCmd */
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Can't parse args: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	//log.Printf("Configured cobra")
}

var rootCmd = &cobra.Command{
	Use:   "pogo",
	Short: "pogo is a pomodoro timer",
	Long: `A client-server pomodoro timer system. Pogo supports multiple
	notification mechanisms and can produce timesheets`,
	Run: func(cmd *cobra.Command, args []string) {
		s := tasks.GetStatus()

		/* TODO: use some library to template this into a txt snippet?
		 * Add colour? */
		if s.State == pb.Status_IDLE {
			fmt.Println("    IDLE")
			listCommand.Run(listCommand, make([]string, 0))
		} else {
			if s.State == pb.Status_TASK {
				fmt.Println("    RUNNING")
			} else if s.State == pb.Status_BREAK {
				fmt.Println("    BREAK")
			}
			// TODO: really does need to be an interal type becuase it needs to be printable
			fmt.Printf("%d: %s [%s]\n", s.Task.Id.Idx, s.Task.Title, s.Task.State)
			fmt.Printf("Time remaining: %d\n", s.RemainingTime)
		}
	},
}
