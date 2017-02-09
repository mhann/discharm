package handlers

import (
	"fmt"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"github.com/mhann/discharm/twitch"
	"github.com/spf13/viper"
	"log"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	afterinit.RegisterAfterInitializationFunction(TwitchRegisterListeners)
}

type twitchChannelCheck struct {
	Name   string `mapstructure:"channel"`
	Notify string `mapstructure:"notify"`
}

var (
	channelTwitchChannelChecks []twitchChannelCheck
)

func TwitchRegisterListeners() {
	eventloops.RegisterTwitchListener(ChannelLive)
	eventloops.RegisterTwitchListener(ChannelOffline)

	err := viper.UnmarshalKey("TwitchChannelChecks", &channelTwitchChannelChecks)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println("Registering channels from configuration file")

	for _, twitchChannelCheck := range channelTwitchChannelChecks {
		twitch.RegisterChannelOnlineCheck(twitchChannelCheck.Name)
		log.Printf("Registering channel online check for channel %s and will notify: %s", twitchChannelCheck.Name, twitchChannelCheck.Notify)
	}
}

/*
 * Called whenever a twitch channel that is checked (see RegisterChannelOnlineCheck above) goes online.
 */
func ChannelLive(channel *twitch.ChannelOnline) {
	channelNotifyNeeded, ChannelToNotify := getDiscordChannelToNotify(channel.Name)

	if channelNotifyNeeded {
		discord := eventloops.GetDiscordSession()
		log.Println("Twitch channel has come online!")
		discord.ChannelMessageSend(ChannelToNotify, fmt.Sprintf("@everyone '%s' just started streaming. Watch live at: https://twitch.tv/%s", channel.Name, channel.Name))

		// embed := &discordgo.MessageEmbed{}
		// embed.URL = fmt.Sprintf("https://twitch.tv/%s", channel.Name)
		// embed.Title = fmt.Sprintf("%s just went live", channel.Name)
		// embed.Description = fmt.Sprintf("%s just went live", channel.Name)
		// embedImage := &discordgo.MessageEmbedImage{}
		// embedImage.URL = "http://assets.barcroftmedia.com.s3-website-eu-west-1.amazonaws.com/assets/images/recent-images-11.jpg"
		// embedFooter := &discordgo.MessageEmbedFooter{}
		// embedFooter.Text = "Powered by discharm"
		// embed.Image = embedImage
		// embed.Footer = embedFooter

		// discord.ChannelMessageSendEmbed("188699044456038400", embed)
	}
}

func getDiscordChannelToNotify(channelName string) (bool, string) {
	for _, channel := range channelTwitchChannelChecks {
		if channel.Name == channelName {
			return true, channel.Notify
		}
	}

	return false, ""
}

/*
 * Called whenever a twitch channel that is checked (see RegisterChannelOnlineCheck above) goes offline.
 */
func ChannelOffline(channel *twitch.ChannelOffline) {
	discord := eventloops.GetDiscordSession()
	log.Println("Twitch channel has come offline!")
	discord.ChannelMessageSend("188699044456038400", fmt.Sprintf("'%s' just finished streaming. Watch their next stream at: https://twitch.tv/%s", channel.Name, channel.Name))
}
