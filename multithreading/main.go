package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilAPIResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	Erro        bool   `json:"erro"`
}

type APIResult struct {
	Source string
	Result string
}

func main() {
	// Verificar argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP>")
		os.Exit(1)
	}
	cep := os.Args[1]

	// Criar contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Canais para receber os resultados
	resultChannel := make(chan APIResult)

	// Iniciar goroutines para chamadas às APIs
	go fetchBrasilAPI(ctx, cep, resultChannel)
	go fetchViaCEP(ctx, cep, resultChannel)

	// Usar select para pegar o primeiro resultado ou detectar timeout
	select {
	case result := <-resultChannel:
		fmt.Printf("Resposta da API: %s\n%s\n", result.Source, result.Result)
	case <-ctx.Done():
		fmt.Println("Erro: Nenhuma API respondeu dentro do tempo limite de 1 segundo.")
	}
}

func fetchBrasilAPI(ctx context.Context, cep string, resultChannel chan<- APIResult) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resultChannel <- APIResult{Source: "BrasilAPI", Result: fmt.Sprintf("Erro ao criar requisição: %v", err)}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resultChannel <- APIResult{Source: "BrasilAPI", Result: fmt.Sprintf("Erro: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultChannel <- APIResult{Source: "BrasilAPI", Result: fmt.Sprintf("Erro: status code %d", resp.StatusCode)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChannel <- APIResult{Source: "BrasilAPI", Result: fmt.Sprintf("Erro ao ler resposta: %v", err)}
		return
	}

	var data BrasilAPIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultChannel <- APIResult{Source: "BrasilAPI", Result: fmt.Sprintf("Erro ao decodificar JSON: %v", err)}
		return
	}

	resultChannel <- APIResult{
		Source: "BrasilAPI",
		Result: fmt.Sprintf("  CEP: %s\n  Estado: %s\n  Cidade: %s\n  Bairro: %s\n  Rua: %s",
			data.CEP, data.State, data.City, data.Neighborhood, data.Street),
	}
}

func fetchViaCEP(ctx context.Context, cep string, resultChannel chan<- APIResult) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resultChannel <- APIResult{Source: "ViaCEP", Result: fmt.Sprintf("Erro ao criar requisição: %v", err)}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		resultChannel <- APIResult{Source: "ViaCEP", Result: fmt.Sprintf("Erro: %v", err)}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultChannel <- APIResult{Source: "ViaCEP", Result: fmt.Sprintf("Erro: status code %d", resp.StatusCode)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChannel <- APIResult{Source: "ViaCEP", Result: fmt.Sprintf("Erro ao ler resposta: %v", err)}
		return
	}

	var data ViaCEPResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		resultChannel <- APIResult{Source: "ViaCEP", Result: fmt.Sprintf("Erro ao decodificar JSON: %v", err)}
		return
	}

	if data.Erro {
		resultChannel <- APIResult{Source: "ViaCEP", Result: "CEP não encontrado."}
		return
	}

	resultChannel <- APIResult{
		Source: "ViaCEP",
		Result: fmt.Sprintf("  CEP: %s\n  Logradouro: %s\n  Complemento: %s\n  Bairro: %s\n  Localidade: %s\n  UF: %s",
			data.CEP, data.Logradouro, data.Complemento, data.Bairro, data.Localidade, data.UF),
	}
}
