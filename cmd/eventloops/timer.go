package eventloops

import (
	"log"
	"time"
)

type Handler struct {
	Callback    callback
	Period      float64
	LastTrigger time.Time
}

var (
	handlers []*Handler
)

type callback func()

func RegisterTimer(Callback callback, period float64) {
	handler := Handler{}
	handler.Callback = Callback
	handler.Period = period
	handler.LastTrigger = time.Now()
	handlers = append(handlers, &handler)
	log.Printf("Registerd timer with period %d", period)
}

func mainLoop() {
	for {
		for _, handler := range handlers {
			if time.Since(handler.LastTrigger).Seconds() > handler.Period {
				go handler.Callback()
				handler.LastTrigger = time.Now()
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func Run() {
	go mainLoop()
}
