package main

import (
	"github.com/0xAX/notificator"
	"github.com/getlantern/systray"
)

func main() {
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
	newSession := systray.AddMenuItem("New session", "Start new pomodoro session")
	quit := systray.AddMenuItem("Quit", "Quit pomodoro session and app")
	var completeCh chan bool

	go func() {
		for {
			select {
			case <-newSession.ClickedCh:
				newSession.Disable()
				session := &session{
					config: config,
					notify: notify,
				}
				completeCh = session.start()

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

}
