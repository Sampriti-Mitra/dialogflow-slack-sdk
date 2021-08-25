# dialogflow-slack-sdk

## Introduction
The goal of this guide is to show you how to set up an integration deployment to link your Dialogflow agent to slack.
If you don't already have a Dialogflow agent, you may create one by following the instructions here or by adding a prebuilt agent.
Although this integration deployment may be set up on any other hosting platform, these instructions will use Google's App Engine/Cloud functions.

## GCP Setup
Step 1: Log in or Sign up for GCP
Log in or sign up to google cloud console using a credit or debit card for a free trial. Create a project and enable billing for that project. Go to Cloud Functions and enable the API, also enable Cloud Build and Deploy.


### gcloud CLI setup

The deployment process for GCP App Engine and Cloud Functions via this README utilizes gcloud CLI commands. Follow the steps below to set up gcloud CLI locally for this deployment.

1. On the gcloud CLI [documentation page](https://cloud.google.com/sdk/docs/quickstarts), select your OS and follow the instructions for the installation.
2. Run ``gcloud config get-value project`` to check the GCP Project configured.

### Service Account Setup (GCP)

For the integration to function properly, it is necessary to create a Service Account in your agent’s GCP Project. See [this page](https://cloud.google.com/dialogflow/docs/quick/setup#sa-create) of the documentation for more details.

Follow the steps below to create a Service Account and set up the integration.

1. Go into your project's settings and click on the Project ID link to open the associated GCP Project.
2. Click on the navigation menu in the GCP console, hover over "IAM & admin", and click "Service accounts".
3. Click on "+ CREATE SERVICE ACCOUNT", fill in the details, and give it the "Dialogflow Client API" role.
4. Click on "+ Create Key" and download the resulting JSON key file.
5. Save the JSON key file inside the functions/config directory of the cloned repo, else set the GOOGLE_APPLICATION_CREDENTIALS environmental variable on the deployment environment to the absolute path of Service Account JSON key file.
   See [this guide](https://cloud.google.com/dialogflow/docs/quick/setup#auth) for details.

### Creating a Slack app
Create a bot in a new Slack Workspace
1. Create or Sign in to Slack<br>
   Create/Sign in to a different slack workspace where we will create our bot and choose a workspace name for it.
2. Create a Slack app<br>
   Let’s check out the slack api https://api.slack.com/apps.
   Go to Create New App, select Create from scratch and give it a name.<br>
3. Add Bot Permissions<br>
   Go to OAuth & Permissions tab and scroll down to Scopes.<br>
   Now add all permissions your bot will need access to. Read the description carefully for all the scopes you’re providing the bot access to.<br>
   Do include app_mentions:read, channels:history, channels:join, chat:write, incoming-webhook<br>
4. Enable event subscription<br>
   On the left, go to the event subscriptions and enable events.<br>
   Now select all events we want to subscribe to from Subscribe to Bot Events. 
   The event subscription will ensure slack sends the events when they occur to the link provided.<br>
   We need to keep the link empty for now.
   Slack authorizes the link we provide, by sending a request with a challenge parameter and the app must respond with the challenge parameter.

### Setup Slack

The integration requires slack credentials from the slack api to function properly.<br>
Follow the steps to obtain the credentials and setup the config/token.go file to deploy and start the integration:<br>
1. In the Slack API, go to the basic information for your app and install the app to your workspace.
2. On installing the app to the workspace, you should be able to see a token in OAuth & Permissions. This is your BOT_TOKEN.
3. In basic information for the app, there should be one APP_TOKEN. Copy both and replace the values in config/token.go
4. On slack, go to the channel(s) you want the slack bot to have access to and invite the bot to the channel. Refer the below image for adding the bot to a channel. Alternatively, you can type /invite on the channel


## Deploying via App Engine

### Setup

1. Go into the project's settings and click on the Project ID link to open its associated GCP Project.
2. Click on the navigation menu in the GCP console and click "Billing". Set up and enable billing for the project.
3. Enable Google App Engine or Google Cloud Functions, according to what you want to use, for the project
   [here](https://console.cloud.google.com/flows/enableapi?apiid=cloudbuild.googleapis.com,run.googleapis.com).
4. Clone this git repository onto your local machine or development environment:
   `git clone [repository url]`
5. Open the root directory of the repository on your local machine or development environment.

### Changes in the app.yaml file

Open the app.yaml(https://github.com/Sampriti-Mitra/dialogflow-slack-sdk/blob/main/app.yaml) in the root directory of the repository, and uncomment the last line when using socket mode

If you have not done so already, copy your Service Account JSON key file to the desired subdirectory.


### Deploying App Engine
1. On the terminal, cd to the root directory of the cloned project and `gcloud app deploy --project [project-name]`
   This will deploy your project.
2. To check the logs, ` gcloud app --project [project-name] logs tail -s default`
3. Plugin the url obtained from 2 into slack's event url. Request URL should get approved if the app was able to successfully respond back with the challenge parameter. 
   This basically meant slack sent your URL some request, and you needed to respond with the challenge parameter, which you did!
4. On slack, Ensure the bot events you need to subscribe to, are all selected. If not, then add and save them.
   Include app_mention and message.im
   
## Post-deployment

### Shutting Down an Integration

In order to shut down an integration set up via the steps in this README, you need to delete the entire project where app engine is hosted.