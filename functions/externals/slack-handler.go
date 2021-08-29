package externals

import (
	"encoding/json"
	"errors"
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
	Body []byte
	*slackevents.EventsAPIEvent
	credentials string
	*slack.InteractionCallback
}

func NewSlackRequest(req *http.Request, credentialsPath string) (*SlackRequest, error) {

	var payload *slack.InteractionCallback

	var body []byte

	if req != nil {
		err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
		if err != nil {
			payload = nil
		}

		if req.Body != nil {
			body, err = ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
		}
	}

	return &SlackRequest{body, nil, credentialsPath, payload}, nil
}

func (slackReq SlackRequest) VerifyIncomingSlackRequests(headers http.Header, body []byte, signingSecret string) (statusCode int, err error) {

	sv, err := slack.NewSecretsVerifier(headers, signingSecret)
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
	return 200, nil
}

func (slackReq *SlackRequest) HandleSlackRequests(body []byte, isIncomingRequestVerified bool) ([]byte, int, error) {

	switch {
	// if interaction callback event
	case slackReq.InteractionCallback != nil:
		return slackReq.HandleInteractionCallbackEvents()
	}
	// if event callback
	return slackReq.HandleEventsApiCallbackEvents(body, isIncomingRequestVerified)
	//return nil, 500, errors.New("Unsupported event")
}

func (slackReq *SlackRequest) HandleInteractionCallbackEvents() ([]byte, int, error) {
	respChat, slackInteractionEventCallbackErr := slackReq.SendSlackInteractionEventToDialogFlowCxAndGetResponse()

	if slackInteractionEventCallbackErr != nil {
		return nil, http.StatusInternalServerError, slackInteractionEventCallbackErr
	}

	slackErr := slackReq.UpdateInteractiveSlackMessage(slackReq.InteractionCallback, respChat)

	if slackErr != nil {
		log.Print(slackErr)
		statusCode := http.StatusInternalServerError
		return nil, statusCode, slackErr
	}

	return []byte("OK"), 200, nil
}

func (slackReq *SlackRequest) HandleEventsApiCallbackEvents(body []byte, isIncomingRequestVerified bool) ([]byte, int, error) {
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

		// if request is not verified, only url verification is allowed,
		//no action to callback events are allowed
		if !isIncomingRequestVerified {
			return nil, 400, errors.New("slack request needs signing secret header")
		}

		respChat, slackEventCallbackErr := slackReq.SendSlackCallbackEventToDialogflowCxAndGetResponse()

		if slackEventCallbackErr != nil {
			log.Print(slackEventCallbackErr)
			statusCode := http.StatusInternalServerError
			return nil, statusCode, slackEventCallbackErr
		}

		slackErr := slackReq.PostMsgToSlack(respChat)

		if slackErr != nil {
			log.Print(slackErr)
			statusCode := http.StatusInternalServerError
			return nil, statusCode, slackErr
		}
	default:
		return nil, 400, errors.New("Type not supported")
	}

	return []byte("OK"), 200, nil
}

func (slackReq *SlackRequest) PostMsgToSlack(responseMessages []*cx.ResponseMessage) error {

	innerEvent := &slackReq.EventsAPIEvent.InnerEvent

	var api = slack.New(config.BOT_TOKEN) // can be moved to SlackRequest

	responseStr := utils.ParseStringFromResponse(responseMessages)

	blocks, _ := utils.ParsePayloadFromResponse(responseMessages)

	// if it is an eventsApi callback event
	// then post as separate message if it is a DM message
	// otherwise post as a reply to the existing thread in channel

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

	return nil
}

func (slackReq *SlackRequest) UpdateInteractiveSlackMessage(interactiveCallbackMessage *slack.InteractionCallback, responseMessages []*cx.ResponseMessage) error {

	responseStr := utils.ParseStringFromResponse(responseMessages)

	blocks, _ := utils.ParsePayloadFromResponse(responseMessages)

	webhookMsg := slack.WebhookMessage{
		Channel: interactiveCallbackMessage.Channel.ID,
		Blocks:  &slack.Blocks{BlockSet: blocks},
		Text:    responseStr,
	}

	// if it is an interaction event, then post as separate message if DM
	// if channel, post as a reply to the thread
	if interactiveCallbackMessage != nil {
		if !interactiveCallbackMessage.Channel.IsIM {
			webhookMsg.ThreadTimestamp = interactiveCallbackMessage.Container.ThreadTs
		}
		slack.PostWebhook(interactiveCallbackMessage.ResponseURL, &webhookMsg)
	}

	return nil
}

func (slackReq *SlackRequest) SendSlackCallbackEventToDialogflowCxAndGetResponse() ([]*cx.ResponseMessage, error) {
	innerEvent := slackReq.EventsAPIEvent.InnerEvent

	var botId, text, user string

	// make a dialogflow request
	dialogflowcxReq := DialogFlowCXRequest{}

	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		botId = ev.BotID
		text = utils.ParseTextFromInput(ev.Text)
		user = ev.User

	case *slackevents.MessageEvent:
		botId = ev.BotID
		text = utils.ParseTextFromInput(ev.Text)
		user = ev.User
	}

	dialogflowcxReq = DialogFlowCXRequest{
		userInput:       text,
		sessionId:       user,
		credentialsPath: slackReq.credentials,
	}

	if botId != "" || text == "" {
		return nil, errors.New("Can't reply to bot message")
	}

	return dialogflowcxReq.GetDialogFlowCXResponse()

}

func (slackReq *SlackRequest) SendSlackInteractionEventToDialogFlowCxAndGetResponse() ([]*cx.ResponseMessage, error) {
	actionCallbacks := slackReq.InteractionCallback.ActionCallback

	// make a dialogflow request
	dialogflowcxReq := DialogFlowCXRequest{}

	for _, blockAction := range actionCallbacks.BlockActions {
		if blockAction != nil && blockAction.Value != "" {
			dialogflowcxReq = DialogFlowCXRequest{
				userInput:       blockAction.Value,
				sessionId:       slackReq.InteractionCallback.User.ID,
				credentialsPath: slackReq.credentials,
			}
			return dialogflowcxReq.GetDialogFlowCXResponse()
		}
	}

	return nil, errors.New("no proper input to dialgflowcx")
}

func (slackReq *SlackRequest) HandleSlackURLVerificationkEvent(body []byte) (*slackevents.ChallengeResponse, error) {
	var r *slackevents.ChallengeResponse
	err := json.Unmarshal(body, &r)
	return r, err
}
