package cmd

import (
	"age/internal/service"
	"age/internal/ui"
	"errors"
	"fmt"
	"strings"
	"time"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update [name]",
	Aliases: []string{"rename"},
	Short:   "Update a saved person",
	Long:    "Update a saved person's name and/or birthday.",
	Example: strings.TrimSpace(`
  age update
  age update Anton
  age update Anton   # opens a form with current values
`),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		oldName, newName, birthday, ok := collectUpdateValues(args)
		if !ok {
			return
		}

		if err := personService.Update(oldName, newName, birthday); err != nil {
			ui.PrintError(fmt.Sprintf("could not update %s: %v", oldName, err))
			return
		}

		ui.PrintSuccess(fmt.Sprintf("updated %s", newName))
	},
}

func collectUpdateValues(args []string) (string, string, time.Time, bool) {
	persons, err := personService.List()
	if err != nil {
		ui.PrintError(fmt.Sprintf("could not load birthdays: %v", err))
		return "", "", time.Time{}, false
	}

	if len(persons) == 0 {
		ui.PrintInfo("there is nobody to update yet")
		return "", "", time.Time{}, false
	}

	selectedPerson, err := pickPersonForUpdate(persons, args)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			ui.PrintInfo("update canceled")
			return "", "", time.Time{}, false
		}

		ui.PrintError(err.Error())
		return "", "", time.Time{}, false
	}

	newName, birthday, err := runUpdateForm(selectedPerson)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			ui.PrintInfo("update canceled")
			return "", "", time.Time{}, false
		}

		ui.PrintError(fmt.Sprintf("could not open update form: %v", err))
		return "", "", time.Time{}, false
	}

	return selectedPerson.Name, newName, birthday, true
}

func pickPersonForUpdate(persons []service.PersonInfo, args []string) (service.PersonInfo, error) {
	if len(args) == 1 {
		name := strings.TrimSpace(args[0])
		person, ok := findPersonByName(persons, name)
		if !ok {
			return service.PersonInfo{}, fmt.Errorf("person %q not found", name)
		}
		return person, nil
	}

	selectedName := ""
	options := make([]huh.Option[string], 0, len(persons))
	for _, person := range persons {
		options = append(options, huh.NewOption(person.Name, person.Name))
	}

	fmt.Println(ui.RenderBanner("Choose a person", "Pick who you want to update."))

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
		return service.PersonInfo{}, err
	}

	person, ok := findPersonByName(persons, selectedName)
	if !ok {
		return service.PersonInfo{}, fmt.Errorf("person %q not found", selectedName)
	}

	return person, nil
}

func runUpdateForm(person service.PersonInfo) (string, time.Time, error) {
	name := person.Name
	birthday := person.Birthday

	fmt.Println(ui.RenderBanner("Update a birthday", "Change the name, the birth date, or both."))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
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
				Value(&birthday).
				Validate(func(value string) error {
					_, err := parseBirthday(value)
					return err
				}),
		),
	).
		WithWidth(56).
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

func findPersonByName(persons []service.PersonInfo, name string) (service.PersonInfo, bool) {
	for _, person := range persons {
		if person.Name == name {
			return person, true
		}
	}

	return service.PersonInfo{}, false
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
