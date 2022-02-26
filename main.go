package main

import (
	"os"
	"os/signal"

	"github.com/Akshit8/simple-redis-clone/data"
	"github.com/Akshit8/simple-redis-clone/server"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ds := data.NewStore("snapshot.json")
	srv := server.NewServer(ds)

	select {
	case <-c:
		srv.Stop()
	}
}
