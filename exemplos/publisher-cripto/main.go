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

func main() {
	addr := getEnv("LB_ADDR", "localhost:8080")

	intervaloBitcoinSeg, err := strconv.Atoi(getEnv("INTERVALO_BITCOIN", "15"))
	if err != nil {
		intervaloBitcoinSeg = 15
	}

	intervaloSelicSeg, err := strconv.Atoi(getEnv("INTERVALO_SELIC", "60"))
	if err != nil {
		intervaloSelicSeg = 60
	}

	clientBitcoin := lib.NewClient(addr)
	if err := clientBitcoin.Connect(); err != nil {
		log.Fatalf("[FATAL] Erro ao conectar (bitcoin): %v", err)
	}
	defer clientBitcoin.Close()

	clientSelic := lib.NewClient(addr)
	if err := clientSelic.Connect(); err != nil {
		log.Fatalf("[FATAL] Erro ao conectar (selic): %v", err)
	}
	defer clientSelic.Close()

	log.Printf("[INFO] Publisher Cripto/Juros iniciado")
	log.Printf("[INFO]   bitcoin → a cada %ds | selic → a cada %ds",
		intervaloBitcoinSeg, intervaloSelicSeg)

	publicarBitcoin := func() {
		c, err := coletor.BuscarBitcoin()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
		payload := map[string]any{
			"preco_brl":     c.Compra,
			"variacao_pct":  c.VariacaoPct,
			"maximo_dia":    c.MaximoDia,
			"minimo_dia":    c.MinimoDia,
			"atualizado_em": c.AtualizadoEm,
		}
		if err := clientBitcoin.Publish("bitcoin", payload); err != nil {
			log.Printf("[ERROR] Falha ao publicar bitcoin: %v", err)
			return
		}
		log.Printf("[INFO] BTC-BRL → R$%s | var %s", c.Compra, c.VariacaoPct)
	}

	publicarSelic := func() {
		s, err := coletor.BuscarSelic()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
		payload := map[string]any{
			"taxa_anual_pct":  s.TaxaAnualPct,
			"data_referencia": s.DataReferencia,
			"fonte":           "Banco Central do Brasil",
		}
		if err := clientSelic.Publish("selic", payload); err != nil {
			log.Printf("[ERROR] Falha ao publicar selic: %v", err)
			return
		}
		log.Printf("[INFO] SELIC → %.2f%% a.a. (ref: %s)", s.TaxaAnualPct, s.DataReferencia)
	}

	publicarBitcoin()
	publicarSelic()

	tickerBitcoin := time.NewTicker(time.Duration(intervaloBitcoinSeg) * time.Second)
	tickerSelic := time.NewTicker(time.Duration(intervaloSelicSeg) * time.Second)
	defer tickerBitcoin.Stop()
	defer tickerSelic.Stop()

	for {
		select {
		case <-tickerBitcoin.C:
			publicarBitcoin()
		case <-tickerSelic.C:
			publicarSelic()
		case <-clientBitcoin.Done():
			return
		}
	}
}
