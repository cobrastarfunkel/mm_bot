package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Conn struct {
	client          *model.Client4
	webSocketClient *model.WebSocketClient
}

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func (conn *Conn) SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if conn.webSocketClient != nil {
				conn.webSocketClient.Close()
			}

			conn.SendMsg("_"+BOTNAME+" has **stopped** running_", "", debuggingChannel.Id)
			os.Exit(0)
		}
	}()
}

func (conn *Conn) SetupClient() {
	conn.client = model.NewAPIv4Client(MMURI)
	conn.client.SetToken(TOKEN)
	logger.Debug("Conn setup")
}

func (conn Conn) MakeSureServerIsRunning() {
	if props, resp := conn.client.GetOldClientConfig(""); resp.Error != nil {
		logger.Error("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		logger.PrintError(resp.Error)
		os.Exit(1)
	} else {
		logger.Info("Server detected and is running version " + props["Version"])
	}
}

func (conn Conn) CreateBotDebuggingChannelIfNeeded(channelName string, botTeamId string) {
	if rchannel, resp := conn.client.GetChannelByName(channelName, botTeamId, ""); resp.Error != nil {
		logger.Error("We failed to get the channels")
		logger.PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = channelName
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = botTeamId
	if rchannel, resp := conn.client.CreateChannel(channel); resp.Error != nil {
		logger.Error("We failed to create the channel " + channelName)
		logger.PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		logger.Error("Looks like this might be the first run so we've created the channel " + channelName)
	}
}

func (conn Conn) SendMsg(msg string, replyToId string, channelId string) {
	post := &model.Post{}
	post.ChannelId = channelId
	post.Message = msg

	post.RootId = replyToId

	if _, resp := conn.client.CreatePost(post); resp.Error != nil {
		logger.Error("We failed to send a message to the logging channel")
		logger.PrintError(resp.Error)
	}
}

func (conn Conn) HandleMsg(event *model.WebSocketEvent, b MMBot) {
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

		chn, _ := conn.client.GetChannel(event.Broadcast.ChannelId, "")
		logger.Debug(fmt.Sprintf("responding to %s channel msg", chn.Name))

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)(alive|up|running)(?:$|\W)`, message); matched {
			conn.SendMsg("Yes I'm running", post.Id, event.Broadcast.ChannelId)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello|hi(?:$|\W)`, message); matched {
			conn.SendMsg("Hi!", post.Id, event.Broadcast.ChannelId)
			return
		}
	}

	conn.SendMsg("I did not understand you!", post.Id, event.Broadcast.ChannelId)
}

func (conn *Conn) StartWebsocketListening(b MMBot) {
	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(WSURI, conn.client.AuthToken)
	if err != nil {
		logger.Error("We failed to connect to the web socket")
		logger.PrintError(err)
	}
	conn.webSocketClient = webSocketClient
	conn.webSocketClient.Listen()

	go func() {
		for resp := range conn.webSocketClient.EventChannel {
			conn.HandleMsg(resp, b)
		}
	}()

	// You can block forever with
	select {}
}

func (c *Conn) init() {
	c.SetupGracefulShutdown()
	c.SetupClient()

	// Lets test to see if the mattermost server is up and running
	c.MakeSureServerIsRunning()
}
