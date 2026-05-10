package main

import (
	"log"
	"net"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	port := getEnv("LB_PORT", ":8080")
	brokersStr := getEnv("BROKERS", "localhost:9000,localhost:9001")
	brokers := strings.Split(brokersStr, ",")

	balancer := NewBalancer(brokers)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[FATAL] Erro ao iniciar load balancer: %v", err)
	}
	defer listener.Close()

	log.Printf("[INFO] Load Balancer iniciado em %s", port)
	log.Printf("[INFO] Brokers disponíveis: %v", brokers)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] Erro ao aceitar conexão: %v", err)
			continue
		}
		go balancer.ProxyConnection(conn)
	}
}
