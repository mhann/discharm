package cmd

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"github.com/mhann/discharm/cmd/handlers"
	"log"
)

var (
	discordToken string // Set by command line flag later - used to authenticate with discord.
	discord      discordgo.Session
)

/*
 * Initialize the commandline flags and parse so they can be used later on.
 */
func initFlags() {
	flag.StringVar(&discordToken, "discordToken", "", "Discord application token")
	flag.Parse()
}

/*
 * Entrypoint into the program.
 */
func Run() {
	log.Println("Bot is now running.  Press CTRL-C to exit.")

	// HORRIBLE! Workaround to allow us to inclue handlers (and therefore run their init functions)
	handlers.DummyFunction()

	eventloops.ConnectToDiscord()
	eventloops.StartLoops()

	afterinit.RunAfterInitializationFunctions()

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
