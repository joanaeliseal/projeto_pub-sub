package coletor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const awesomeAPIURL = "https://economia.awesomeapi.com.br/json/last"

type Cotacao struct {
	Compra       string
	Venda        string
	VariacaoPct  string
	MaximoDia    string
	MinimoDia    string
	AtualizadoEm string
}

type respostaAPI struct {
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	High       string `json:"high"`
	Low        string `json:"low"`
	PctChange  string `json:"pctChange"`
	CreateDate string `json:"create_date"`
}

func buscarPar(par string) (*Cotacao, error) {
	resp, err := http.Get(awesomeAPIURL + "/" + par)
	if err != nil {
		return nil, fmt.Errorf("erro de conexão com AwesomeAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AwesomeAPI retornou status %d para o par %s", resp.StatusCode, par)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]respostaAPI
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta da AwesomeAPI: %w", err)
	}

	for _, r := range result {
		return &Cotacao{
			Compra:       r.Bid,
			Venda:        r.Ask,
			VariacaoPct:  r.PctChange + "%",
			MaximoDia:    r.High,
			MinimoDia:    r.Low,
			AtualizadoEm: r.CreateDate,
		}, nil
	}
	return nil, fmt.Errorf("AwesomeAPI retornou resposta vazia para o par %s", par)
}

func BuscarDolar() (*Cotacao, error) {
	return buscarPar("USD-BRL")
}

func BuscarEuro() (*Cotacao, error) {
	return buscarPar("EUR-BRL")
}

func BuscarBitcoin() (*Cotacao, error) {
	return buscarPar("BTC-BRL")
}
