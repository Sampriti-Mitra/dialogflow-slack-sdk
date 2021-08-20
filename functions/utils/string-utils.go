package utils

import "strings"

func ParseTextFromInput(text string)string{
	ind:= strings.Index(text,">")
	return text[ind+1:]
}
