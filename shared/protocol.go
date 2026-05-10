package shared

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	TypePublish     = "publish"
	TypeSubscribe   = "subscribe"
	TypeUnsubscribe = "unsubscribe"
	TypeMessage     = "message"
	TypeAck         = "ack"
	TypeError       = "error"
)

const (
	ErrInvalidJSON   = "INVALID_JSON"
	ErrUnknownType   = "UNKNOWN_TYPE"
	ErrTopicRequired = "TOPIC_REQUIRED"
	ErrNoSubscribers = "NO_SUBSCRIBERS"
	ErrBufferFull    = "BUFFER_FULL"
)

type Message struct {
	Type      string          `json:"type"`
	Topic     string          `json:"topic,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

type Response struct {
	Type    string `json:"type"`
	Success bool   `json:"success,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewPublish(topic string, payload any) (*Message, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Message{
		Type:    TypePublish,
		Topic:   topic,
		Payload: data,
	}, nil
}

func NewSubscribe(topic string) *Message {
	return &Message{
		Type:  TypeSubscribe,
		Topic: topic,
	}
}

func NewUnsubscribe(topic string) *Message {
	return &Message{
		Type:  TypeUnsubscribe,
		Topic: topic,
	}
}

func NewDelivery(topic string, payload json.RawMessage) *Message {
	return &Message{
		Type:      TypeMessage,
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func NewAck(topic string) *Response {
	return &Response{
		Type:    TypeAck,
		Success: true,
		Topic:   topic,
	}
}

func NewError(code, message string) *Response {
	return &Response{
		Type:    TypeError,
		Success: false,
		Code:    code,
		Message: message,
	}
}

func EncodeMessage(w io.Writer, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

func DecodeMessage(r *bufio.Reader) (*Message, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := json.Unmarshal(line, &msg); err != nil {
		return nil, fmt.Errorf("json inválido: %w", err)
	}
	return &msg, nil
}

func EncodeResponse(w io.Writer, resp *Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

func DecodeResponse(r *bufio.Reader) (*Response, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var resp Response
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
