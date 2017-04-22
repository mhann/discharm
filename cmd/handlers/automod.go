package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/database"
	"github.com/mhann/discharm/cmd/eventloops"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Timeout struct {
	UserID   string
	Start    time.Time
	Duration time.Duration
}

var (
	timeouts     map[int]Timeout
	timeoutsLock *sync.Mutex
)

func init() {
	afterinit.RegisterAfterInitializationFunction(loggingRegisterListeners)
	timeoutsLock = &sync.Mutex{}
}

func loggingRegisterListeners() {
	eventloops.RegisterDiscordListener(handleTimeoutMessage)
}

func handleTimeoutMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!timeout") {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			log.Printf("Couldn't get channel from channel id: #s", err)
		}

		userCanTimeout, err := canUserTimeout(s, m.Author, m.ChannelID, channel.GuildID)
		if err != nil {
			log.Printf("Couldn't determine if user can set timeouts: #s", err)
		}

		if userCanTimeout {
			// Process the timeout request
			processTimeoutRequest(s, m)
			log.Printf("Processing the timeout request")
			return
		}
	}

	timedout, err := isUserTimedout(m.Author, m.ChannelID)
	if err != nil {
		log.Println("Unable to tell if user is timedout, assuming not")
		return
	}

	if timedout {
		log.Printf("Deleting message!")
		// checkUserIsAuthorized(s, m, m.Author.ID)
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Printf("Error deleting channel message: %s", err)
		}
	}

}

func processTimeoutRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	db = database.GetConnection()
	commandParts := strings.Fields(m.Message.Content)

	/*
	 * @FIXME MAH:
	 *   We need to check if the array commandParts is the correct length here otherwise we could get a SEGV below.
	 */
	if len(commandParts) != 3 {
		log.Println("There were not enough parts to the timeout command - giving up.")
		printTimeoutHelp(s, m)
		return
	}

	/*
	 * Order of parts is:
	 *   1 - !timeout (ignored as this is always static)
	 *   2 - <@[discord-id]>
	 *   3 - timeout time in seconds
	 */
	toTimeoutIdString := commandParts[1]
	timeoutLengthString := commandParts[2]

	log.Printf("Raw user to timeout: %s", toTimeoutIdString)
	log.Printf("Raw time to timeout: %s", timeoutLengthString)

	/*
	 * Check that timeout user is in correct format, and then process.
	 */
	validRawUserId, err := regexp.Match("^<@[0-9]{18}>$", []byte(toTimeoutIdString))
	if err != nil {
		log.Fatalln("There was an error error checking the format of the raw user id")
	}

	if !validRawUserId {
		/*
		 * We could probably do something more useful here.
		 */
		s.ChannelMessageSend(m.ChannelID, "The user format was in the wrong format.")
		return
	}

	/*
	 * Carve the middle part of the id string out.
	 */
	processedTimeoutIdString := toTimeoutIdString[2:20]

	log.Printf("Processed user id to timeout: %s", processedTimeoutIdString)

	timeoutLengthInt, err := strconv.Atoi(timeoutLengthString)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please make sure that the timeout length is a valid integer")
	}

	log.Printf("Processed timeout length: %d", timeoutLengthInt)

	/*
	 * We have now validated the input of the user to the best of our ability.
	 *   Next step is to add the timeout to the database.
	 */
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO timeouts (target_user_id, creator_user_id, guild_id, created, length_seconds, known_expired) VALUES ($1, $2, $3, NOW(), $4, false);")
	if err != nil {
		log.Println(err)
		return
	}
	defer stmt.Close()

	channel, _ := s.Channel(m.ChannelID)

	_, err = stmt.Exec(processedTimeoutIdString, m.Author.ID, channel.GuildID, timeoutLengthInt)
	if err != nil {
		log.Println(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	/*
	 * Finally, we wish to inform the user that they have been timed out in a PM.
	 */
}

func printTimeoutHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	HelpMessage := "```" +
		"Discharm timeout usage: \n" +
		"   - Cooldown User   - !cooldown @<UserId> <timeout time in seconds> \n" +
		// "   - View countdowns  - !countdown list \n" +
		// "   - Get Help         - !countdown help \n" +
		// "   - Delete Countdown - !countdown delete <id>" +
		"```"
	s.ChannelMessageSend(m.ChannelID, HelpMessage)
}

func canUserTimeout(s *discordgo.Session, user *discordgo.User, channel_id string, guild_id string) (bool, error) {
	guildMember, err := s.GuildMember(guild_id, user.ID)
	if err != nil {
		log.Printf("Error getting guild member from user: %s", err)
		return false, err
	}

	for role := range guildMember.Roles {
		log.Println(role)
	}

	return true, nil
}

func isUserTimedout(user *discordgo.User, channel_id string) (bool, error) {
	db = database.GetConnection()

	rows, err := db.Query("SELECT created, length_seconds FROM timeouts WHERE target_user_id = $1 AND known_expired != true AND guild_id = $2", user.ID, channel_id)
	if err != nil {
		log.Println("Failed to get current timeouts from database")
		log.Println(err)
		return false, err
	}

	log.Println("Checking if user is in timedout table")
	for rows.Next() {
		log.Println("Checking next row")
		var created time.Time
		var length_seconds int

		err = rows.Scan(&created, &length_seconds)
		if err != nil {
			log.Println("Error scanning rows")
			return false, err
		}

		log.Println(created)
		log.Println(length_seconds)
		log.Println(created.Add(time.Second * time.Duration(length_seconds)))
		log.Println(time.Now())

		if time.Now().Before(created.Add(time.Second * time.Duration(length_seconds))) {
			log.Println("User is currently timedout")
			return true, nil
		} else {
			/*
			 * @FIXME:
			 *   This needs to be set known_expired to true.
			 */
		}
	}

	return false, nil
}

func checkUserIsAuthorized(s *discordgo.Session, m *discordgo.MessageCreate, UserID string) (bool, error) {
	channel, err := s.Channel(m.ChannelID)
	// Check for errors
	if err != nil {
		log.Printf("Error getting the channel from channel ID '%s': %s", m.ChannelID, err)
	}

	member, err := s.State.Member(channel.GuildID, UserID)
	// Check for errors
	if err != nil {
		log.Printf("Error getting the member from userID '%s' and GuildID '%s': %s", UserID, channel.GuildID, err)
	}

	guildRoles, err := s.GuildRoles(channel.GuildID)
	// Check for errors
	if err != nil {
		log.Printf("Error getting list of roles for GuildID '%s': %s", channel.GuildID, err)
	}

	for _, role := range guildRoles {
		if stringInSlice(role.ID, member.Roles) {

		}
		log.Printf("Role id: %s", role.ID)
	}

	return false, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
