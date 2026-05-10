package coletor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const bcbSelicURL = "https://api.bcb.gov.br/dados/serie/bcdata.sgs.11/dados/ultimos/1?formato=json"

type DadoSelic struct {
	TaxaAnualPct   float64
	DataReferencia string
}

func BuscarSelic() (*DadoSelic, error) {
	resp, err := http.Get(bcbSelicURL)
	if err != nil {
		return nil, fmt.Errorf("erro de conexão com API do BCB: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API do BCB retornou status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dados []struct {
		Data  string `json:"data"`
		Valor string `json:"valor"`
	}
	if err := json.Unmarshal(body, &dados); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta do BCB: %w", err)
	}

	if len(dados) == 0 {
		return nil, fmt.Errorf("BCB retornou lista vazia")
	}

	taxa, err := strconv.ParseFloat(dados[0].Valor, 64)
	if err != nil {
		return nil, fmt.Errorf("valor da SELIC inválido (%q): %w", dados[0].Valor, err)
	}

	return &DadoSelic{
		TaxaAnualPct:   taxa,
		DataReferencia: dados[0].Data,
	}, nil
}
