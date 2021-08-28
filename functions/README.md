# Deploying with Cloud Functions
We can also deploy this using cloud functions deployment.
Follow the instructions on the ReadMe at the root of this project.
For deploying with Cloud Functions, Go to Cloud Functions and enable the API, also enable Cloud Build and Deploy.

## gcloud CLI set up
The deployment process for GCP Cloud functions via this README utilizes gcloud CLI commands. Follow the steps in the main README to set up gcloud CLI locally for this deployment.


## Deploying Cloud Functions
1. Ensure in the functions directory, you have a go.mod. If you don't, 
   a. delete the go.mod in the root directory and cd to functions directory<br>
   b. go mod init
   c. go mod tidy
1. On the terminal, cd to the functions directory of the cloned project and `gcloud functions deploy SimplestBotFunction --runtime go113 --trigger-http --allow-unauthenticated`
   This will deploy your project.
   
2. You need to enable Events Subscription and turn off socket mode(if it is turned on) for deploying with cloud functions.

3. Don't forget to copy the BOT_TOKEN in basic information for the app and replace the value in config/token.go

On slack, go to the channel(s) you want the slack bot to have access to and invite the bot to the channel. Refer the below image for adding the bot to a channel. Alternatively, you can type /invite on the channel

## Post-deployment

### Shutting Down an Integration

In order to shut down an integration set up via the steps in this README, you need to ` gcloud functions delete SimplestBotFunction `. This will delete your cloud functions.