package service

import (
	"fmt"
	"time"
)

const (
	OneMonthThreshold = 32
)

func calculateAge(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	// If your birthday hasn't arrived yet this year, subtract 1.
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

func calculateDaysUntilBD(birthday time.Time) (string, int) {
	now := time.Now().Truncate(24 * time.Hour)
	nextBD := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)

	if nextBD.Before(now) {
		nextBD = nextBD.AddDate(1, 0, 0)
	}

	days := int(nextBD.Sub(now).Hours() / 24)

	switch {
	case days == 0:
		return "Birthday today!", days
	case days == 1:
		return "Tomorrow!", days
	case days < 7:
		return fmt.Sprintf("%d days", days), days
	case days < OneMonthThreshold:
		weeks := days / 7
		remDays := days % 7
		if remDays == 0 {
			return fmt.Sprintf("%d weeks", weeks), days
		}
		return fmt.Sprintf("%d weeks, %d days", weeks, remDays), days
	default:
		months := calculateMonthsUntil(now, nextBD)
		afterMonths := now.AddDate(0, months, 0)
		remainingDays := int(nextBD.Sub(afterMonths).Hours() / 24)

		if remainingDays == 0 {
			return fmt.Sprintf("%d months", months), days
		}

		weeks := remainingDays / 7
		remDays := remainingDays % 7

		if weeks == 0 {
			return fmt.Sprintf("%d months, %d days", months, remDays), days
		}
		if remDays == 0 {
			return fmt.Sprintf("%d months, %d weeks", months, weeks), days
		}
		return fmt.Sprintf("%d months, %d weeks, %d days", months, weeks, remDays), days
	}
}

func calculateMonthsUntil(start, end time.Time) int {
	months := 0
	for start.AddDate(0, months+1, 0).Before(end) {
		months++
	}
	return months
}
