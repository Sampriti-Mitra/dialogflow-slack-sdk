# dialogflow-slack-sdk

## Introduction
The goal of this guide is to show you how to set up an integration deployment to link your Dialogflow agent to slack.
If you don't already have a Dialogflow agent, you may create one by following the instructions here or by adding a prebuilt agent.
Although this integration deployment may be set up on any other hosting platform, these instructions will use Google's App Engine/Cloud functions.

## GCP Setup

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

### Setup Slack

The integration requires slack credentials from the slack api to function properly.

Follow the steps to obtain the credentials and setup the config/token.go file to deploy and start the integration:

## Post-deployment

### Shutting Down an Integration

In order to shut down an integration set up via the steps in this README, you need to delete the entire project where app engine is hosted.