package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	countdowns map[int]*Countdown
)

type Countdown struct {
	Id         int
	Channel    string
	Message    string
	Period     time.Duration
	End        time.Time
	LastNotify time.Time
}

/*
 * Called automatically by go when this file is included.
 */
func init() {
	/*
	 * We register a function to be called after common dependencies are initiated.
	 */
	afterinit.RegisterAfterInitializationFunction(CountdownRegisterListeners)
	countdowns = make(map[int]*Countdown)
}

/*
 * Registered in init function. Will be called after common dependencies are initiated.
 */
func CountdownRegisterListeners() {
	eventloops.RegisterDiscordListener(DiscordCountdownCommand)
	eventloops.RegisterTimerListener(TimerListener, 1) // Call TimerListener every second
}

/*
 * Regisered in CountdownRegisterListeners above and called on every new discord message.
 *
 * At the moment, the only commands processed by this function are:
 *   - Add countdowns   - !countdown start <channel-id> <message> <period> <countdown-end>
 *   - View countdowns  - !countdown list
 *   - Get Help         - !countdown help
 *   - Delete Countdown - !countdown delete <id>
 */
func DiscordCountdownCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!countdown") {
		if strings.HasPrefix(m.Message.Content, "!countdown start") {
			startCountdown(s, m)
		} else if strings.HasPrefix(m.Message.Content, "!countdown list") {

		} else if strings.HasPrefix(m.Message.Content, "!countdown delete") {

		} else if strings.HasPrefix(m.Message.Content, "!countdown help") {
			printCountdownHelp(s, m)
		}
	}
}

func startCountdown(s *discordgo.Session, m *discordgo.MessageCreate) {

	commandParts := strings.Fields(m.Message.Content)

	if len(commandParts) != 6 {
		printCountdownHelp(s, m)
		return
	}

	/*
	 * commandParts now contains:
	 *   [0] - !countdown
	 *   [1] - start
	 *   [2] - Channel ID
	 *   [3] - Message       - with %formatted-time% to be replaced with the time left
	 *   [4] - Countdown End - in the format ???
	 */
	channelId := commandParts[2]
	message := commandParts[3]
	periodString := commandParts[4]
	countdownEnd := commandParts[5]

	validChannelId, err := regexp.Match("^[0-9]{18}$", []byte(channelId))
	if err != nil {
		log.Fatalln("There was an error error checking the channel id")
	}

	if !validChannelId {
		s.ChannelMessageSend(m.ChannelID, "Please make sure that the channel id is 18 numbers")
		return
	}

	periodInt, err := strconv.Atoi(periodString)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please make sure that the period is a valid integer")
	}

	parsedLayout, err := time.Parse(viper.GetString("UserEnteredDateFormat"), countdownEnd)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please make sure that the date you put was in the format: DD/MM/YYYY-HH:MM:SS")
		log.Printf("Error parsing date: %s", err)
	}

	newCountdown := Countdown{}
	newCountdown.Id = len(countdowns) + 1
	newCountdown.Channel = m.ChannelID
	newCountdown.Message = message
	newCountdown.Period = time.Second * time.Duration(periodInt)
	newCountdown.End = parsedLayout
	newCountdown.LastNotify = time.Now()
	countdowns[newCountdown.Id] = &newCountdown
}

/*
 * post the help page for the countdown function of the bot to discord.
 */
func printCountdownHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	HelpMessage := "```" +
		"Discharm countdown usage: \n" +
		"   - Add countdowns   - !countdown start <channel-id> <message> <countdown-end> \n" +
		"   - View countdowns  - !countdown list \n" +
		"   - Get Help         - !countdown help \n" +
		"   - Delete Countdown - !countdown delete <id>" +
		"```"
	s.ChannelMessageSend(m.ChannelID, HelpMessage)
}

/*
 * Regisered in CountdownRegisterListeners above and called every 120 seconds.
 * Checks all countdowns for any needed notifications
 */
func TimerListener() {
	// In here, we will check all countdowns to check for any needed notifications
	for key, countdown := range countdowns {
		if time.Since(countdown.End) > 0 {
			delete(countdowns, key)
			return
		}

		if time.Since(countdown.LastNotify) > countdown.Period {
			s := eventloops.GetDiscordSession()
			timeLeft := countdown.End.Sub(time.Now())
			log.Printf("Time left: %s", timeLeft)
			s.ChannelMessageSend(countdown.Channel, strings.Replace(countdown.Message, viper.GetString("CountdownTimeRemainingPlaceholder"), fmt.Sprintf("%s", timeLeft), -1))
			countdown.LastNotify = time.Now()
		}
	}
}
