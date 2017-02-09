package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"log"
	"runtime"
	"strings"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	afterinit.RegisterAfterInitializationFunction(ManagementRegisterListeners)
}

func ManagementRegisterListeners() {
	eventloops.RegisterDiscordListener(ManagementStatus)
	eventloops.RegisterDiscordListener(PingPing)
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
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Go Version: %s", runtime.Version()))
		memstats := runtime.MemStats{}
		runtime.ReadMemStats(&memstats)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Memory Allocated (mb): %.2f", float64(memstats.Alloc)/1024/1024))
	}
}

func PingPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!ping") {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
