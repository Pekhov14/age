package cmd

import (
	"age/internal/ui"
	"errors"
	"fmt"
	"strings"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

var deleteWithoutConfirm bool

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a saved birthday",
	Long:  "Remove a person from the birthday list by name, or choose one interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, ok := collectDeleteName(args)
		if !ok {
			return
		}

		if !deleteWithoutConfirm {
			confirmed, err := runDeleteConfirm(name)
			if err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					ui.PrintInfo("delete canceled")
					return
				}

				ui.PrintError(fmt.Sprintf("could not open confirmation: %v", err))
				return
			}

			if !confirmed {
				ui.PrintInfo("delete canceled")
				return
			}
		}

		if err := personService.Delete(name); err != nil {
			ui.PrintError(fmt.Sprintf("could not delete %s: %v", name, err))
			return
		}

		ui.PrintSuccess(fmt.Sprintf("deleted %s", name))
	},
}

func collectDeleteName(args []string) (string, bool) {
	if len(args) == 1 {
		name := strings.TrimSpace(args[0])
		if name == "" {
			ui.PrintError("name is required")
			return "", false
		}
		return name, true
	}

	persons, err := personService.List()
	if err != nil {
		ui.PrintError(fmt.Sprintf("could not load birthdays: %v", err))
		return "", false
	}

	if len(persons) == 0 {
		ui.PrintInfo("there is nobody to delete yet")
		return "", false
	}

	selectedName := ""
	options := make([]huh.Option[string], 0, len(persons))
	for _, person := range persons {
		options = append(options, huh.NewOption(person.Name, person.Name))
	}

	fmt.Println(ui.RenderBanner("Choose a person", "Pick who you want to delete."))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Saved people").
				Options(options...).
				Value(&selectedName),
		),
	).
		WithWidth(56).
		WithTheme(huh.ThemeFunc(huh.ThemeCharm))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			ui.PrintInfo("delete canceled")
			return "", false
		}

		ui.PrintError(fmt.Sprintf("could not open picker: %v", err))
		return "", false
	}

	return selectedName, true
}

func runDeleteConfirm(name string) (bool, error) {
	confirmed := false

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Delete %s?", name)).
				Description("This will remove the person from your saved birthday list.").
				Affirmative("Delete").
				Negative("Keep").
				Value(&confirmed),
		),
	).
		WithWidth(56).
		WithTheme(huh.ThemeFunc(huh.ThemeCharm))

	if err := form.Run(); err != nil {
		return false, err
	}

	return confirmed, nil
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteWithoutConfirm, "yes", "y", false, "Delete without confirmation")
	rootCmd.AddCommand(deleteCmd)
}
