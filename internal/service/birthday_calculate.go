package service

import (
	"fmt"
	"time"
)

func calculateDaysUntilBD(birthday time.Time) (string, int) {
	now := time.Now()
	nextBD := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)

	if nextBD.Before(now) {
		nextBD = nextBD.AddDate(1, 0, 0)
	}

	days := int(nextBD.Sub(now).Hours() / 24)

	switch {
	case days == 0:
		return "Birthday today!", days
	case days == 1:
		return "Tomorrow birthday!", days
	case days < 7:
		return fmt.Sprintf("%d days", days), days
	case days < 30:
		weeks := days / 7
		if weeks == 1 {
			return "1 week ", days
		}
		return fmt.Sprintf("%d weeks", weeks), days
	default:
		months := days / 30
		remainingWeeks := (days - months*30) / 7

		monthStr := "1 month"
		if months > 1 {
			monthStr = fmt.Sprintf("%d months", months)
		}

		if remainingWeeks == 0 {
			return fmt.Sprintf("%s", monthStr), days
		}
		weekStr := "1 week"
		if remainingWeeks > 1 {
			weekStr = fmt.Sprintf("%d weeks", remainingWeeks)
		}
		return fmt.Sprintf("%s %s", monthStr, weekStr), days
	}
}

func calculateAge(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	// If your birthday hasn't arrived yet this year, subtract 1.
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}
