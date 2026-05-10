package main

import (
	"log"
	"sync"
)

// Broker é o componente central do sistema pub/sub
type Broker struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

// NewBroker cria uma nova instância do broker
func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

// GetOrCreateTopic obtém um tópico existente ou cria um novo
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

// GetTopic obtém um tópico existente
func (b *Broker) GetTopic(name string) *Topic {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.topics[name]
}

// TopicExists verifica se um tópico existe
func (b *Broker) TopicExists(name string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, exists := b.topics[name]
	return exists
}

// RemoveTopicIfEmpty remove um tópico se não houver subscribers
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

// ListTopics retorna a lista de tópicos ativos
func (b *Broker) ListTopics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	names := make([]string, 0, len(b.topics))
	for name := range b.topics {
		names = append(names, name)
	}
	return names
}
