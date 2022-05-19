package limoBookingApp

import "regexp"

//StripHtmlRegex strips any HTML tags from the given string.
func StripHtmlRegex(s string) string {
	const regex = `<.*?>`
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}
