package utils

import "strings"

func IsCapitalized(str string) bool {
	if len(str) == 0 {
		return false
	}
	return str[0] >= 'A' && str[0] <= 'Z'
}

func ToSentenceCase(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func ToUpperCase(str string) string {
	return strings.ToUpper(str)
}

func ToLowerCase(str string) string {
	return strings.ToLower(str)
}

func Plural(singular, plural string, count int) string {
	if count == 1 {
		return singular
	}
	return plural
}
