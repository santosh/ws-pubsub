package pubsub

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
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

func (p *PubSub) GetSubscriptions(topic string, client *Client) []Subscription {
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
	clientSubs := p.GetSubscriptions(topic, client)
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

func (p *PubSub) Unsubscribe(client *Client, topic string) *PubSub {
	for idx, sub := range p.Subscriptions {
		if sub.Client.Id == client.Id && sub.Topic == topic {
			// found this subscription from client and we need to remove if
			p.Subscriptions = append(p.Subscriptions[:idx], p.Subscriptions[idx+1:]...)
		}
	}

	return p
}

func (p *PubSub) Publish(topic string, message []byte, excludeClient *Client) *PubSub {
	subscriptions := p.GetSubscriptions(topic, nil)
	for _, sub := range subscriptions {
		sub.Client.Send(message)
	}
	log.Printf("published '%s' to %s", message, topic)
	return p
}

func (c *Client) Send(message []byte) error {
	return c.Connection.WriteMessage(websocket.TextMessage, message)
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
		p.Publish(m.Topic, m.Message, &c)
	case SUBSCRIBE:
		p.Subscribe(&c, m.Topic)
	case UNSUBSCRIBE:
		p.Unsubscribe(&c, m.Topic)
	default:
		log.Println("unknown action type")
	}

	return p
}
