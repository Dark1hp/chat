package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
)

type Connect struct {
	Clientid  string
	Clientmsg string
}

var (
	dirPath   string
	clients   = make(map[*websocket.Conn]bool)
	varDetect string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sendAll(msg string) {

	// Enumerating clients form "clients" map and sending msg to them
	for client := range clients {
		if err := websocket.Message.Send(client, msg); err != nil {
			delete(clients, client)
			return
		}
	}
}

func saveImg(imgType string, sizeData []byte) {
	fo, err := os.Create(fmt.Sprintf("./_site/assets/img/cache/%d."+imgType, time.Now().UnixNano()))
	check(err)
	_, err = fo.Write(sizeData)
	check(err)
	fo.Close()
	varDetect = "msg"
}

func ChatServer(ws *websocket.Conn) {

	// Connecting to MongoDB, collection History
	session, err := mgo.Dial("mongodb://darik954:darik954@ds045795.mongolab.com:45795/catalog")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("catalog").C("History")

	var historyMsg []Connect
	c.Find(bson.M{}).Sort("_id").All(&historyMsg)
	for _, connect := range historyMsg {
		if err := websocket.Message.Send(ws, connect.Clientmsg); err != nil {
			log.Fatal(err)
		}
	}

	// Adding clients to the map
	clientId := ws.RemoteAddr().String()
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()
	clients[ws] = true

	// Loop for receiving msg
	for {
		var payload []byte
		err := websocket.Message.Receive(ws, &payload)
		if err != nil {
			delete(clients, ws)
			return
		}
		sizeData := append(payload)
		fmt.Println(len(sizeData))

		filetype := http.DetectContentType(payload[:512])
		fmt.Println(filetype)

		switch filetype {
		case "image/jpg":
			saveImg("jpg", sizeData)
			varDetect = "jpg"

		case "image/gif":
			saveImg("gif", sizeData)
			varDetect = "gif"

		case "image/png":
			saveImg("png", sizeData)
			varDetect = "png"

		default:
			switch varDetect {
			case "jpg":
				saveImg("jpg", sizeData)
			case "gif":
				saveImg("gif", sizeData)
			case "png":
				saveImg("png", sizeData)
			default:
				sendAll(string(sizeData))
				err = c.Insert(&Connect{clientId, string(sizeData)})
				if err != nil {
					log.Fatal(err)
				}
				varDetect = "msg"
			}
		}

	}
}

/*	for {
	var msg string
	// If can not read msg - delete client from map
	if err := websocket.Message.Receive(ws, &msg); err != nil {
		delete(clients, ws)
		return
	}
	sendAll(msg)

}*/

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: chatExample <dir>")
	}

	dirPath = os.Args[1]

	fmt.Println("Starting... ")

	http.Handle("/", http.FileServer(http.Dir(dirPath)))
	http.Handle("/ws", websocket.Handler(ChatServer))

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Fatal("ListenAndServe ", err)
	}
}
