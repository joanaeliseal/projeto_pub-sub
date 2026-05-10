package shared

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Tipos de mensagem do protocolo
const (
	TypePublish     = "publish"
	TypeSubscribe   = "subscribe"
	TypeUnsubscribe = "unsubscribe"
	TypeMessage     = "message"
	TypeAck         = "ack"
	TypeError       = "error"
)

// Códigos de erro
const (
	ErrInvalidJSON   = "INVALID_JSON"
	ErrUnknownType   = "UNKNOWN_TYPE"
	ErrTopicRequired = "TOPIC_REQUIRED"
	ErrNoSubscribers = "NO_SUBSCRIBERS"
)

// Message representa uma mensagem do protocolo
type Message struct {
	Type      string          `json:"type"`
	Topic     string          `json:"topic,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

// Response representa uma resposta do broker
type Response struct {
	Type    string `json:"type"`
	Success bool   `json:"success,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewPublish cria uma mensagem de publish
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

// NewSubscribe cria uma mensagem de subscribe
func NewSubscribe(topic string) *Message {
	return &Message{
		Type:  TypeSubscribe,
		Topic: topic,
	}
}

// NewUnsubscribe cria uma mensagem de unsubscribe
func NewUnsubscribe(topic string) *Message {
	return &Message{
		Type:  TypeUnsubscribe,
		Topic: topic,
	}
}

// NewDelivery cria uma mensagem para entregar ao subscriber
func NewDelivery(topic string, payload json.RawMessage) *Message {
	return &Message{
		Type:      TypeMessage,
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewAck cria uma resposta de sucesso
func NewAck(topic string) *Response {
	return &Response{
		Type:    TypeAck,
		Success: true,
		Topic:   topic,
	}
}

// NewError cria uma resposta de erro
func NewError(code, message string) *Response {
	return &Response{
		Type:    TypeError,
		Success: false,
		Code:    code,
		Message: message,
	}
}

// EncodeMessage serializa uma mensagem para JSON com newline
func EncodeMessage(w io.Writer, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

// DecodeMessage lê e deserializa uma mensagem JSON
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

// EncodeResponse serializa uma resposta para JSON com newline
func EncodeResponse(w io.Writer, resp *Response) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}

// DecodeResponse lê e deserializa uma resposta JSON
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
