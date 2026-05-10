// Package main é o ponto de entrada do Load Balancer.
//
// ARQUITETURA:
// O Load Balancer distribui conexões entre múltiplos brokers.
// Funciona como um proxy TCP simples.
//
// ESTRATÉGIA (escolher uma):
// - Round-robin: distribui conexões sequencialmente entre brokers
// - Hash de tópico: mesmo tópico sempre vai para o mesmo broker
//
// IMPORTANTE (conforme CLAUDE.md):
// NÃO implementar:
// - Sincronização de estado entre brokers
// - Replicação de mensagens
// - Consenso distribuído
// - Failover complexo
//
// EXECUÇÃO:
//   go run ./loadbalancer
//
// PRÓXIMOS PASSOS (feature/load-balancer):
// - Implementar estratégia round-robin
// - Implementar proxy de conexões
// - Configurar lista de brokers
package main

import (
	"log"
)

func main() {
	log.Println("[INFO] Load Balancer iniciando...")

	// TODO: Implementar na branch feature/load-balancer
	// 1. Carregar configuração de brokers
	// 2. Iniciar listener TCP
	// 3. Para cada conexão, selecionar broker (round-robin ou hash)
	// 4. Fazer proxy da conexão

	log.Println("[INFO] Load Balancer finalizado")
}
