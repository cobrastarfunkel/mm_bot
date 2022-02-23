package main

import (
	"os"
	"os/signal"

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

func (c *Conn) init() {
	c.SetupGracefulShutdown()
	c.SetupClient()

	// Lets test to see if the mattermost server is up and running
	c.MakeSureServerIsRunning()
}
