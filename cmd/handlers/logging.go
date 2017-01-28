package handlers

import (
	"discharm/cmd/eventloops"
	"github.com/bwmarrin/discordgo"
	"log"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	eventloops.RegisterDiscordListener(MessageCreateLog)
}

func MessageCreateLog(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("%20s %20s > %s\n", m.ChannelID, m.Author.Username, m.Content)
}
