package main

import (
	"log"
	"sync"
)

type Broker struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

func (b *Broker) GetOrCreateTopic(name string) *Topic {
	b.mu.Lock()
	defer b.mu.Unlock()
	if topic, exists := b.topics[name]; exists {
		return topic
	}
	topic := NewTopic(name)
	b.topics[name] = topic
	log.Printf("[INFO] Tópico criado: %s", name)
	return topic
}

func (b *Broker) GetTopic(name string) *Topic {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.topics[name]
}

func (b *Broker) TopicExists(name string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, exists := b.topics[name]
	return exists
}

func (b *Broker) RemoveTopicIfEmpty(name string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if topic, exists := b.topics[name]; exists {
		if !topic.HasSubscribers() {
			topic.Close()
			delete(b.topics, name)
			log.Printf("[INFO] Tópico removido: %s", name)
		}
	}
}

func (b *Broker) ListTopics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	names := make([]string, 0, len(b.topics))
	for name := range b.topics {
		names = append(names, name)
	}
	return names
}
