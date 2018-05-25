package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pb "github.com/mt-inside/pogo/api"
	"github.com/mt-inside/pogo/pkg/pogo/tasks"
)

func init() {
	rootCmd.AddCommand(categoryCommand)
	categoryCommand.AddCommand(categorySetCommand)
	categoryCommand.AddCommand(categoryShowCommand)
	categoryCommand.AddCommand(categoryListCommand)
}

/* This layer should: unmarshall user input, checks & sanitises, send on
* - next layer is responsible for turning into internal types and e.g.
* finding existing object
* There are no interal types atm, so this turns straight into PBs
 */

var categoryCommand = &cobra.Command{
	Use:   "category",
	Short: "Manages categories",
}

var categorySetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set current Category",
	Long:  "TODO",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("category", args[0])
		//err := viper.WriteConfigAs("/home/matt/.pogo/config.yaml") //FIXME
		err := viper.WriteConfig()
		if err != nil {
			log.Println("Couldn't write config file: ", err)
		}
	},
}

var categoryShowCommand = &cobra.Command{
	Use:   "show",
	Short: "Show current Category",
	Long:  "TODO",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("category"))
	},
}

var categoryListCommand = &cobra.Command{
	Use:   "list",
	Short: "List existing Categories",
	Long:  "TODO",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		/* TODO only the non-completed ones? add --completed flag for all of them? */
		/* TODO: service for this on the server so less has to go across the wire */
		cs := map[string]bool{}
		for _, t := range tasks.ListTasks(&pb.TaskFilter{}) { // FIXME factor out
			cs[t.Category] = true
		}
		for c := range cs {
			fmt.Println(c)
		}
	},
}
