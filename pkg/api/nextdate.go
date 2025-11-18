package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// function for comparing the start date and time from which the nearest date is searched
func afterNow(date, now time.Time) bool {
	return date.Before(now) || date.Equal(now)
}

// Task repetition function
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(fdate, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("an empty rule has been passed")
	}

	rule := parts[0]

	switch rule {
	case "d":
		return nextDays(start, now, parts)
	case "y":
		return nextYear(start, now)
	case "w":
		return nextWeekDay(start, now, parts)
	case "m":
		return nextMonth(start, now, parts)
	default:
		return "", fmt.Errorf("unknown rule %q, allowed: d, y, w", rule)
	}
}

// The nextDays function defines the next day
func nextDays(start, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect format for repeating the date")
	}

	d, err := strconv.Atoi(parts[1])
	if err != nil || d < 0 || d > 400 {
		return "", fmt.Errorf("incorrect format for repeating the date: %w", err)
	}

	if d == 0 {
		d = 1
	}

	date := start
	for {
		date = date.AddDate(0, 0, d)
		if date.After(now) {
			return date.Format(fdate), nil
		}
	}
}

// The nextYear function determines the next year
func nextYear(start, now time.Time) (string, error) {
	//реализация правила
	date := start
	for {
		date = date.AddDate(1, 0, 0)
		if !afterNow(date, now) {
			return date.Format(fdate), nil
		}
	}
}

// The nextWeekDay function determines the next day of the week
func nextWeekDay(start, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("incorrect format for repeating the date")
	}

	weekday := strings.Split(parts[1], ",")

	weekdayInt := []int{}
	for _, w := range weekday {
		d, err := strconv.Atoi(w)
		if err != nil || d > 7 || d < 0 {
			return "", fmt.Errorf("incorrect format for repeating the date: %w", err)
		}
		weekdayInt = append(weekdayInt, d)
	}

	date := start
	for {
		date = date.AddDate(0, 0, 1)
		dateWeekday := int(date.Weekday())
		if dateWeekday == 0 {
			dateWeekday = 7
		}
		if !afterNow(date, now) {
			for _, wd := range weekdayInt {
				if dateWeekday == wd {
					return date.Format(fdate), nil
				}
			}
		}
	}
}

// The function nextMonth will determine the next month and day of the month
func nextMonth(start, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("incorrect format for repeating the date")
	}

	days := []int{}
	months := []int{}

	for _, d := range strings.Split(parts[1], ",") {
		dInt, err := strconv.Atoi(d)
		if err != nil || dInt < -2 || dInt > 31 || dInt == 0 {
			return "", fmt.Errorf("incorrect format for repeating the date: %w", err)
		}
		days = append(days, dInt)
	}

	if len(parts) == 3 {
		for _, m := range strings.Split(parts[2], ",") {
			mInt, err := strconv.Atoi(m)
			if err != nil || mInt < 1 || mInt > 12 {
				return "", fmt.Errorf("incorrect format for repeating the date: %w", err)
			}
			months = append(months, mInt)
		}
	}

	date := start
	for {
		date = date.AddDate(0, 0, 1)
		if !date.After(now.Truncate(24 * time.Hour)) {
			continue
		}
		if len(months) > 0 {
			ok := false
			for _, m := range months {
				if int(date.Month()) == m {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		for _, d := range days {
			if d < 0 {
				y, m, _ := date.Date()
				lastday := time.Date(y, m+1, 0, 0, 0, 0, 0, date.Location())
				searchDay := lastday.Day() + d + 1
				if date.Day() == searchDay {
					return date.Format(fdate), nil
				}
			}
			if date.Day() == d {
				return date.Format(fdate), nil
			}
		}
	}
}

// A handler nextDayHandler that accepts a request in the format: "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>" and returns the date of the next task execution in the format: "20060102"
func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	now, err := time.Parse(fdate, r.FormValue("now"))
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	dateRes, err := NextDate(now, date, repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusFailedDependency)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(dateRes))
}
