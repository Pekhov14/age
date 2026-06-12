package cmd

import (
	"age/internal/ui"
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [birthday]",
	Short: "Add a person and their birthday",
	Long:  "Add a birthday reminder. You can pass name and date as arguments or launch an interactive form.",
	Example: strings.TrimSpace(`
  age add Anton 1999-07-21
  age add Anton
  age add
`),
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name, birthday, ok := collectPersonForAdd(args)
		if !ok {
			return
		}

		if err := personService.Add(name, birthday); err != nil {
			ui.PrintError(fmt.Sprintf("could not save %s: %v", name, err))
			return
		}

		fmt.Println(ui.RenderAddSuccess(name, birthday))
	},
}

func collectPersonForAdd(args []string) (string, time.Time, bool) {
	if len(args) == 2 {
		birthday, err := parseBirthday(args[1])
		if err != nil {
			ui.PrintError(err.Error())
			return "", time.Time{}, false
		}
		return strings.TrimSpace(args[0]), birthday, true
	}

	prefilledName := ""
	if len(args) == 1 {
		prefilledName = strings.TrimSpace(args[0])
	}

	name, birthday, err := runAddForm(prefilledName)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			ui.PrintInfo("add canceled")
			return "", time.Time{}, false
		}

		ui.PrintError(fmt.Sprintf("could not open add form: %v", err))
		return "", time.Time{}, false
	}

	return name, birthday, true
}

func runAddForm(prefilledName string) (string, time.Time, error) {
	name := prefilledName
	birthday := ""

	fmt.Println(ui.RenderBanner("Add a birthday", "Fill in the form or press Ctrl+C to cancel."))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Placeholder("Anton").
				Value(&name).
				Validate(func(value string) error {
					if strings.TrimSpace(value) == "" {
						return fmt.Errorf("name is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Birthday").
				Description("Use format YYYY-MM-DD").
				Placeholder("1999-07-21").
				Value(&birthday).
				Validate(func(value string) error {
					_, err := parseBirthday(value)
					return err
				}),
		),
	).
		WithWidth(48).
		WithTheme(huh.ThemeFunc(huh.ThemeCharm))

	if err := form.Run(); err != nil {
		return "", time.Time{}, err
	}

	parsedBirthday, err := parseBirthday(birthday)
	if err != nil {
		return "", time.Time{}, err
	}

	return strings.TrimSpace(name), parsedBirthday, nil
}

func parseBirthday(value string) (time.Time, error) {
	birthday, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, fmt.Errorf("birthday must be in format YYYY-MM-DD")
	}

	return birthday, nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
