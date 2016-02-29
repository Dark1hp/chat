package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "io/ioutil"
	"log"
	"net/http"
	"os"
	// "strconv"
	"time"
)

var (
	dirPath string
	test    byte
	lenFile int
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ChatServer(ws *websocket.Conn) {

	defer ws.Close()

	for {
		payload := make([]byte, 32*8192)
		n, err := ws.Read(payload)
		if err != nil {
			log.Println(err)
		}
		lenFile += n
		// fmt.Println(lenFile)
		test := append(payload[:n])
		fmt.Println(len(test))

		fo, err := os.Create(fmt.Sprintf("./%d.png", time.Now().UnixNano()))
		check(err)
		_, err = fo.Write(test[:lenFile])
		check(err)
		fo.Close()
	}
	log.Print("DONE")

}

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
