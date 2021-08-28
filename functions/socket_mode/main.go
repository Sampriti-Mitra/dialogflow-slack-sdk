package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"net/http"
	"os"
	"weekend.side/dialogFlowSlackSdk/functions/config"
	"weekend.side/dialogFlowSlackSdk/functions/externals"
)

func main() {
	botToken := config.BOT_TOKEN
	appToken := config.APP_TOKEN

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				continue
			case socketmode.EventTypeConnectionError:
				continue
			case socketmode.EventTypeConnected:
				continue
			case socketmode.EventTypeHello:
				continue
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					slackReq, slackReqErr := externals.NewSlackRequest(nil, config.CREDENTIALS_PATH)
					if slackReqErr != nil {
						//log.Print(slackReqErr)
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}
					slackReq.EventsAPIEvent = &eventsAPIEvent

					respChat, slackEventCallbackErr := slackReq.SendSlackCallbackEventToDialogflowCxAndGetResponse()

					if slackEventCallbackErr != nil {
						log.Print(slackEventCallbackErr)
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}

					slackErr := slackReq.PostMsgToSlack(&slackReq.EventsAPIEvent.InnerEvent, nil, respChat)

					if slackErr != nil {
						log.Print(slackErr)
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}
				default:
					client.Debugf("unsupported Events API event received")
				}

			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slack.InteractionCallback)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)

					continue
				}
				client.Ack(*evt.Request)

				var payload interface{}

				switch callback.Type {
				case slack.InteractionTypeBlockActions:

					slackReq, slackReqErr := externals.NewSlackRequest(nil, config.CREDENTIALS_PATH)
					if slackReqErr != nil {
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}
					slackReq.InteractionCallback = &callback

					respChat, slackInteractionEventCallbackErr := slackReq.SendSlackInteractionEventToDialogFlowCxAndGetResponse()

					if slackInteractionEventCallbackErr != nil {
						log.Print(slackInteractionEventCallbackErr)
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}

					slackErr := slackReq.UpdateInteractiveSlackMessage(&callback, respChat)

					if slackErr != nil {
						log.Print(slackErr)
						statusCode := http.StatusInternalServerError
						fmt.Print(statusCode)
						continue
					}
				case slack.InteractionTypeShortcut:
				case slack.InteractionTypeViewSubmission:
					// See https://api.slack.com/apis/connections/socket-implement#modal
				case slack.InteractionTypeDialogSubmission:
				default:

				}

				client.Ack(*evt.Request, payload)

			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				client.Debugf("Slash command received: %+v", cmd)
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()

}
