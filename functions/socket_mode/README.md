# Socket Mode
Socket Mode allows your app to communicate with Slack via a WebSocket URL. WebSockets use a bidirectional stateful protocol with low latency to communicate between two parties—in this case, Slack and your app. <br>

Unlike a public HTTP endpoint, the WebSocket URL you listen to is not static. It's created at runtime by calling the apps.connections.open method, and it refreshes regularly.<br>

Because the URL isn't static and is created at runtime, it allows for greater security in some cases, and it allows you to develop behind a firewall.<br>

In Socket Mode, your app still uses the very same Events API and interactive components of the Slack platform. The only difference is the communication protocol.<br>


### Setup Slack

Follow the Slack setup instructions in the root ReadME of the project.
In the Slack Api bot page, select the Socket Mode on the left tab. Click on enable socket mode.


## Deploying via App Engine

### Setup

Follow the App Engine set up instructions in the root ReadME of the project.

### Changes in the app.yaml file

1. Open the app.yaml(https://github.com/Sampriti-Mitra/dialogflow-slack-sdk/blob/main/app.yaml) in the root directory of the repository.
2. Uncomment the last line when using socket mode

If you have not done so already, copy your Service Account JSON key file to the desired subdirectory and export the credentials.


### Deploying App Engine
1. On the terminal, cd to the root directory of the cloned project and `gcloud app deploy --project [project-name]`
   This will deploy your project.
2. To check the logs, ` gcloud app --project [project-name] logs tail -s default`

3. Don't forget to copy the APP_TOKEN in basic information for the app and replace the value in config/token.go

On slack, go to the channel(s) you want the slack bot to have access to and invite the bot to the channel. Refer the below image for adding the bot to a channel. Alternatively, you can type /invite on the channel

## Post-deployment

### Shutting Down an Integration

In order to shut down an integration set up via the steps in this README, you need to delete the entire project where app engine is hosted.