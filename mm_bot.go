package main

import (
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
)

type MMBot struct {
	botUser  *model.User
	botName  string
	botTeam  *model.Team
	teamName string
}

func (b *MMBot) SetupBot(client model.Client4) {
	bot, res := client.GetUserByUsername(b.botName, "")

	if res.Error != nil {
		logger.Error("Error in SetupBot")
		logger.PrintError(res.Error)
		os.Exit(1)
	}
	b.botUser = bot
}

func (b *MMBot) FindBotTeam(client model.Client4) {
	if team, resp := client.GetTeamByName(b.teamName, ""); resp.Error != nil {
		logger.Error("We failed to get the initial load")
		logger.Error("or we do not appear to be a member of the team '" + b.teamName + "'")
		logger.PrintError(resp.Error)
		os.Exit(1)
	} else {
		b.botTeam = team
	}
}

func (b *MMBot) init(c model.Client4) {
	b.SetupBot(c)
	b.FindBotTeam(c)
}
