package main

import (
	"flag"
	"fmt"

	"github.com/0xAX/notificator"
	"github.com/getlantern/systray"
)

var (
	currentSession *session
	stepCommand    string
)

func main() {
	flag.StringVar(&stepCommand, "c", "", "command to execute between steps")
	flag.Parse()

	systray.Run(onReady, onExit)
}

func onReady() {
	// TODO: add icon
	systray.SetTitle("I'm going to work hard")

	var notify *notificator.Notificator
	// TODO: add icon
	notify = notificator.New(notificator.Options{
		AppName: "Pogo",
	})

	config := getConfig()
	newSession := systray.AddMenuItem("Start session", "Start new pomodoro session")
	quit := systray.AddMenuItem("Quit", "Quit pomodoro session and app")
	var completeCh chan struct{}

	go func() {
		for {
			select {
			case <-newSession.ClickedCh:
				newSession.Disable()
				currentSession = &session{
					config:      config,
					notify:      notify,
					stepCommand: stepCommand,
				}
				fmt.Printf(stepCommand)
				completeCh = currentSession.start()

			case <-completeCh:
				newSession.Enable()

			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	if currentSession != nil {
		currentSession.rest()
	}
}
