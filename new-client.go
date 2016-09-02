package main

import (
	"bufio"
	"encoding/json"
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
	Name string
	Room string
	Sent int
	Text string
}

var client ClientInfo

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
	client = ClientInfo{name: name, room: room}
}

//UserInputHandler is a func which reads the keybord input and tryes to send it
//via websocket connection
func UserInputHandler() {
	for {
		if input, err := scanLine(scanner); err == nil {
			if ws != nil {
				fmt.Println(client)
				timestamp := time.Now().Nanosecond()
				msg := Msg{Name: client.name, Room: client.room, Sent: timestamp, Text: input}
				log.Printf("sending: %+v", msg)
				data, err := json.Marshal(msg)
				if err != nil {
					log.Print(err)
				} else {
					if _, err := ws.Write(data); err != nil {
						log.Print(err)
					}
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
		data := make([]byte, 1024)
		n, err := ws.Read(data)
		if err != nil {
			log.Print("Connection closed...")
			break
		}
		var msg Msg
		err = json.Unmarshal(data[:n], &msg)
		if err != nil {
			fmt.Println("error:", err)
		}
		log.Printf("Received: %+v", msg)
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
