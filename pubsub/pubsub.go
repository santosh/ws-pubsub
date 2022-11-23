package pubsub

import "github.com/gorilla/websocket"

type PubSub struct {
	Clients []Client
}

type Client struct {
	Id         string
	Connection *websocket.Conn
}

// AddClient adds a client to PubSub instance
func (p *PubSub) AddClient(c Client) {
	p.Clients = append(p.Clients, c)
}
