package main

import (
	"time"

	"github.com/0xAX/notificator"
	"github.com/getlantern/systray"
	"github.com/lobre/goodhosts"
)

type session struct {
	config *Config
	notify *notificator.Notificator
}

func (session *session) start() chan bool {
	complete := make(chan bool)

	// TODO: add quit channel

	go func(chan bool) {
		for repeat := 1; repeat <= session.config.Pomodoro.Repeats; repeat++ {
			lockTimer := time.NewTimer(time.Duration(session.config.Pomodoro.Work) * time.Minute)
			// TODO: use icons instead of text
			systray.SetTitle("I'm working hard")
			session.notify.Push("Work hard", "", "", notificator.UR_NORMAL)
			session.work()

			<-lockTimer.C

			restTimer := time.NewTimer(time.Duration(session.config.Pomodoro.Rest) * time.Minute)
			systray.SetTitle("I'm resting")
			session.notify.Push("Take a rest", "", "", notificator.UR_NORMAL)
			session.rest()

			<-restTimer.C
		}

		complete <- true

	}(complete)

	return complete
}

func (session *session) rest() {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		panic(err)
	}

	for _, bad := range session.config.Domains.Bad {
		if hosts.Has("127.0.0.1", bad) {
			hosts.Remove("127.0.0.1", bad)
			hosts.Remove("127.0.0.1", "www."+bad)
		}
	}

	if err := hosts.Flush(); err != nil {
		panic(err)
	}
}

func (session *session) work() {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		panic(err)
	}

	for _, bad := range session.config.Domains.Bad {
		if !hosts.Has("127.0.0.1", bad) {
			hosts.Add("127.0.0.1", bad)
			hosts.Add("127.0.0.1", "www."+bad)
		}
	}

	if err := hosts.Flush(); err != nil {
		panic(err)
	}
}
