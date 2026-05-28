package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		birthday, err := time.Parse("2006-01-02", args[1])
		if err != nil {
			fmt.Println("Not a correct date, example 2006-01-30")
			return
		}

		err = personService.Add(args[0], birthday)
		if err != nil {
			fmt.Println("Error: ")
			fmt.Println(err)
			return
		}

		fmt.Printf("Person %s added\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
