package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mateja97/toggl/toggl"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	t := toggl.Toggl{}
	t.Init(":8080")
	go t.Run()
	done := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
	t.Stop()
}
