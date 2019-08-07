package main

import (
	"fmt"
	"time"
)

func weekdayIndex(t time.Time) int {
	switch t.Weekday() {
	case time.Monday:
		return 0
	case time.Tuesday:
		return 1
	case time.Wednesday:
		return 2
	case time.Thursday:
		return 3
	case time.Friday:
		return 4
	default:
		return 0
	}
}

func weekdayName(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "Montag"
	case time.Tuesday:
		return "Dienstag"
	case time.Wednesday:
		return "Mittwoch"
	case time.Thursday:
		return "Donnerstag"
	case time.Friday:
		return "Freitag"
	case time.Saturday:
		return "Samstag"
	case time.Sunday:
		return "Sonntag"
	default:
		return "??"
	}
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%s, den %s", weekdayName(t), t.Format("02.01.2006"))
}
