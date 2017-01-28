package handlers

import (
	"discharm/cmd/eventloops"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	eventloops.RegisterDiscordListener(ManagementStatus)
}

/*
 * Get the current status of the bot and send in discord via channel message.
 */
func ManagementStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!statustest") {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Fatalln("There was an error getting the channel: ", err)
		}
		log.Printf("%s asked for the status in channel %s, we are obliging!", m.Message.Author.Username, channel.Name)
		s.ChannelMessageSend(m.ChannelID, "Discord Status:")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Shard: %d", s.ShardID))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Shard Count: %d", s.ShardCount))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Uptime: %d", s.ShardCount))
	}
}
