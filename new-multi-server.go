package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var cons []*websocket.Conn

// MultiServer sends the data received on the WebSocket to all connections.
func MultiServer(ws *websocket.Conn) {
	log.Printf("Connected %v", ws)
	cons = append(cons, ws)
	handleMessages(ws)
}

func splice(cons []*websocket.Conn, ws *websocket.Conn) []*websocket.Conn {
	newCons := []*websocket.Conn{}
	for _, con := range cons {
		if con != ws {
			newCons = append(newCons, con)
		} else {
			log.Printf("Disconnected %v", ws)
		}
	}
	return newCons
}

func handleMessages(ws *websocket.Conn) {
	for {
		msg := make([]byte, 1024)
		if _, err := ws.Read(msg); err != nil {
			ws.Close()
			cons = splice(cons, ws)
			break
		}
		log.Printf("Received: %s", msg)
		for _, wsconn := range cons {
			if _, err := wsconn.Write(msg); err != nil {
				wsconn.Close()
				cons = splice(cons, wsconn)
				break
			}
		}
	}
}

func main() {
	log.Print("Starting server...")
	http.Handle("/", websocket.Handler(MultiServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
