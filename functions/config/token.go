package config

const (
	BOT_TOKEN = "" // required for events api callback events
	APP_TOKEN = "" // required for socket mode

	// if below line is uncommented, dialogflowcx,json must be present at the location relative to root
	//CREDENTIALS_PATH     = "functions/config/dialogflowcx.json"

	SLACK_SIGNING_SECRET = ""
	// set this to false when url verification is not yet done
	VERIFY_SECRET = true
	PROJECT_ID    = ""
	AGENT_NAME    = ""
)
