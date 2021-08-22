package externals

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"io/ioutil"
	"log"
	"net/http"
	"weekend.side/dialogFlowSlackSdk/functions/config"
	"weekend.side/dialogFlowSlackSdk/functions/utils"
)

type SlackRequest struct {
	*http.Request
	*slackevents.EventsAPIEvent
	credentials string
	*slack.InteractionCallback
}

func NewSlackRequest(req *http.Request, credentialsPath string) SlackRequest {
	return SlackRequest{req, nil, credentialsPath, nil}
}

func (slackReq SlackRequest) VerifyAndParseIncomingSlackRequests(signingSecret string, verifySecret bool) (respBody []byte, statusCode int, err error) {
	body, err := ioutil.ReadAll(slackReq.Body)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}

	if !verifySecret { // in case of url verification, secret header is not passed
		return body, 200, nil
	}

	sv, err := slack.NewSecretsVerifier(slackReq.Header, signingSecret)
	if err != nil {
		log.Print("error in secret  verification ", err)
		statusCode = http.StatusBadRequest
		return
	}

	if _, err = sv.Write(body); err != nil {
		log.Print("error in sv ", err)
		statusCode = http.StatusInternalServerError
		return
	}
	if err = sv.Ensure(); err != nil {
		log.Print("error in secret  ensure ", err)
		statusCode = http.StatusUnauthorized
		return
	}
	return body, 200, nil
}

func (slackReq *SlackRequest) HandleSlackRequests(body []byte) ([]byte, int, error) {

	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Print("error in parse event ", err)
		statusCode := http.StatusInternalServerError
		return nil, statusCode, err
	}

	slackReq.EventsAPIEvent = &eventsAPIEvent

	switch eventsAPIEvent.Type {

	case slackevents.URLVerification:

		r, slackUrlErr := slackReq.HandleSlackURLVerificationkEvent(body)

		if slackUrlErr != nil {
			statusCode := http.StatusInternalServerError
			return nil, statusCode, slackUrlErr
		}

		return []byte(r.Challenge), 200, nil

	case slackevents.CallbackEvent:

		respChat, slackEventCallbackErr := slackReq.HandleSlackCallbackEvent()

		if slackEventCallbackErr != nil {
			log.Print(slackEventCallbackErr)
			statusCode := http.StatusInternalServerError
			return nil, statusCode, slackEventCallbackErr
		}

		slackErr := slackReq.PostMsgToSlack(&slackReq.EventsAPIEvent.InnerEvent, nil, respChat)

		if slackErr != nil {
			log.Print(slackErr)
			statusCode := http.StatusInternalServerError
			return nil, statusCode, slackErr
		}
	}

	return []byte("OK"), 200, nil
}

func (slackReq *SlackRequest) PostMsgToSlack(innerEvent *slackevents.EventsAPIInnerEvent, interactiveCallbackMessage *slack.InteractionCallback, responseMessages []*cx.ResponseMessage) error {

	var api = slack.New(config.BOT_TOKEN) // can be moved to SlackRequest

	responseStr := utils.ParseStringFromResponse(responseMessages)

	blocks, _ := utils.ParsePayloadFromResponse(responseMessages)

	if innerEvent != nil {
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			if ev.ThreadTimeStamp == "" {
				api.PostMessage(ev.Channel, slack.MsgOptionTS(ev.TimeStamp), slack.MsgOptionText(responseStr, true), slack.MsgOptionBlocks(blocks...))
			} else {
				api.PostMessage(ev.Channel, slack.MsgOptionTS(ev.ThreadTimeStamp), slack.MsgOptionText(responseStr, true), slack.MsgOptionBlocks(blocks...))
			}
		case *slackevents.MessageEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText(responseStr, true), slack.MsgOptionBlocks(blocks...))
		}
	}

	if interactiveCallbackMessage != nil {
		api.PostMessage(interactiveCallbackMessage.Channel.ID, slack.MsgOptionText(responseStr, true), slack.MsgOptionBlocks(blocks...))
	}

	return nil
}

func (slackReq *SlackRequest) HandleSlackCallbackEvent() ([]*cx.ResponseMessage, error) {
	innerEvent := slackReq.EventsAPIEvent.InnerEvent

	isBot := ""

	// make a dialogflow request
	dialogflowcxReq := DialogFlowCXRequest{}

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		isBot = ev.BotID
		dialogflowcxReq = DialogFlowCXRequest{
			userInput:       utils.ParseTextFromInput(ev.Text),
			sessionId:       ev.User,
			credentialsPath: slackReq.credentials,
		}
	case *slackevents.MessageEvent:
		isBot = ev.BotID
		dialogflowcxReq = DialogFlowCXRequest{
			userInput:       utils.ParseTextFromInput(ev.Text),
			sessionId:       ev.User,
			credentialsPath: slackReq.credentials,
		}
	}

	if isBot != "" {
		return nil, errors.New("Can't reply to bot message")
	}

	return dialogflowcxReq.GetDialogFlowCXResponse()

}

func (slackReq *SlackRequest) HandleSlackInteractionEvent() ([]*cx.ResponseMessage, error) {
	actionCallbacks := slackReq.InteractionCallback.ActionCallback

	// make a dialogflow request
	dialogflowcxReq := DialogFlowCXRequest{}

	for _, blockAction := range actionCallbacks.BlockActions {
		fmt.Print("block action :", *blockAction, "\n\n\n\n")
		if blockAction != nil {
			dialogflowcxReq = DialogFlowCXRequest{
				userInput:       blockAction.Value,
				sessionId:       slackReq.InteractionCallback.User.ID,
				credentialsPath: slackReq.credentials,
			}
			return dialogflowcxReq.GetDialogFlowCXResponse()
		}
	}

	return nil, nil
}

func (slackReq *SlackRequest) HandleSlackURLVerificationkEvent(body []byte) (*slackevents.ChallengeResponse, error) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal(body, &r)
	return r, err
}
