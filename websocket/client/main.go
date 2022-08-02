package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:80", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	conn := c.UnderlyingConn()
	go func() { io.Copy(conn, os.Stdin) }()
	io.Copy(os.Stdout, conn)
}