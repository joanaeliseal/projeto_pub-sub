package lib

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"sync"

	"pubsub/shared"
)

type MessageHandler func(topic string, payload json.RawMessage)

type Client struct {
	addr       string
	conn       net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	mu         sync.Mutex
	handlers   map[string]MessageHandler
	handlersMu sync.RWMutex
	done       chan struct{}
	closeOnce  sync.Once
}

func NewClient(addr string) *Client {
	return &Client{
		addr:     addr,
		handlers: make(map[string]MessageHandler),
		done:     make(chan struct{}),
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	go c.readLoop()
	return nil
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
	})
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) Done() <-chan struct{} {
	return c.done
}

func (c *Client) Publish(topic string, payload any) error {
	msg, err := shared.NewPublish(topic, payload)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := shared.EncodeMessage(c.writer, msg); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	c.handlersMu.Lock()
	c.handlers[topic] = handler
	c.handlersMu.Unlock()

	msg := shared.NewSubscribe(topic)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := shared.EncodeMessage(c.writer, msg); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) Unsubscribe(topic string) error {
	c.handlersMu.Lock()
	delete(c.handlers, topic)
	c.handlersMu.Unlock()

	msg := shared.NewUnsubscribe(topic)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := shared.EncodeMessage(c.writer, msg); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) readLoop() {
	for {
		line, err := c.reader.ReadBytes('\n')
		if err != nil {
			select {
			case <-c.done:
			default:
				log.Printf("[LIB] Conexão encerrada inesperadamente: %v", err)
			}
			return
		}

		var raw struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		switch raw.Type {
		case shared.TypeMessage:
			var msg shared.Message
			if err := json.Unmarshal(line, &msg); err != nil {
				continue
			}
			c.handlersMu.RLock()
			h := c.handlers[msg.Topic]
			c.handlersMu.RUnlock()
			if h != nil {
				go h(msg.Topic, msg.Payload)
			}
		case shared.TypeError:
			var resp shared.Response
			if err := json.Unmarshal(line, &resp); err != nil {
				continue
			}
			log.Printf("[LIB] Aviso do broker [%s]: %s", resp.Code, resp.Message)
		}
	}
}
