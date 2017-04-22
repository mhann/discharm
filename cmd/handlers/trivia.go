package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mhann/discharm/cmd/afterinit"
	"github.com/mhann/discharm/cmd/eventloops"
	"html"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

type Trivia struct {
	Id              int
	Channel         string
	Test            int
	Scores          map[string]int
	QuestionLineup  []question
	CurrentQuestion question
	Mutex           *sync.Mutex
}

type opentdbResponse struct {
	ResponseCode int        `json:"response_code"`
	Results      []question `json:"results"`
}

type question struct {
	Category         string   `json:"category"`
	CorrectAnswer    string   `json:"correct_answer"`
	Difficulty       string   `json:"difficulty"`
	IncorrectAnswers []string `json:"incorrect_answers"`
	Question         string   `json:"question"`
	Type             string   `json:"type"`
}

var (
	runningTrivias map[string]*Trivia
	db             *sql.DB
)

func getRunningTrivias() ([]*Trivia, error) {
	log.Println("Getting running trivias from the database")
	rows, err := db.Query("select Id, Channel, CurrentQuestion from Trivia")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	runningTrivias := []*Trivia{}

	for rows.Next() {
		runningTrivia := Trivia{}
		err := rows.Scan(&runningTrivia.Id, &runningTrivia.Channel, &runningTrivia.Test)
		if err != nil {
			return nil, err
		}
		log.Printf("%s", runningTrivia.Channel)
		runningTrivias = append(runningTrivias, &runningTrivia)
	}

	return runningTrivias, nil
}

/*
 * Called automatically by go when this file is included.
 */
func init() {
	afterinit.RegisterAfterInitializationFunction(TriviaRegisterListeners)
}

func TriviaRegisterListeners() {
	// 	eventloops.RegisterTimerListener(timer, 10)
	eventloops.RegisterDiscordListener(discordMessage)
	runningTrivias = make(map[string]*Trivia)

	localDb, err := sql.Open("mysql", "root:J@spercat@tcp(127.0.0.1:3306)/discharmdev")
	if err != nil {
		log.Fatal(err)
	}

	err = localDb.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db = localDb

	_, err = getRunningTrivias()
	if err != nil {
		log.Fatalln(err)
	}
}

func discordMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*
	 * Temporary check to make sure that we are not processing our own messages.
	 */
	if "273124631656005633" != m.Author.ID {
		lowerMessage := strings.ToLower(m.Message.Content)

		if strings.HasPrefix(lowerMessage, "!trivia start") {
			log.Println("Starting trivia")
			startTrivia(s, m)
			return
		} else if strings.HasPrefix(lowerMessage, "!trivia help") {
			printHelp(s, m)
			return
		} else if strings.HasPrefix(lowerMessage, "!trivia running") {
			sendRunningTrivias(s, m)
			return
		} else {
			/*
			 * Should only really be run if trivia is running in channel.
			 */
			processAnswer(s, m.ChannelID, lowerMessage, m.Author)
		}
	}
}

func sendRunningTrivias(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Currently Running Trivias:")

	for channel, trivia := range runningTrivias {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Channel: %s, Players: %d", channel, len(trivia.Scores)))
	}
}

func printHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Trivia Commands:")
	s.ChannelMessageSend(m.ChannelID, " - !trivia start - start a game of trivia")
	s.ChannelMessageSend(m.ChannelID, " - !trivia help - show this help")
	s.ChannelMessageSend(m.ChannelID, "To answer, just type the answer in the chat (not including numbers)")
}

func startTrivia(s *discordgo.Session, m *discordgo.MessageCreate) {
	trivia := Trivia{}
	trivia.Channel = m.ChannelID
	trivia.Scores = make(map[string]int)
	trivia.QuestionLineup = getQuestions(10).Results
	trivia.Mutex = &sync.Mutex{}
	runningTrivias[trivia.Channel] = &trivia
	sendNextTriviaQuestion(s, m.ChannelID)
}

