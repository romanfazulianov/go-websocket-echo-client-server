package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

var ws *websocket.Conn

//UserInputHandler is a func which reads the keybord input and tryes to send it
//via websocket connection
func UserInputHandler() {
	reader := bufio.NewReader(os.Stdin)
	for {
		if input, err := reader.ReadString('\n'); err != nil {
			log.Print(err)
			break
		} else {
			if ws != nil {
				if _, err := ws.Write([]byte(input)); err != nil {
					log.Print(err)
					break
				}
			} else {
				log.Print("Connection closed. Try later!")
			}
		}
	}
}

//ServerAnswerHandler reads websocket connection output and log it into terminal
func ServerAnswerHandler() {
	for {
		msg := make([]byte, 1024)
		if _, err := ws.Read(msg); err != nil {
			log.Print("Connection closed...")
			break
		} else {
			log.Printf("Received: %s", msg)
		}
	}
}

func main() {
	var err error
	origin := "http://localhost/"
	url := "ws://localhost:12345/echo"
	go UserInputHandler()
	for {
		log.Print("Trying to connect...")
		if ws, err = websocket.Dial(url, "", origin); err != nil {
			log.Print(err)
			time.Sleep(time.Second * 5)
		} else {
			log.Print("Connection established...")
			ServerAnswerHandler()
		}
	}
}
