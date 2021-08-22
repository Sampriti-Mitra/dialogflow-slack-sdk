package utils

import (
	"encoding/json"
	"github.com/slack-go/slack"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"log"
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

func ParsePayloadFromResponse(responseMessages []*cx.ResponseMessage) ([]slack.Block, bool) {
	var blocksResp []slack.Block
	var anyCustomPayloadExists bool

	for _, responseMessage := range responseMessages {
		bytes, err := responseMessage.GetPayload().MarshalJSON()
		if err != nil {
			continue
		}
		if string(bytes) != "{}" {
			blocks := ParsePayloadMessage(bytes)
			anyCustomPayloadExists = true
			return blocks, anyCustomPayloadExists
		}
	}
	return blocksResp, anyCustomPayloadExists
}

func GetStringFromSlice(strSlice []string) string {
	var respStr string
	for _, str := range strSlice {
		respStr += str
	}
	return respStr
}

func ParsePayloadMessage(bytes []byte) []slack.Block {
	// we convert the view into a message struct
	views := slack.Msg{}

	err := json.Unmarshal(bytes, &views)

	if err != nil {
		log.Print("error in unmarshalling payload to message ", err)
		return nil
	}

	return views.Blocks.BlockSet
}
