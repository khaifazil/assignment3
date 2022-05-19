package limoBookingApp

import (
	"errors"
	"time"
)

//IsAlphabetic checks if the given string is only Alphabetic, Upper & Lowercase. Returns false if it is not.
func IsAlphabetic(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

//IsAlphanumeric checks if the given string is only Alphanumeric, Upper & Lowercase. Returns false if it is not.
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < 0 || r > 9) {
			return false
		}
	}
	return true
}

//CheckDate validates the user date input to make sure that given date has not passed. Returns error if date has passed.
func CheckDate(date string) error {
	parsedDate, err := time.Parse(timeFormat, date)
	if err != nil {
		return err
	} else if parsedDate.Before(time.Now()) {
		return errors.New("date given has passed")
	}
	return nil
}

//CheckContactLen validates the user contact input. If length is more or lest than 8 digits, false is returned.
func CheckContactLen(n int) bool {
	if n > 99999999 || n < 10000000 {
		return false
	}
	return true
}
