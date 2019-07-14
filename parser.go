package main

import (
	"regexp"
	"strings"
)

type dailyMenu struct {
	M1 menuItem
	M2 menuItem
}

type menuItem struct {
	title    string
	subtitle string
	price    string
}

var weekdayNames = []string{"Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag"}

const weekdays = 5

func parseMenu(text string) ([]dailyMenu, error) {
	result := make([]dailyMenu, 0)
	lines := toLines(text)
	idx := 0

	for idx = 0; idx < len(lines) && len(result) < weekdays; idx++ {
		line := lines[idx]
		if isWeekday(line) {
			var m1, m2 menuItem

			idx = skip(lines, idx+1)
			m1, idx = readMenuItem(lines, idx)
			idx = skip(lines, idx)
			m2, idx = readMenuItem(lines, idx)

			result = append(result, dailyMenu{
				M1: m1,
				M2: m2,
			})
		}
	}

	itemIdx := 0

	for ; idx < len(lines) && itemIdx < weekdays; idx++ {
		idx = skip(lines, idx)

		line := lines[idx]
		if isPrice(line) {
			result[itemIdx].M1.price = line
			idx++

			idx = skip(lines, idx)
			if idx < len(lines) && isPrice(lines[idx]) {
				result[itemIdx].M2.price = lines[idx]
				idx++
			}

			itemIdx++
		}
	}

	return result, nil
}

func toLines(text string) []string {
	lines := strings.Split(text, "\n")

	for idx, line := range lines {
		lines[idx] = strings.TrimSpace(line)
	}

	return lines
}

func isWeekday(line string) bool {
	for _, weekday := range weekdayNames {
		if line == weekday {
			return true
		}
	}

	return false
}

var price = regexp.MustCompile(`^[0-9,.]+ €$`)

func isPrice(line string) bool {
	return price.MatchString(line)
}

func skip(lines []string, idx int) int {
	for idx < len(lines) && lines[idx] == "" {
		idx++
	}

	return idx
}

func readMenuItem(lines []string, idx int) (menuItem, int) {
	item := menuItem{}

	if idx < len(lines) && lines[idx] != "" {
		item.title = lines[idx]
		idx++
	}

	if idx < len(lines) && lines[idx] != "" {
		item.subtitle = lines[idx]
		idx++
	}

	return item, idx
}
