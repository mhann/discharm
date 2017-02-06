package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/eventloops"
	"github.com/mhann/discharm/twitch"
	"log"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	eventloops.RegisterTwitchListener(ChannelLive)
	eventloops.RegisterTwitchListener(ChannelOffline)
	twitch.RegisterChannelOnlineCheck("marcushann")
	twitch.RegisterChannelOnlineCheck("mrsyncz")
}

/*
 * Called whenever a twitch channel that is checked (see RegisterChannelOnlineCheck above) goes online.
 */
func ChannelLive(channel *twitch.ChannelOnline) {
	discord := eventloops.GetDiscordSession()
	log.Println("Twitch channel has come online!")
	discord.ChannelMessageSend("188699044456038400", fmt.Sprintf("@everyone '%s' just started streaming. Watch live at: https://twitch.tv/%s", channel.Name, channel.Name))

	embed := &discordgo.MessageEmbed{}
	embed.URL = fmt.Sprintf("https://twitch.tv/%s", channel.Name)
	embed.Title = fmt.Sprintf("%s just went live", channel.Name)
	embed.Description = fmt.Sprintf("%s just went live", channel.Name)
	embedImage := &discordgo.MessageEmbedImage{}
	embedImage.URL = "http://assets.barcroftmedia.com.s3-website-eu-west-1.amazonaws.com/assets/images/recent-images-11.jpg"
	embedFooter := &discordgo.MessageEmbedFooter{}
	embedFooter.Text = "Powered by discharm"
	embed.Image = embedImage
	embed.Footer = embedFooter

	discord.ChannelMessageSendEmbed("188699044456038400", embed)
}

/*
 * Called whenever a twitch channel that is checked (see RegisterChannelOnlineCheck above) goes offline.
 */
func ChannelOffline(channel *twitch.ChannelOffline) {
	discord := eventloops.GetDiscordSession()
	log.Println("Twitch channel has come offline!")
	discord.ChannelMessageSend("188699044456038400", fmt.Sprintf("'%s' just finished streaming. Watch their next stream at: https://twitch.tv/%s", channel.Name, channel.Name))
}
