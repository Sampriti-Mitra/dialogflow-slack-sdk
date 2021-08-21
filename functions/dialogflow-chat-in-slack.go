package functions

import (
	"fmt"
	"log"
	"net/http"
	"weekend.side/dialogFlowSlackSdk/functions/config"
	"weekend.side/dialogFlowSlackSdk/functions/externals"
)

func SimplestBotFunction(w http.ResponseWriter, r *http.Request) {

	signingSecret := config.SLACK_SIGNING_SECRET

	verifySecret := config.VERIFY_SECRET

	credentialsPath := config.CREDENTIALS_PATH

	slackReq := externals.NewSlackRequest(r, credentialsPath)

	body, statusCode, err := slackReq.VerifyAndParseIncomingSlackRequests(signingSecret, verifySecret)

	if err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, err)
		log.Print(err)
		return
	}

	log.Print(r.Header, string(body))

	w.Header().Set("X-Slack-No-Retry", "1")

	resp, statusCode, err := slackReq.HandleSlackRequests(body)

	w.WriteHeader(statusCode)

	if err != nil {
		log.Print("error in handling slack request ", err)
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, string(resp))

}
