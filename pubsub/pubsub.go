package pubsub

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

const (
	PUBLISH   = "publish"
	SUBSCRIBE = "subscribe"
)

type PubSub struct {
	Clients       []Client
	Subscriptions []Subscription
}

type Client struct {
	Id         string
	Connection *websocket.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Topic  string
	Client *Client
}

// AddClient adds a client to PubSub instance
func (p *PubSub) AddClient(c Client) {
	p.Clients = append(p.Clients, c)
}

func (p *PubSub) GetSubscription(topic string, client *Client) []Subscription {
	var subscriptionList []Subscription

	for _, subscription := range p.Subscriptions {
		if client != nil {
			if subscription.Client.Id == client.Id {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

func (p *PubSub) Subscribe(client *Client, topic string) *PubSub {
	clientSubs := p.GetSubscription(topic, client)
	if len(clientSubs) > 0 {
		return p
	}

	newSubscription := Subscription{
		Topic:  topic,
		Client: client,
	}

	p.Subscriptions = append(p.Subscriptions, newSubscription)
	log.Printf("%s subscribed to %s", client.Id, topic)
	return p
}

// HandleReceiveMessage fetches message and acts based on Action of the Message
func (p *PubSub) HandleReceiveMessage(c Client, messageType int, payload []byte) *PubSub {
	m := Message{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Printf("wrong message payload: %v", err)
	}

	switch m.Action {
	case PUBLISH:
		log.Println("need to handle publish action")
	case SUBSCRIBE:
		p.Subscribe(&c, m.Topic)
	default:
		log.Println("unknown action type")
	}

	return p
}
