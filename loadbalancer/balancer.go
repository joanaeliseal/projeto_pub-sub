// balancer.go contém a lógica de balanceamento de carga.
//
// ARQUITETURA:
// O Balancer mantém uma lista de brokers disponíveis
// e seleciona qual broker atenderá cada nova conexão.
//
// ESTRATÉGIAS SUPORTADAS:
// 1. Round-robin (recomendado para simplicidade)
//    - Distribui conexões sequencialmente
//    - Fácil de implementar e entender
//
// 2. Hash de tópico (alternativa)
//    - Garante que mesmo tópico vai para mesmo broker
//    - Requer extração do tópico da primeira mensagem
//
// IMPORTANTE:
// Esta é uma implementação SIMPLES para fins acadêmicos.
// NÃO é um sistema de produção com failover real.
//
// PRÓXIMOS PASSOS (feature/load-balancer):
// - Implementar struct Balancer
// - Implementar AddBroker
// - Implementar SelectBroker (round-robin)
// - Implementar health check básico (opcional)
package main

// TODO: Implementar na branch feature/load-balancer
//
// Estrutura esperada:
//
// type Balancer struct {
//     mu      sync.Mutex
//     brokers []string  // Lista de endereços "host:port"
//     current int       // Índice atual para round-robin
// }
//
// func NewBalancer() *Balancer
// func (b *Balancer) AddBroker(addr string)
// func (b *Balancer) SelectBroker() string  // Round-robin
// func (b *Balancer) ProxyConnection(clientConn net.Conn)
