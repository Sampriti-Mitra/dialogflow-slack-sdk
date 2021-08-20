package externals

import (
	dialogflowcx "cloud.google.com/go/dialogflow/cx/apiv3"
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3"
	"log"
	"weekend.side/dialogFlowSlackSdk/functions/config"
)

type DialogFlowCXRequest struct {
	credentialsPath string
	userInput       string
	sessionId       string
}

func (dialogflowcxReq DialogFlowCXRequest) GetDialogFlowCXResponse() (interface{}, error) {

	ProjectId := config.PROJECT_ID // project id

	agent := config.AGENT_NAME

	ctx := context.Background()

	sa := option.WithCredentialsFile(dialogflowcxReq.credentialsPath)

	detectIntentReq := cx.DetectIntentRequest{
		Session: agent + "/sessions/" + ProjectId + dialogflowcxReq.sessionId,
		QueryInput: &cx.QueryInput{
			Input: &cx.QueryInput_Text{
				&cx.TextInput{
					Text: dialogflowcxReq.userInput,
				},
			},
			LanguageCode: "en",
		},
	}

	opts := []option.ClientOption{
		internaloption.WithDefaultEndpoint("us-central1-dialogflow.googleapis.com:443"),
		internaloption.WithDefaultAudience("https://us-central1-dialogflow.googleapis.com/"),
	}

	dialogFlowClient, err := dialogflowcx.NewSessionsClient(ctx, append(opts, sa)...)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	resp, err := dialogFlowClient.DetectIntent(ctx, &detectIntentReq)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	queryResult := resp.GetQueryResult()

	responseMessages := queryResult.GetResponseMessages()

	return responseMessages, nil
}
