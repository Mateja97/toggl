package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mateja97/toggl/toggl"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	port := os.Getenv("PORT")
	t := toggl.Toggl{}
	if err := t.Init(port); err != nil {
		log.Fatal("[ERROR] Service failed to init, error: ", err)
	}
	go t.Run()
	done := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
	t.Stop()
}
