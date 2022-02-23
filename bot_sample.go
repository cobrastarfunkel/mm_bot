package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	BOTNAME = "bobbot"
	TOKEN   = "j4m37kzow78zprw5bwgtysz89a"
	MMURI   = "http://localhost:8065"
	WSURI   = "ws://localhost:8065"
)

var debuggingChannel *model.Channel

var logger = Logger{}

func main() {
	logger.init(DEBUG)
	logger.Info("Starting Bobbot")

	c := Conn{}
	c.init()

	b := MMBot{botName: BOTNAME, teamName: "botsample", conn: &c}
	// Lets find our bot team
	b.init()

	logger.Debug(fmt.Sprintf("Client: %v", c.client))
	logger.Debug(fmt.Sprintf("Bot: %v", b))
	// Lets create a bot channel for logging debug messages into
}
