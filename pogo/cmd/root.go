package cmd

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	/* Config sources, in viper presidence order */

	/* Defaults */
	viper.SetDefault("category", "default")

	/* Config file */
	viper.SetConfigName("config") /* e.g. config.yaml, config.json */
	viper.AddConfigPath("$HOME/.pogo/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		//Also goes off for no config file
		//log.Fatalf("Fatal error in config file: %s \n", err)
	}
	/* ...with auto-reload */
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		viper.Debug()
	})

	/* Environment */
	viper.SetEnvPrefix("pogo")
	viper.AutomaticEnv() /* e.g. POGO_CATEGORY */

	/* Command-line args */
	// category flag is under taskCommand
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
			fmt.Println("Category: ", viper.GetString("category"))
			listCommand.Run(listCommand, make([]string, 0))
		} else {
			if s.State == pb.Status_TASK {
				fmt.Println("    RUNNING")
			} else if s.State == pb.Status_BREAK {
				fmt.Println("    BREAK")
			}
			// TODO: really does need to be an interal type becuase it needs to be printable
			fmt.Println("Category: ", viper.GetString("category"))
			fmt.Println(RenderTask(s.Task))
			fmt.Println("Time remaining: ", s.RemainingTime)
		}
	},
}
