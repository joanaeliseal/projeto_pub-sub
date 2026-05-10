package main

import (
	"log"
	"sync"

	"pubsub/shared"
)

const BufferSize = 100

type Topic struct {
	name        string
	mu          sync.RWMutex
	subscribers map[*Client]bool
	messages    chan *shared.Message
	done        chan struct{}
}

func NewTopic(name string) *Topic {
	t := &Topic{
		name:        name,
		subscribers: make(map[*Client]bool),
		messages:    make(chan *shared.Message, BufferSize),
		done:        make(chan struct{}),
	}
	go t.dispatchMessages()
	return t
}

func (t *Topic) AddSubscriber(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subscribers[c] = true
	log.Printf("[INFO] Cliente %s inscrito no tópico %s", c.RemoteAddr(), t.name)
}

func (t *Topic) RemoveSubscriber(c *Client) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.subscribers, c)
	log.Printf("[INFO] Cliente %s removido do tópico %s", c.RemoteAddr(), t.name)
}

func (t *Topic) HasSubscribers() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.subscribers) > 0
}

func (t *Topic) Publish(msg *shared.Message) bool {
	select {
	case t.messages <- msg:
		return true
	default:
		log.Printf("[WARN] Buffer cheio no tópico %s", t.name)
		return false
	}
}

func (t *Topic) dispatchMessages() {
	for {
		select {
		case msg := <-t.messages:
			t.mu.RLock()
			for client := range t.subscribers {
				delivery := shared.NewDelivery(t.name, msg.Payload)
				if err := client.Send(delivery); err != nil {
					log.Printf("[ERROR] Erro ao enviar para %s: %v", client.RemoteAddr(), err)
				}
			}
			t.mu.RUnlock()
		case <-t.done:
			return
		}
	}
}

func (t *Topic) Close() {
	close(t.done)
}
