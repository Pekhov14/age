package cmd

import (
	"age/internal/ui"
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show saved birthdays",
	Long:  "Display all saved people with their birth date, age, and countdown to the next birthday.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		persons, err := personService.List()
		if err != nil {
			ui.PrintError(fmt.Sprintf("could not load birthdays: %v", err))
			return
		}

		fmt.Println(ui.RenderList(persons))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
