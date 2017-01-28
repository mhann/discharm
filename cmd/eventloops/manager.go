package eventloops

import (
	"github.com/mhann/discharm/twitch"
	"github.com/bwmarrin/discordgo"
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

/*
 * Called by subscribers to register each function that should be triggered by the twitch event loops.
 *  In the future, I would like this to be called from with the subscribing file, not a whole list here, but will do for now.
 */
func RegisterTwitchListener(listener interface{}) {
	// 	queuedTwitchListeners = append(queuedTwitchListeners, listener)
	twitch.RegisterHandler(listener)
}

func init() {
	_discord, err := discordgo.New("marcus@hannmail.co.uk", "J@spercat")
	if err != nil {
		log.Fatalln("Could not connect to discord: ", err)
	}

	discord = _discord

	log.Println("Discord event loop successfully initialized")
	// 	registerQueuedHandlers()
}

// func registerQueuedHandlers() {
// 	log.Println("Registering discord event handlers")
// 	for handler := range queuedDiscordListeners {
// 		log.Println("Registering a discord event handler")
// 		discord.AddHandler(handler)
// 	}

// 	log.Println("Registing twitch event handlers")
// 	for handler := range queuedTwitchListeners {
// 		log.Println("Registering a twitch event handler")
// 		discord.AddHandler(handler)
// 	}
// }
