package cmd

import (
	"age/internal/service"
	"age/internal/ui"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var personService *service.PersonService

var rootCmd = &cobra.Command{
	Use:           "age [birthday]",
	Short:         "Keep track of birthdays without doing mental math",
	Long:          "Age is a small CLI for saving birthdays, checking current ages, and seeing how long remains until the next celebration.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		persons, err := personService.List()
		if err != nil {
			ui.PrintError(fmt.Sprintf("could not load birthdays: %v", err))
			return
		}

		fmt.Println(ui.RenderDashboard(persons))
	},
}

func Execute(serv *service.PersonService) {
	personService = serv

	if tryBirthdayPreview(os.Args[1:]) {
		return
	}

	if err := rootCmd.Execute(); err != nil {
		ui.PrintError(fmt.Sprintf("command failed: %v", err))
		os.Exit(1)
	}
}

func tryBirthdayPreview(args []string) bool {
	if len(args) != 1 {
		return false
	}

	candidate := strings.TrimSpace(args[0])
	if candidate == "" || strings.HasPrefix(candidate, "-") || isSubcommand(candidate) {
		return false
	}

	birthday, err := parseBirthday(candidate)
	if err != nil {
		return false
	}

	fmt.Println(ui.RenderBirthdayPreview(personService.PreviewBirthday(birthday)))
	return true
}

func isSubcommand(name string) bool {
	for _, command := range rootCmd.Commands() {
		if command.Name() == name {
			return true
		}

		for _, alias := range command.Aliases {
			if alias == name {
				return true
			}
		}
	}

	return false
}
