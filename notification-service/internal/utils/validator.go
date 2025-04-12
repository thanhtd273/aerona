package utils

import "regexp"

func ValidatePhoneNumber(phoneNumber string) bool {
	if m, _ := regexp.MatchString("^0\\d{9}$", phoneNumber); !m {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	if m, _ := regexp.MatchString("^[_A-Za-z0-9-]+(\\.[_A-Za-z0-9-]+)*@[A-Za-z0-9]+(\\.[A-Za-z0-9]+)*(\\.[A-Za-z]{2,})$", email); !m {
		return false
	}
	return true
}
