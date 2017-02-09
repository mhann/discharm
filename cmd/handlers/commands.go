package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
)

/*
 * Called automatically by go when this file is included.
 */
func init() {
	afterinit.RegisterAfterInitializationFunction(CommandsRegisterListeners)
}

func CommandsRegisterListeners() {
	eventloops.RegisterDiscordListener(PrintCommands)
}

func PrintCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if "273124631656005633" != m.Author.ID {
		//log.Printf("Commands are: ")
	}
}
