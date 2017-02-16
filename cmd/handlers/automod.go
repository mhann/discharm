package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"log"
)

func init() {
	afterinit.RegisterAfterInitializationFunction(loggingRegisterListeners)
}

func loggingRegisterListeners() {
	eventloops.RegisterDiscordListener(messageCreateLog)
}

func messageCreateLog(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Deleting message!")
	// err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	// if err != nil {
	// 	log.Printf("Error deleting channel message: %s", err)
	// }
}
