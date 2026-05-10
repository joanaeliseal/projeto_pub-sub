// Package lib fornece a biblioteca cliente para o middleware Pub/Sub.
//
// ARQUITETURA:
// Esta biblioteca abstrai toda a comunicação TCP com o broker,
// permitindo que aplicações cliente utilizem uma API simples:
//
//   client := lib.NewClient("localhost:9000")
//   client.Connect()
//   client.Subscribe("orders", handler)
//   client.Publish("orders", data)
//   client.Unsubscribe("orders")
//   client.Close()
//
// RESPONSABILIDADES:
// - Gerenciar conexão TCP com broker
// - Serializar/deserializar mensagens usando shared.protocol
// - Abstrair detalhes de rede da aplicação
// - Fornecer callbacks para mensagens recebidas
//
// PRÓXIMOS PASSOS (feature/pubsub-lib):
// - Implementar struct Client
// - Implementar Connect/Close
// - Implementar Publish
// - Implementar Subscribe/Unsubscribe
// - Implementar loop de leitura com goroutine
// - Implementar sistema de callbacks
package lib

// TODO: Implementar na branch feature/pubsub-lib
//
// Estrutura esperada:
//
// type MessageHandler func(topic string, payload []byte)
//
// type Client struct {
//     addr    string
//     conn    net.Conn
//     reader  *bufio.Reader
//     writer  *bufio.Writer
//     mu      sync.Mutex
//     handlers map[string]MessageHandler
//     done    chan struct{}
// }
//
// func NewClient(addr string) *Client
// func (c *Client) Connect() error
// func (c *Client) Close() error
// func (c *Client) Publish(topic string, payload any) error
// func (c *Client) Subscribe(topic string, handler MessageHandler) error
// func (c *Client) Unsubscribe(topic string) error
