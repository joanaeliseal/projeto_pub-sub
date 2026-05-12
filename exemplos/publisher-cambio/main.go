package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"pubsub/exemplos/coletor"
	"pubsub/lib"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func cotacaoParaPayload(c *coletor.Cotacao) map[string]any {
	return map[string]any{
		"compra":        c.Compra,
		"venda":         c.Venda,
		"variacao_pct":  c.VariacaoPct,
		"maximo_dia":    c.MaximoDia,
		"minimo_dia":    c.MinimoDia,
		"atualizado_em": c.AtualizadoEm,
	}
}

func main() {
	addr := getEnv("LB_ADDR", "localhost:8080")
	intervaloSeg, err := strconv.Atoi(getEnv("INTERVALO_CAMBIO", "10"))
	if err != nil {
		intervaloSeg = 10
	}
	intervalo := time.Duration(intervaloSeg) * time.Second

	client := lib.NewClient(addr)
	defer client.Close()

	log.Printf("[INFO] Publisher Câmbio iniciado — publicando dolar e euro a cada %v", intervalo)

	ciclo := func() {
		if c, err := coletor.BuscarDolar(); err != nil {
			log.Printf("[ERROR] %v", err)
		} else if err := client.Publish("dolar", cotacaoParaPayload(c)); err != nil {
			log.Printf("[ERROR] Falha ao publicar dolar: %v", err)
		} else {
			log.Printf("[INFO] Dólar → compra R$%s | var %s", c.Compra, c.VariacaoPct)
		}

		if c, err := coletor.BuscarEuro(); err != nil {
			log.Printf("[ERROR] %v", err)
		} else if err := client.Publish("euro", cotacaoParaPayload(c)); err != nil {
			log.Printf("[ERROR] Falha ao publicar euro: %v", err)
		} else {
			log.Printf("[INFO] Euro  → compra R$%s | var %s", c.Compra, c.VariacaoPct)
		}
	}

	ciclo()

	ticker := time.NewTicker(intervalo)
	defer ticker.Stop()

	for range ticker.C {
		ciclo()
	}
}
