package eventloops

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/twitch"
	"github.com/spf13/viper"
	"log"
)

var (
	discord                *discordgo.Session
	queuedDiscordListeners []interface{}
	queuedTwitchListeners  []interface{}
)

func StartLoops() {
	twitch.Run()

	err := discord.Open()
	if err != nil {
		log.Fatalln("Could not connect to discord: ", err)
	}

	Run()
}

func GetDiscordSession() *discordgo.Session {
	return discord
}

/*
 * Called by subscribers to register each function that should be triggered by the discord event loops.
 *  In the future, I would like this to be called from with the subscribing file, not a whole list here, but will do for now.
 */
func RegisterDiscordListener(listener interface{}) {
	// 	queuedDiscordListeners = append(queuedDiscordListeners, listener)
	discord.AddHandler(listener)
}

func RegisterTimerListener(Callback callback, period float64) {
	// 	queuedDiscordListeners = append(queuedDiscordListeners, listener)
	RegisterTimer(Callback, period)
}

/*
 * Called by subscribers to register each function that should be triggered by the twitch event loops.
 *  In the future, I would like this to be called from with the subscribing file, not a whole list here, but will do for now.
 */
func RegisterTwitchListener(listener interface{}) {
	// 	queuedTwitchListeners = append(queuedTwitchListeners, listener)
	twitch.RegisterHandler(listener)
}

func ConnectToDiscord() {
	log.Printf("Connecting to discord with bot id: '%s'\n", viper.GetString("BotID"))
	_discord, err := discordgo.New(viper.GetString("BotID"))
	if err != nil {
		log.Fatalln("Could not connect to discord: ", err)
	}

	discord = _discord

	log.Println("Discord event loop successfully initialized")
}
