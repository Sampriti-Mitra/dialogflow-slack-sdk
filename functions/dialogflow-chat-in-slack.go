package functions

import (
	"fmt"
	"log"
	"net/http"
	"weekend.side/dialogFlowSlackSdk/functions/config"
	"weekend.side/dialogFlowSlackSdk/functions/externals"
)

func SimplestBotFunction(w http.ResponseWriter, r *http.Request) {

	var isIncomingRequestVerified bool

	signingSecret := config.SLACK_SIGNING_SECRET

	credentialsPath := config.CREDENTIALS_PATH

	slackReq, err := externals.NewSlackRequest(r, credentialsPath)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	statusCode, err := slackReq.VerifyIncomingSlackRequests(r.Header, slackReq.Body, signingSecret)

	if err == nil {
		isIncomingRequestVerified = true
	}

	w.Header().Set("X-Slack-No-Retry", "1")

	resp, statusCode, err := slackReq.HandleSlackRequests(slackReq.Body, isIncomingRequestVerified)

	w.WriteHeader(statusCode)

	if err != nil {
		log.Print("error in handling slack request ", err)
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(resp))

}
