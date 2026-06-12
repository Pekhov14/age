package ui

import (
	"age/internal/service"
	"fmt"
	"slices"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

var (
	accentColor  = lipgloss.Color("#A78BFA")
	mutedColor   = lipgloss.Color("#94A3B8")
	borderColor  = lipgloss.Color("#7C3AED")
	successColor = lipgloss.Color("#4ADE80")
	warningColor = lipgloss.Color("#FBBF24")
	errorColor   = lipgloss.Color("#F87171")
	infoColor    = lipgloss.Color("#60A5FA")

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	headerCellStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	cardValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	commandStyle = lipgloss.NewStyle().
			Foreground(infoColor)

	nameCellStyle = lipgloss.NewStyle().Width(18)
	dateCellStyle = lipgloss.NewStyle().Width(12).Foreground(mutedColor)
	ageCellStyle  = lipgloss.NewStyle().Width(10)
	daysCellStyle = lipgloss.NewStyle().Width(8)
	nextCellStyle = lipgloss.NewStyle().Width(24)

	successTextStyle = lipgloss.NewStyle().Bold(true).Foreground(successColor)
	errorTextStyle   = lipgloss.NewStyle().Bold(true).Foreground(errorColor)
	infoTextStyle    = lipgloss.NewStyle().Bold(true).Foreground(infoColor)
)

func RenderBanner(title, subtitle string) string {
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		subtitleStyle.Render(subtitle),
	)

	return panelStyle.Render(body)
}

func RenderBirthdayPreview(person service.PersonInfo) string {
	cards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderStatCard("Age", fmt.Sprintf("%d years", person.Age)),
		renderStatCard("Next birthday", compactCountdownLabel(person.DaysUntilBD)),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		cards,
	)
}

func RenderDashboard(persons []service.PersonInfo) string {
	if len(persons) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			RenderBanner("No birthdays yet", "Start by adding someone with `age add`."),
			// renderQuickActions(),
		)
	}

	sorted := sortedByUpcoming(persons)
	todayCount := 0
	for _, person := range persons {
		if person.DaysToBirth == 0 {
			todayCount++
		}
	}

	nextUp := fmt.Sprintf("%s in %s", sorted[0].Name, compactCountdownLabel(sorted[0].DaysUntilBD))
	cards := lipgloss.JoinHorizontal(
		lipgloss.Top,
		renderStatCard("Saved", fmt.Sprintf("%d people", len(persons))),
		renderStatCard("Today", fmt.Sprintf("%d birthdays", todayCount)),
		renderStatCard("Next up", nextUp),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		RenderBanner("Birthday dashboard", "A quick glance at the people you track."),
		cards,
		// renderQuickActions(),
	)
}

func RenderList(persons []service.PersonInfo) string {
	if len(persons) == 0 {
		return RenderBanner("No birthdays yet", "Add someone with `age add` and your reminders will show up here.")
	}

	persons = sortedByUpcoming(persons)
	headers := lipgloss.JoinHorizontal(
		lipgloss.Top,
		headerCellStyle.Width(18).Render("Name"),
		headerCellStyle.Width(12).Render("Birthday"),
		headerCellStyle.Width(10).Render("Age"),
		headerCellStyle.Width(8).Render("Days"),
		headerCellStyle.Width(24).Render("Next birthday"),
	)

	separator := lipgloss.NewStyle().Foreground(borderColor).Render(strings.Repeat("─", 72))

	rows := make([]string, 0, len(persons))
	for _, person := range persons {
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			nameCellStyle.Render(person.Name),
			dateCellStyle.Render(person.Birthday),
			ageCellStyle.Render(fmt.Sprintf("%d years", person.Age)),
			daysCellStyle.Render(fmt.Sprintf("%d", person.DaysToBirth)),
			nextCellStyle.Render(renderCountdownBadge(person.DaysToBirth, compactCountdownLabel(person.DaysUntilBD))),
		)

		rows = append(rows, row)
	}

	summary := subtitleStyle.Render(fmt.Sprintf("Tracking %d birthdays", len(persons)))
	nextUpBlock := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Next birthday"),
		subtitleStyle.Render(fmt.Sprintf("%s %s", persons[0].Name, compactCountdownLabel(persons[0].DaysUntilBD))),
	)
	headerBlock := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Upcoming birthdays"),
			summary,
		),
		lipgloss.NewStyle().MarginLeft(8).Render(nextUpBlock),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBlock,
		"",
		headers,
		separator,
		strings.Join(rows, "\n"),
	)

	return panelStyle.Render(content)
}

func PrintSuccess(message string) {
	fmt.Println(successTextStyle.Render("✓ " + message))
}

func PrintError(message string) {
	fmt.Println(errorTextStyle.Render("✗ " + message))
}

func PrintInfo(message string) {
	fmt.Println(infoTextStyle.Render("• " + message))
}

func RenderAddSuccess(name string, birthday time.Time) string {
	return RenderBanner(
		"Birthday saved",
		fmt.Sprintf("%s · %s", name, birthday.Format("2006-01-02")),
	)
}

// func renderQuickActions() string {
// 	lines := []string{
// 		commandStyle.Render("age 2000-09-04") + subtitleStyle.Render(" quick age preview without saving"),
// 		commandStyle.Render("age list") + subtitleStyle.Render("        show all saved birthdays"),
// 		commandStyle.Render("age add") + subtitleStyle.Render("         open the interactive add form"),
// 		commandStyle.Render("age update") + subtitleStyle.Render("      update a saved name and birthday"),
// 		commandStyle.Render("age delete") + subtitleStyle.Render("      delete a saved person with confirmation"),
// 	}

// 	return panelStyle.Render(lipgloss.JoinVertical(
// 		lipgloss.Left,
// 		titleStyle.Render("Quick actions"),
// 		strings.Join(lines, "\n"),
// 	))
// }

func renderStatCard(label, value string) string {
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		subtitleStyle.Render(label),
		cardValueStyle.Render(value),
	)

	return cardStyle.Width(22).Render(body)
}

func sortedByUpcoming(persons []service.PersonInfo) []service.PersonInfo {
	sorted := slices.Clone(persons)
	slices.SortFunc(sorted, func(a, b service.PersonInfo) int {
		if a.DaysToBirth != b.DaysToBirth {
			return a.DaysToBirth - b.DaysToBirth
		}
		return strings.Compare(a.Name, b.Name)
	})
	return sorted
}

func renderCountdownBadge(days int, label string) string {
	style := lipgloss.NewStyle().Padding(0, 1).Bold(true)

	switch {
	case days == 0:
		style = style.Foreground(successColor)
	case days <= 7:
		style = style.Foreground(warningColor)
	default:
		style = style.Foreground(infoColor)
	}

	return style.Render(label)
}

func compactCountdownLabel(label string) string {
	replacer := strings.NewReplacer(
		"Birthday today!", "today",
		"Tomorrow!", "tomorrow",
		" months", "mo",
		" month", "mo",
		" weeks", "w",
		" week", "w",
		" days", "d",
		" day", "d",
		", ", " ",
	)

	return replacer.Replace(label)
}
