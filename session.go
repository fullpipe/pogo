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

func (session *session) start() chan struct{} {
	complete := make(chan struct{})

	// TODO: add quit channel

	go func(chan struct{}) {
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

		complete <- struct{}{}

	}(complete)

	return complete
}

func (session *session) rest() {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		// TODO: may be it`s better just notify and stop
		panic(err)
	}

	for _, bad := range session.config.Domains.Bad {
		if hosts.Has(session.config.WorkIp, bad) {
			hosts.Remove(session.config.WorkIp, bad)
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
		if !hosts.Has(session.config.WorkIp, bad) {
			hosts.Add(session.config.WorkIp, bad)
		}
	}

	if err := hosts.Flush(); err != nil {
		panic(err)
	}
}
