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
		s := tasks.State()

		/* TODO: use some library to template this into a txt snippet?
		 * Add colour? */
		if s.State == pb.PogoState_IDLE {
			fmt.Println("    IDLE")
			listCommand.Run(listCommand, make([]string, 0))
		} else {
			if s.State == pb.PogoState_TASK {
				fmt.Println("    RUNNING")
			} else if s.State == pb.PogoState_BREAK {
				fmt.Println("    BREAK")
			}
			fmt.Println()
			fmt.Println()
			fmt.Printf("        %v\n", s.Task)
			fmt.Println()
			fmt.Printf("        Time remaining: %d\n", s.RemainingTime)
		}
	},
}
