package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		persons, _ := personService.List()
		if len(persons) == 0 {
			fmt.Println("No persons found")
		}

		fmt.Printf("%-10s %-15s %-4s %-10s\n", "Name", "Birthday", "Age", "Next Birthday")
		fmt.Println(strings.Repeat("-", 50))

		for _, person := range persons {
			fmt.Printf("%-10s %-15s %-4d %-10s \n",
				person.Name,
				person.Birthday,
				person.Age,
				person.DaysUntilBD,
			)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
