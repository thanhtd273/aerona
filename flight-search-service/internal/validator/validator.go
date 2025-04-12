package validator

import (
	"regexp"
	"time"
)

func ValidateDateString(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func IsNumeric(str string) bool {
	return regexp.MustCompile(`\d`).MatchString(str)
}
