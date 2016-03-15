package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
)

type HistoryMsg struct {
	Id   string
	Msg  string
	Type string
}

var (
	dirPath    string
	clients    = make(map[*websocket.Conn]bool)
	typeDetect string
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (group *HistoryMsg) sendAll() {
	// Enumerating clients from "clients" map and sending msg to them
	for client := range clients {
		err := client.WriteJSON(group)
		if err != nil {
			delete(clients, client)
			return
		}
	}
}

func ChatServer(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Connecting to MongoDB, collection History
	session, err := mgo.Dial("mongodb://darik954:darik954@ds045795.mongolab.com:45795/catalog")
	if err != nil {
		check(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("catalog").C("History")

	var historyMsg []HistoryMsg
	c.Find(bson.M{}).Sort("_id").All(&historyMsg)
	for _, connect := range historyMsg {
		err = conn.WriteJSON(connect)
		if err != nil {
			check(err)
		}
	}

	// Adding clients to the map
	clientId := conn.RemoteAddr().String()
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()
	clients[conn] = true

	// Loop for receiving msg
	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			check(err)
		}
		fmt.Println(messageType)

		switch messageType {

		case 1:
			clientMsg := string(r)

			group := HistoryMsg{
				Id:   clientId,
				Msg:  clientMsg,
				Type: "Text",
			}

			group.sendAll()

			err = c.Insert(group)
			if err != nil {
				check(err)
			}

		case 2:
			imgPath := fmt.Sprintf("img/cache/%d.png", time.Now().UnixNano())

			fo, err := os.Create(dirPath + imgPath)
			if err != nil {
				check(err)
			}

			if _, err := fo.Write(r); err != nil {
				check(err)
			}

			fo.Close()

			group := HistoryMsg{
				Id:   clientId,
				Msg:  imgPath,
				Type: "Image",
			}

			group.sendAll()

			err = c.Insert(group)
			if err != nil {
				check(err)
			}

			fmt.Println("Ok")
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: chatExample <dir>")
	}

	dirPath = os.Args[1]

	fmt.Println("Starting... ")

	http.Handle("/", http.FileServer(http.Dir(dirPath)))
	http.HandleFunc("/ws", ChatServer)

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}
