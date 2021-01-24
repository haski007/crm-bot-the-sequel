package validate

import (
	"regexp"
	"strings"
)

// PhoneNumber returns valid number or empty string
func PhoneNumber(input string) string {

	// ---> Get all digits from string
	reg := regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*`)
	phone := strings.Join(reg.FindAllString(input, -1), "")

	// ---> Validate number
	reg = regexp.MustCompile(`^(\+?3?8?)0\d{9}$`)

	if !reg.MatchString(phone) {
		return ""
	}

	code := reg.FindStringSubmatch(phone)[1]

	switch code {
	case "":
		phone = "+38" + phone
	case "38":
		phone = "+" + phone
	case "8":
		phone = "+3" + phone
	}

	return phone
}
