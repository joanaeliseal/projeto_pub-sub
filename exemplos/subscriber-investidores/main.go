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

	client := lib.NewClient(addr)
	defer client.Close()

	if err := client.Subscribe("bitcoin", handleBitcoin); err != nil {
		log.Fatalf("[FATAL] Erro ao subscrever em bitcoin: %v", err)
	}

	if err := client.Subscribe("selic", handleSelic); err != nil {
		log.Fatalf("[FATAL] Erro ao subscrever em selic: %v", err)
	}

	log.Println("[INFO] Painel de Investidores aguardando: bitcoin, selic")

	<-client.Done()
}

func handleBitcoin(topic string, payload json.RawMessage) {
	var d map[string]any
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("[ERROR] Payload inválido: %v", err)
		return
	}
	log.Printf("[BITCOIN] R$%v | var %v | máx R$%v | mín R$%v",
		d["preco_brl"], d["variacao_pct"], d["maximo_dia"], d["minimo_dia"])
}

func handleSelic(topic string, payload json.RawMessage) {
	var d map[string]any
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("[ERROR] Payload inválido: %v", err)
		return
	}
	log.Printf("[SELIC]   %.2f%% a.a. (referência: %v)",
		d["taxa_anual_pct"], d["data_referencia"])
}
