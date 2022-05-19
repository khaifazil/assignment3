package main

import (
	"errors"
	"time"
)

func IsAlphabetic(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < 0 || r > 9) {
			return false
		}
	}
	return true
}

func checkDate(date string) error {
	parsedDate, err := time.Parse(timeFormat, date)
	if err != nil {
		return err
	} else if parsedDate.Before(time.Now()) {
		return errors.New("date given has passed")
	}
	return nil
}

func checkContactLen(n int) bool {
	if n > 99999999 || n < 10000000 {
		return false
	}
	return true
}
