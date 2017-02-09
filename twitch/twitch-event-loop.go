package twitch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	listeners []interface{}
	handlers  = map[string][]*eventHandlerInstance{}
	channels  []*checkChannel
)

// eventHandlerInstance is a wrapper around an event handler, as functions
// cannot be compared directly.
type eventHandlerInstance struct {
	eventHandler EventHandler
}

type checkChannel struct {
	Name       string
	LastStatus bool
	FirstCheck bool
}

type twitchChannel struct {
	Mature              bool   `json:"mature"`
	Status              string `json:"status"`
	BroadcasterLanguage string `json:"broadcaster_language"`
	Game                string `json:"game"`
	Language            string `json:"language"`
	Name                string `json:"name"`
	Url                 string `json:"url"`
}

type twitchStreamChannel struct {
	Stream *twitchStream `json:"stream"`
}

type twitchStream struct {
	Channel twitchChannel `json:"channel"`
}

type EventHandler interface {
	// Type returns the type of event this handler belongs to.
	Type() string

	// Handle is called whenever an event of Type() happens.
	// It is the recievers responsibility to type assert that the interface
	// is the expected struct.
	Handle(interface{})
}

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(*ChannelOnline):
		return channelOnlineEventHandler(v)
	case func(*ChannelOffline):
		return channelOfflineEventHandler(v)
	}

	return nil
}

const (
	EventChannelOnline  = iota
	EventChannelOffline = iota
)

func RegisterHandler(handler interface{}) func() {
	eh := handlerForInterface(handler)

	if eh == nil {
		log.Println("Invalid handler type registered, handler will never be called")
		return func() {}
	}

	return addEventHandler(eh)
}

func addEventHandler(eventHandler EventHandler) func() {
	log.Println("Handler is for type:")
	ehi := &eventHandlerInstance{eventHandler}

	handlers[eventHandler.Type()] = append(handlers[eventHandler.Type()], ehi)
	log.Println("Twitch event handler added")

	return func() {
		return
	}
}

func handle(t string, i interface{}) {
	log.Println("Running handlers for: ", t)
	for _, eh := range handlers[t] {
		log.Println("Running a handler")
		go eh.eventHandler.Handle(i)
	}
}

func checkTwitchStream(channelName string) twitchStreamChannel {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/streams/"+channelName, nil)
	if err != nil {
		log.Fatalln("Error talking to twitch api: ", err)
	}

	req.Header.Add("Client-ID", `sshkf19j8b7may6p14i92898a88rtj`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error talking to twitch api: ", err)
	}

	var stream twitchStreamChannel
	respBuf, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error decoding json: ", err)
	}

	// log.Println(string(respBuf))
	err = json.Unmarshal(respBuf, &stream)
	if err != nil {
		log.Fatalln("Error decoding json: ", err)
	}

	return stream
}

func RegisterChannelOnlineCheck(channelName string) {
	log.Printf("Possibly registering channel online check for '%s'", channelName)
	if !isChannelRegistered(channelName) {
		log.Printf("Channel %s not already registered, registering")
		channel := checkChannel{}
		channel.Name = channelName
		channel.FirstCheck = true
		channels = append(channels, &channel)
	} else {
		log.Printf("Channel %s already registered, not registering")
	}
}

func isChannelRegistered(channelName string) bool {
	for _, channel := range channels {
		if channel.Name == channelName {
			return true
		}
	}

	return false
}

func mainLoop() {
	for {
		for _, channel := range channels {
			twitchStreamChannelInstance := checkTwitchStream(channel.Name)

			if twitchStreamChannelInstance.Stream != nil && channel.LastStatus != true {

				channelOnlineInstance := &ChannelOnline{}
				channelOnlineInstance.Channel = twitchStreamChannelInstance
				channelOnlineInstance.Name = channel.Name

				log.Printf("Channel '%s' has gone from offline to online", channel.Name)
				if !channel.FirstCheck {
					log.Printf("This is not the first check of channel '%s' - notifying", channel.Name)
					handle(channelOnlineEventType, channelOnlineInstance)
				}

				channel.LastStatus = true
			} else if twitchStreamChannelInstance.Stream == nil && channel.LastStatus != false {

				channelOfflineInstance := &ChannelOffline{}
				channelOfflineInstance.Channel = twitchStreamChannelInstance
				channelOfflineInstance.Name = channel.Name

				log.Printf("Channel '%s' has gone from online to offline", channel.Name)
				if !channel.FirstCheck {
					log.Printf("This is not the first check of channel '%s' - notifying", channel.Name)
					handle(channelOfflineEventType, channelOfflineInstance)
				}
				channel.LastStatus = false
			}

			if channel.FirstCheck {
				log.Printf("Completed first check for '%s'", channel.Name)
				channel.FirstCheck = false
			}

			time.Sleep(1000 * time.Millisecond)
		}

		time.Sleep(30000 * time.Millisecond)
	}
}

func Run() {
	go mainLoop()
}
