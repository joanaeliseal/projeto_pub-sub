package main

import (
	"io"
	"log"
	"net"

	"pubsub/shared"
)

const DefaultPort = ":9000"

var broker *Broker

func main() {
	broker = NewBroker()

	listener, err := net.Listen("tcp", DefaultPort)
	if err != nil {
		log.Fatalf("[FATAL] Erro ao iniciar servidor: %v", err)
	}
	defer listener.Close()

	log.Printf("[INFO] Broker iniciado na porta %s", DefaultPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] Erro ao aceitar conexão: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	client := NewClient(conn)
	defer cleanup(client)

	log.Printf("[INFO] Cliente conectado: %s", client.RemoteAddr())

	for {
		msg, err := client.Read()
		if err != nil {
			if err != io.EOF {
				log.Printf("[ERROR] Erro ao ler mensagem: %v", err)
			}
			return
		}

		handleMessage(client, msg)
	}
}

func handleMessage(client *Client, msg *shared.Message) {
	switch msg.Type {
	case shared.TypeSubscribe:
		handleSubscribe(client, msg)
	case shared.TypeUnsubscribe:
		handleUnsubscribe(client, msg)
	case shared.TypePublish:
		handlePublish(client, msg)
	default:
		resp := shared.NewError(shared.ErrUnknownType, "tipo de mensagem desconhecido")
		client.SendResponse(resp)
	}
}

func handleSubscribe(client *Client, msg *shared.Message) {
	if msg.Topic == "" {
		resp := shared.NewError(shared.ErrTopicRequired, "tópico é obrigatório")
		client.SendResponse(resp)
		return
	}

	topic := broker.GetOrCreateTopic(msg.Topic)
	topic.AddSubscriber(client)
	client.AddTopic(msg.Topic)

	resp := shared.NewAck(msg.Topic)
	client.SendResponse(resp)
}

func handleUnsubscribe(client *Client, msg *shared.Message) {
	if msg.Topic == "" {
		resp := shared.NewError(shared.ErrTopicRequired, "tópico é obrigatório")
		client.SendResponse(resp)
		return
	}

	topic := broker.GetTopic(msg.Topic)
	if topic != nil {
		topic.RemoveSubscriber(client)
		client.RemoveTopic(msg.Topic)
		broker.RemoveTopicIfEmpty(msg.Topic)
	}

	resp := shared.NewAck(msg.Topic)
	client.SendResponse(resp)
}

func handlePublish(client *Client, msg *shared.Message) {
	if msg.Topic == "" {
		resp := shared.NewError(shared.ErrTopicRequired, "tópico é obrigatório")
		client.SendResponse(resp)
		return
	}

	topic := broker.GetTopic(msg.Topic)
	if topic == nil || !topic.HasSubscribers() {
		log.Printf("[INFO] Mensagem descartada - sem subscribers no tópico %s", msg.Topic)
		resp := shared.NewError(shared.ErrNoSubscribers, "mensagem descartada: sem subscribers")
		client.SendResponse(resp)
		return
	}

	topic.Publish(msg)
	log.Printf("[INFO] Mensagem publicada no tópico %s", msg.Topic)

	resp := shared.NewAck(msg.Topic)
	client.SendResponse(resp)
}

func cleanup(client *Client) {
	for _, topicName := range client.GetTopics() {
		topic := broker.GetTopic(topicName)
		if topic != nil {
			topic.RemoveSubscriber(client)
			broker.RemoveTopicIfEmpty(topicName)
		}
	}
	client.Close()
	log.Printf("[INFO] Cliente desconectado: %s", client.RemoteAddr())
}
