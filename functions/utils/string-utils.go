package utils

import (
	"google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"strings"
)

func ParseTextFromInput(text string) string {
	ind := strings.Index(text, ">")
	return text[ind+1:]
}

func ParseStringFromResponse(responseMessages []*cx.ResponseMessage) string {
	var str string
	for _, responseMessage := range responseMessages {
		str += GetStringFromSlice(responseMessage.GetText().GetText())
	}
	return str
}

func GetStringFromSlice(strSlice []string) string {
	var respStr string
	for _, str := range strSlice {
		respStr += str
	}
	return respStr
}
