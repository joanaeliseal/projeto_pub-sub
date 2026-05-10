package main

import (
	"encoding/json"
	"log"
	"os"

	"pubsub/lib"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	addr := getEnv("LB_ADDR", "localhost:8080")

	clientDolar := lib.NewClient(addr)
	if err := clientDolar.Connect(); err != nil {
		log.Fatalf("[FATAL] Erro ao conectar (dolar): %v", err)
	}
	defer clientDolar.Close()

	clientEuro := lib.NewClient(addr)
	if err := clientEuro.Connect(); err != nil {
		log.Fatalf("[FATAL] Erro ao conectar (euro): %v", err)
	}
	defer clientEuro.Close()

	if err := clientDolar.Subscribe("dolar", handleDolar); err != nil {
		log.Fatalf("[FATAL] Erro ao subscrever em dolar: %v", err)
	}

	if err := clientEuro.Subscribe("euro", handleEuro); err != nil {
		log.Fatalf("[FATAL] Erro ao subscrever em euro: %v", err)
	}

	log.Println("[INFO] Importadora aguardando cotações de: dolar, euro")

	<-clientDolar.Done()
}

func handleDolar(topic string, payload json.RawMessage) {
	var d map[string]any
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("[ERROR] Payload inválido: %v", err)
		return
	}
	log.Printf("[DÓLAR]  compra R$%v | venda R$%v | var %v | máx R$%v | mín R$%v",
		d["compra"], d["venda"], d["variacao_pct"], d["maximo_dia"], d["minimo_dia"])
}

func handleEuro(topic string, payload json.RawMessage) {
	var d map[string]any
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("[ERROR] Payload inválido: %v", err)
		return
	}
	log.Printf("[EURO]   compra R$%v | venda R$%v | var %v | máx R$%v | mín R$%v",
		d["compra"], d["venda"], d["variacao_pct"], d["maximo_dia"], d["minimo_dia"])
}
