package main

import (
	"bufio"
	"net"
	"sync"

	"pubsub/shared"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	mu     sync.Mutex
	topics map[string]bool
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		topics: make(map[string]bool),
	}
}

func (c *Client) Send(msg *shared.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := shared.EncodeMessage(c.writer, msg); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) SendResponse(resp *shared.Response) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := shared.EncodeResponse(c.writer, resp); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Client) Read() (*shared.Message, error) {
	return shared.DecodeMessage(c.reader)
}

func (c *Client) AddTopic(topic string) {
	c.topics[topic] = true
}

func (c *Client) RemoveTopic(topic string) {
	delete(c.topics, topic)
}

func (c *Client) GetTopics() []string {
	topics := make([]string, 0, len(c.topics))
	for t := range c.topics {
		topics = append(topics, t)
	}
	return topics
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
