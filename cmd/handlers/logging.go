package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"log"
)

func init() {
	afterinit.RegisterAfterInitializationFunction(LoggingRegisterListeners)
}

func LoggingRegisterListeners() {
	eventloops.RegisterDiscordListener(MessageCreateLog)
}

func MessageCreateLog(s *discordgo.Session, m *discordgo.MessageCreate) {
	if "273124631656005633" != m.Author.ID {
		log.Printf("%20s %20s > %s\n", m.ChannelID, m.Author.Username, m.Content)
	}
}
