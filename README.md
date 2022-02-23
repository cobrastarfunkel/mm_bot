# Mattermost Bot Sample

## Overview

This sample Bot shows how to use the Mattermost [Go driver](https://github.com/mattermost/mattermost-server/blob/master/model/client4.go) to interact with a Mattermost server, listen to events and respond to messages. Documentation for the Go driver can be found [here](https://godoc.org/github.com/mattermost/mattermost-server/model#Client).

Highlights of APIs used in this sample:
 - Log in to the Mattermost server
 - Create a channel
 - Modify user attributes 
 - Connect and listen to WebSocket events for real-time responses to messages
 - Post a message to a channel

This Bot Sample was tested with Mattermost server version 3.10.0.

## Setup Server Environment

### Via Docker And Docker-Compose
1 - Ensure [Docker](https://www.docker.com/get-started) and [Docker-Compose](https://docs.docker.com/compose/install/) are installed for your system

2 - Run `docker-compose up -d --build` and the mattermost client will be built and will expose the port `8065` to your system's localhost

3 - Run `./add_users.sh`. The login information for the Mattermost client will be printed

4 - Start the Bot.
```
make run
```
You can verify the Bot is running when 
  - `Server detected and is running version X.Y.Z` appears on the command line.
  - `Mattermost Bot Sample has started running` is posted in the `Debugging For Sample Bot` channel.

See "Test the Bot" for testing instructions

## Test the Bot

1 - Log in to the Mattermost server as `bill@example.com` and `Password1!`.

2 - Join the `Debugging For Sample Bot` channel.

3 - Post a message in the channel such as `are you running?` to see if the Bot responds. You should see a response similar to `Yes I'm running` if the Bot is running.

## Stop the Bot

1 - In the terminal window, press `CTRL+C` to stop the bot. You should see `Mattermost Bot Sample has stopped running` posted in the `Debugging For Sample Bot` channel.
