package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

var ws *websocket.Conn

//ClientInfo contain the name and room of client
type ClientInfo struct {
	name string
	room string
}

//Msg contain a message
type Msg struct {
	client *ClientInfo
	sent   time.Time
	text   string
}

var client *ClientInfo

type termScan bufio.Scanner

var scanner = bufio.NewScanner(os.Stdin)

func scanLine(s *bufio.Scanner) (string, error) {
	s.Scan()
	err := s.Err()
	if err != nil {
		log.Print(err)
	}
	return s.Text(), err
}

func fillInfo() {
	fmt.Println("Name yourself:")
	name, err := scanLine(scanner)
	if err != nil {
		return
	}
	fmt.Println("What room do you want 2 connect:")
	room, err := scanLine(scanner)
	if err != nil {
		return
	}
	client = &ClientInfo{name, room}
}

//UserInputHandler is a func which reads the keybord input and tryes to send it
//via websocket connection
func UserInputHandler() {
	for {
		if input, err := scanLine(scanner); err == nil {
			if ws != nil {
				if _, err := ws.Write([]byte(input)); err != nil {
					log.Print(err)
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
	fillInfo()
	config, err := websocket.NewConfig("ws://localhost:8080/", "http://localhost/")
	if err != nil {
		log.Fatal("Unable to create config for client!")
	}

	go UserInputHandler()

	for {
		log.Print("Trying to connect...")
		if ws, err = websocket.DialConfig(config); err != nil {
			log.Print(err)
			time.Sleep(time.Second * 5)
		} else {
			log.Print("Connection established...")
			ServerAnswerHandler()
		}
	}
}
