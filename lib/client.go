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

type topicConn struct {
	nc     net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mu     sync.Mutex
}

func newTopicConn(addr string) (*topicConn, error) {
	nc, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &topicConn{
		nc:     nc,
		reader: bufio.NewReader(nc),
		writer: bufio.NewWriter(nc),
	}, nil
}

func (tc *topicConn) send(msg *shared.Message) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	if err := shared.EncodeMessage(tc.writer, msg); err != nil {
		return err
	}
	return tc.writer.Flush()
}

func (tc *topicConn) close() {
	tc.nc.Close()
}

type Client struct {
	addr       string
	conns      map[string]*topicConn
	connsMu    sync.Mutex
	handlers   map[string]MessageHandler
	handlersMu sync.RWMutex
	done       chan struct{}
	closeOnce  sync.Once
}

func NewClient(addr string) *Client {
	return &Client{
		addr:     addr,
		conns:    make(map[string]*topicConn),
		handlers: make(map[string]MessageHandler),
		done:     make(chan struct{}),
	}
}

func (c *Client) Done() <-chan struct{} {
	return c.done
}

func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
		c.connsMu.Lock()
		defer c.connsMu.Unlock()
		for _, tc := range c.conns {
			tc.close()
		}
	})
	return nil
}


func (c *Client) connForTopic(topic string) (*topicConn, error) {
	c.connsMu.Lock()
	defer c.connsMu.Unlock()

	if tc, ok := c.conns[topic]; ok {
		return tc, nil
	}
	tc, err := newTopicConn(c.addr)
	if err != nil {
		return nil, err
	}
	c.conns[topic] = tc
	go c.readLoop(tc)
	return tc, nil
}

func (c *Client) Subscribe(topic string, handler MessageHandler) error {
	c.handlersMu.Lock()
	c.handlers[topic] = handler
	c.handlersMu.Unlock()

	tc, err := c.connForTopic(topic)
	if err != nil {
		return err
	}
	return tc.send(shared.NewSubscribe(topic))
}

func (c *Client) Unsubscribe(topic string) error {
	c.handlersMu.Lock()
	delete(c.handlers, topic)
	c.handlersMu.Unlock()

	c.connsMu.Lock()
	tc, ok := c.conns[topic]
	delete(c.conns, topic)
	c.connsMu.Unlock()

	if ok {
		_ = tc.send(shared.NewUnsubscribe(topic))
		tc.close()
	}
	return nil
}

func (c *Client) Publish(topic string, payload any) error {
	msg, err := shared.NewPublish(topic, payload)
	if err != nil {
		return err
	}
	tc, err := c.connForTopic(topic)
	if err != nil {
		return err
	}
	return tc.send(msg)
}

func (c *Client) readLoop(tc *topicConn) {
	for {
		line, err := tc.reader.ReadBytes('\n')
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
