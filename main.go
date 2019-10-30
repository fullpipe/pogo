package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/0xAX/notificator"
	"github.com/getlantern/systray"
	"github.com/sevlyar/go-daemon"
)

var (
	signal = flag.String("h", "", `Send signal to the daemon:
  quit — graceful shutdown
  stop — fast shutdown
  reload — reloading the configuration file`)
	start = false
)

func main() {
	flag.Parse()
	if len(os.Args[1:]) > 0 {
		fmt.Println(os.Args[1:])
		start = os.Args[1:][0] == "start"
	}
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)
	daemon.AddCommand(daemon.BoolFlag(&start), syscall.SIGHUP, reloadHandler)

	cntxt := &daemon.Context{
		PidFileName: "pogo.pid",
		PidFilePerm: 0644,
		LogFileName: "pogo.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"pogo"},
	}

	fmt.Println(daemon.ActiveFlags())
	fmt.Println(os.Args[1:])
	if len(os.Args[1:]) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	go worker()
	go systray.Run(onReady, onExit)

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")

}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
LOOP:
	for {
		time.Sleep(time.Second) // this is work to be done by worker.
		select {
		case <-stop:
			//systray.Run(onReady, onExit)
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}

func onReady() {
	// TODO: add icon
	systray.SetTitle("I'm going to work hard")
	//systray.SetTooltip("Look at me, I'm a tooltip!")

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
				// TODO: enable new session button on session complete
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
