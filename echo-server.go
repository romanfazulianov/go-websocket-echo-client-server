package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var ws *websocket.Conn

// EchoServer echoes the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	for {
		msg := make([]byte, 1024)
		if _, err := ws.Read(msg); err != nil {
			log.Print(err)
		}
		log.Printf("Received: %s", msg)
		if _, err := ws.Write(msg); err != nil {
			log.Print(err)
		}
	}
}

// This example demonstrates a trivial echo server.
func main() {
	log.Print("Starting server...")
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
