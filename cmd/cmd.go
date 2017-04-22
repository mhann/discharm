package cmd

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	_ "github.com/mhann/discharm/cmd/handlers"
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

	// /*
	//  * HORRIBLE! Workaround to allow us to inclue handlers (and therefore run their init functions)
	//  *   This is because go won't let us include without using - even though we ARE using init functions.
	//  *   It may be able to be achieved using _ <import> above - see here:
	//  *     http://stackoverflow.com/questions/21220077/what-does-an-underscore-in-front-of-an-import-statement-mean-in-golang
	//  */
	// handlers.DummyFunction()

	eventloops.ConnectToDiscord()
	eventloops.StartLoops()

	afterinit.RunAfterInitializationFunctions()

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