func processAnswer(discordsession *discordgo.Session, channelID string, answer string, user *discordgo.User) {
	if trivia, ok := runningTrivias[channelID]; ok {
		trivia.Mutex.Lock()
		question := trivia.CurrentQuestion
		log.Printf(">  Incorrect answers at process answer: %d", len(question.IncorrectAnswers))
		printArrayOfStrings(question.IncorrectAnswers)
		if answer == strings.ToLower(question.CorrectAnswer) {
			discordsession.ChannelMessageSend(channelID, "Correct!")

			if _, ok := trivia.Scores[user.ID]; ok {
				log.Println("Incrementing score for user")
				trivia.Scores[user.ID] = trivia.Scores[user.ID] + 1
			} else {
				trivia.Scores[user.ID] = 1
			}

			discordsession.ChannelMessageSend(channelID, fmt.Sprintf("%s's score is now %d", user.Username, trivia.Scores[user.ID]))

			sendNextTriviaQuestion(discordsession, channelID)
		} else if isInAnswerList(answer, append(trivia.CurrentQuestion.IncorrectAnswers, trivia.CurrentQuestion.CorrectAnswer)) {
			discordsession.ChannelMessageSend(channelID, "Incorrect!")
			sendNextTriviaQuestion(discordsession, channelID)
		}
		trivia.Mutex.Unlock()
	}
}

func isInAnswerList(answer string, answers []string) bool {
	log.Printf("checking %s against %d possible answers", answer, len(answers))
	for _, possibleAnswer := range answers {
		log.Printf("Checking %s againt %s", answer, strings.ToLower(possibleAnswer))
		if answer == strings.ToLower(possibleAnswer) {
			return true
		}
	}

	return false
}

func printArrayOfStrings(toPrintStrings []string) {
	for _, toPrintString := range toPrintStrings {
		log.Printf(">  %s", toPrintString)
	}
}

func sendNextTriviaQuestion(discordsession *discordgo.Session, channelID string) {
	log.Printf(">  ")
	log.Printf(">  ------------ Asking trivia question -------------------")

	trivia := runningTrivias[channelID]

	question := trivia.QuestionLineup[0]
	log.Printf(">  Trivia question lineup before pop: %d", len(trivia.QuestionLineup))
	trivia.QuestionLineup = trivia.QuestionLineup[1:]
	log.Printf(">  Trivia question lineup after pop: %d", len(trivia.QuestionLineup))

	if len(trivia.QuestionLineup) == 0 {
		trivia.QuestionLineup = getQuestions(10).Results
	}

	discordsession.ChannelMessageSend(channelID, question.Question)

	answers := []string{}
	for _, incorrectAnswer := range question.IncorrectAnswers {
		answers = append(answers, incorrectAnswer)
	}

	log.Printf(">  Incorrect answers before appending: %d", len(question.IncorrectAnswers))
	log.Printf(">  All answers before appending: %d", len(answers))
	printArrayOfStrings(question.IncorrectAnswers)
	answers = append(answers, question.CorrectAnswer)
	log.Printf(">  Incorrect answers after appending: %d", len(question.IncorrectAnswers))
	log.Printf(">  All answers after appending: %d", len(answers))
	printArrayOfStrings(question.IncorrectAnswers)
	log.Printf(">  Correct Answer: '%s'", question.CorrectAnswer)
	log.Printf(">  Total Possible Answers: '%d'", len(answers))

	for i := range answers {
		j := rand.Intn(i + 1)
		answers[i], answers[j] = answers[j], answers[i]
	}

	log.Printf(">  Sort done")
	log.Printf(">  Incorrect answers after sort: %d", len(question.IncorrectAnswers))
	printArrayOfStrings(question.IncorrectAnswers)

	for index, answer := range answers {
		discordsession.ChannelMessageSend(channelID, fmt.Sprintf("%d. %s", index+1, answer))
		log.Printf(">  Possible Answer: '%s'", answer)
	}

	log.Printf(">  Incorrect answers after print: %d", len(question.IncorrectAnswers))
	printArrayOfStrings(question.IncorrectAnswers)

	trivia.CurrentQuestion = question

	log.Printf(">  Incorrect answers after function: %d", len(question.IncorrectAnswers))
	printArrayOfStrings(question.IncorrectAnswers)
	log.Printf(">  -------------------------------------------------------")
	log.Printf(">  ")
}

func getQuestions(count int) opentdbResponse {
	url := "https://opentdb.com/api.php?amount=10"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error talking to otdb api: ", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error talking to otdb api: ", err)
	}

	var response opentdbResponse
	respBuf, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error getting tvdb response: ", err)
	}

	// log.Println(string(respBuf))

	respBuf = jsonEncodeString(respBuf)
	respBuf = []byte(html.UnescapeString(string(respBuf)))

	err = json.Unmarshal(respBuf, &response)
	if err != nil {
		log.Fatalln("Error decoding json: ", err)
	}

	return response
}

/*
 * HORRIBLE way to replace html characters in the string.
 */
func jsonEncodeString(input []byte) []byte {
	input = bytes.Replace(input, []byte("&quot;"), []byte("\\\""), -1)
	input = bytes.Replace(input, []byte("&#039;"), []byte("'"), -1)
	return input
}
