package schedule

import (
	"os"
	"time"
)

type Event struct {
	time            time.Time
	callbackChannel chan Event
}

type Schedule struct {
	ticker   *time.Ticker
	shutdown chan os.Signal
	events   []Event
}

func Init() Schedule {
	return Schedule{
		shutdown: make(chan os.Signal, 1),
	}
}

func (s *Schedule) Start() {
	s.ticker = s.initTicker()

	s.startTicker()
}

func (s *Schedule) Stop() {
	s.shutdown <- os.Kill
}

func (s *Schedule) NewEventChannel() chan Event {
	c := make(chan Event)

	return c
}

func (s *Schedule) InsertEvent(eventDelay time.Duration, callbackChannel chan Event) {
	s.events = append(s.events, Event{
		time:            time.Now().Add(eventDelay),
		callbackChannel: callbackChannel,
	})
}

func (s *Schedule) initTicker() *time.Ticker {
	return time.NewTicker(time.Second)
}

func (s *Schedule) startTicker() {
	for {
		select {
		case <-s.ticker.C:
			s.checkEvents()
		case <-s.shutdown:
			s.ticker.Stop()
			return
		}
	}
}

func (s *Schedule) checkEvents() {
	var unpoppedEvents []Event

	for _, e := range s.events {
		if time.Now().After(e.time) {
			e.callbackChannel <- e
		} else {
			unpoppedEvents = append(unpoppedEvents, e)
		}
	}

	s.events = unpoppedEvents
}
