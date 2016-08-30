package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var cons []*websocket.Conn

// EchoServer echoes the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	log.Print(ws)
	cons = append(cons, ws)
	handleEchoes(ws)
}

func handleEchoes(ws *websocket.Conn) {
	for {
		msg := make([]byte, 1024)
		if _, err := ws.Read(msg); err != nil {
			log.Print(err)
		}
		log.Printf("Received: %s", msg)
		for _, wsconn := range cons {
			if _, err := wsconn.Write(msg); err != nil {
				log.Print(err)
				wsconn.Close()
			}
		}
	}
}

// This example demonstrates a trivial echo server.
func main() {
	log.Print("Starting server...")
	http.Handle("/multi", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
