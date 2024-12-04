package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	// Canais para receber os resultados
	resultChannel := make(chan APIResult)

	// Iniciar goroutines para chamadas às APIs
	go fetchViaCEP(cep, resultChannel)
	go fetchBrasilAPI(cep, resultChannel)

	// Usar select para pegar o primeiro resultado
	select {
	case result := <-resultChannel:
		fmt.Printf("Resposta da API: %s\n%s\n", result.Source, result.Result)
	}
}

func fetchBrasilAPI(cep string, resultChannel chan<- APIResult) {
	//time.Sleep(time.Second)
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
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

func fetchViaCEP(cep string, resultChannel chan<- APIResult) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
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
