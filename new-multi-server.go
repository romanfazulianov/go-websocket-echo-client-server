package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var cons []*websocket.Conn

//ClientInfo contain the name and room of client
type ClientInfo struct {
	name string
	room string
}

//Msg contain a message
type Msg struct {
	Name string
	Room string
	Sent int
	Text string
}

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
		data := make([]byte, 1024)
		n, err := ws.Read(data)
		if err != nil {
			ws.Close()
			cons = splice(cons, ws)
			break
		}
		var msg Msg
		err = json.Unmarshal(data[:n], &msg)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Printf("%+v", msg)
		for _, wsconn := range cons {
			if _, err := wsconn.Write(data[:n]); err != nil {
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
