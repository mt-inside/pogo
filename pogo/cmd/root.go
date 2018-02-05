package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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
		fmt.Println("TODO: quick status here")
	},
}
