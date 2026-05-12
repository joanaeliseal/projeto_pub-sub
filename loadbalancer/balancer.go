package main

import (
	"bufio"
	"encoding/json"
	"hash/fnv"
	"io"
	"log"
	"net"
	"sync"
)

type Balancer struct {
	mu      sync.RWMutex
	brokers []string
}

func NewBalancer(brokers []string) *Balancer {
	return &Balancer{brokers: brokers}
}

func (b *Balancer) SelectBroker(topic string) string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	h := fnv.New32a()
	h.Write([]byte(topic))
	idx := int(h.Sum32()) % len(b.brokers)
	return b.brokers[idx]
}

func (b *Balancer) ProxyConnection(clientConn net.Conn) {
	defer clientConn.Close()

	reader := bufio.NewReader(clientConn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Printf("[LB] Erro ao ler primeira mensagem de %s: %v", clientConn.RemoteAddr(), err)
		return
	}

	var raw struct {
		Topic string `json:"topic"`
	}
	json.Unmarshal(line, &raw)

	brokerAddr := b.SelectBroker(raw.Topic)
	log.Printf("[LB] %s → tópico '%s' → %s", clientConn.RemoteAddr(), raw.Topic, brokerAddr)

	brokerConn, err := net.Dial("tcp", brokerAddr)
	if err != nil {
		log.Printf("[LB] Erro ao conectar ao broker %s: %v", brokerAddr, err)
		return
	}
	defer brokerConn.Close()

	if _, err := brokerConn.Write(line); err != nil {
		log.Printf("[LB] Erro ao reenviar mensagem ao broker: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(brokerConn, reader)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientConn, brokerConn)
	}()

	wg.Wait()
}
