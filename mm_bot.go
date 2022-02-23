package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

type MMBot struct {
	botUser  *model.User
	botName  string
	botTeam  *model.Team
	teamName string
	conn     *Conn
}

func (b *MMBot) SetupBot() {
	bot, res := b.conn.client.GetUserByUsername(b.botName, "")

	if res.Error != nil {
		logger.Error("Error in SetupBot")
		logger.PrintError(res.Error)
		os.Exit(1)
	}
	b.botUser = bot
}

func (b *MMBot) FindBotTeam() {
	if team, resp := b.conn.client.GetTeamByName(b.teamName, ""); resp.Error != nil {
		logger.Error("We failed to get the initial load")
		logger.Error("or we do not appear to be a member of the team '" + b.teamName + "'")
		logger.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.botTeam = team
	}
}

func (b MMBot) HandleMsg(event *model.WebSocketEvent) {
	// Lets only reponded to messaged posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))

	if post != nil {
		logger.Debug(fmt.Sprintf("Post ID: %s Bot Id: %s", post.UserId, b.botUser.Id))
		message := strings.ToLower(post.Message)

		if !strings.Contains(message, "@"+b.botUser.Username) {
			return
		}

		// ignore my events
		if post.UserId == b.botUser.Id {
			return
		}

		chn, _ := b.conn.client.GetChannel(event.Broadcast.ChannelId, "")
		logger.Debug(fmt.Sprintf("responding to %s channel msg", chn.Name))

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)(alive|up|running)(?:$|\W)`, message); matched {
			b.conn.SendMsg("Yes I'm running", post.Id, event.Broadcast.ChannelId)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello|hi(?:$|\W)`, message); matched {
			b.conn.SendMsg("Hi!", post.Id, event.Broadcast.ChannelId)
			return
		}
	}

	b.conn.SendMsg("I did not understand you!", post.Id, event.Broadcast.ChannelId)
}

func (b MMBot) CreateBotDebuggingChannelIfNeeded(channelName string) {
	botTeamId := b.botTeam.Id

	if rchannel, resp := b.conn.client.GetChannelByName(channelName, botTeamId, ""); resp.Error != nil {
		logger.Error("We failed to get the channels")
		logger.PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		return
	}

	channel := &model.Channel{
		Name:        channelName,
		DisplayName: "Debugging For Sample Bot",
		Purpose:     "This is used as a test channel for logging bot debug messages",
		Type:        model.CHANNEL_OPEN,
		TeamId:      botTeamId,
	}

	if rchannel, resp := b.conn.client.CreateChannel(channel); resp.Error != nil {
		logger.Error("We failed to create the channel " + channelName)
		logger.PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		logger.Error("Looks like this might be the first run so we've created the channel " + channelName)
	}
}

func (b *MMBot) StartWebsocketListening() {
	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(WSURI, b.conn.client.AuthToken)
	if err != nil {
		logger.Error("We failed to connect to the web socket")
		logger.PrintError(err)
	}
	b.conn.webSocketClient = webSocketClient
	b.conn.webSocketClient.Listen()

	go func() {
		for resp := range b.conn.webSocketClient.EventChannel {
			b.HandleMsg(resp)
		}
	}()

	// You can block forever with
	select {}
}

func (b *MMBot) init() {
	b.SetupBot()
	b.FindBotTeam()
	b.CreateBotDebuggingChannelIfNeeded("debugging-for-sample-bot")
	b.conn.SendMsg("_"+BOTNAME+" has **started** running_", "", debuggingChannel.Id)

	b.StartWebsocketListening()
}
