package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/santosh/pubsub-youtube/pubsub"
)

var ps = &pubsub.PubSub{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := pubsub.Client{
		Id:         newUUID(),
		Connection: conn,
	}

	ps.AddClient(client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static")
	})

	http.HandleFunc("/ws", WsHandler)
	http.ListenAndServe(":3000", nil)
}
