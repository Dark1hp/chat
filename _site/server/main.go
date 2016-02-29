package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	// "io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type (
	Msg struct {
		clientKey string
		text      string
	}
	NewClientEvent struct {
		clientKey string
		msgChan   chan *Msg
	}
)

var (
	dirPath           string
	clientRequests    = make(chan *NewClientEvent, 100)
	clientDisconnects = make(chan string, 100)
	messages          = make(chan *Msg, 100)
)

/*func IndexPage(w http.ResponseWriter, req *http.Request) {

	fp, err := os.Open(dirPath + "/index.html")
	if err != nil {
		log.Println("Could not open file", err.Error())
		w.Write([]byte("500 internal server error"))
		return
	}

	defer fp.Close()

	_, err = io.Copy(w, fp)
	if err != nil {
		log.Println("Could not send file contents", err.Error())
		w.Write([]byte("500 internal server error"))
		return
	}

}*/

func ChatServer(ws *websocket.Conn) {
	lenBuf := make([]byte, 2048)
	// ws.SetDeadline(TIMEOUT)

	msgChan := make(chan *Msg, 100)
	clientKey := ws.RemoteAddr().String()
	clientRequests <- &NewClientEvent{clientKey, msgChan}
	defer func() { clientDisconnects <- clientKey }()

	go func() {
		for msg := range msgChan {
			ws.Write([]byte(msg.text))
		}
	}()

	for {
		_, err := ws.Read(lenBuf)
		if err != nil {
			log.Println("Error: ", err.Error())
			return
		}

		length, _ := strconv.Atoi(strings.TrimSpace(string(lenBuf)))
		if length > 65536 {
			log.Println("Error: too big length: ", length)
			return
		}

		if length <= 0 {
			log.Println("Empty length: ", length)
			return
		}

		buf := make([]byte, length)
		_, err = ws.Read(buf)

		if err != nil {
			log.Println("Could not read ", length, " bytes", err.Error())
			return
		}

		messages <- &Msg{clientKey, string(buf)}
	}

}

func router() {
	clients := make(map[string]chan *Msg)

	for {
		select {
		case req := <-clientRequests:
			clients[req.clientKey] = req.msgChan
			log.Println("Websocket connected: " + req.clientKey)

		case clientKey := <-clientDisconnects:
			close(clients[clientKey])
			delete(clients, clientKey)
			log.Println("Websocket disconnected: " + clientKey)

		case msg := <-messages:
			for _, msgChan := range clients {
				if len(msgChan) < cap(msgChan) {
					msgChan <- msg
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: chatExample <dir>")
	}

	dirPath = os.Args[1]

	fmt.Println("Starting... ")

	go router()

	http.Handle("/", http.FileServer(http.Dir(dirPath)))
	http.Handle("/ws", websocket.Handler(ChatServer))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}
