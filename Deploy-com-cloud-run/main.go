package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"
)

var viacepURL = "https://viacep.com.br/ws/"
var weatherAPIURL = "http://api.weatherapi.com/v1/current.json"

// Estrutura da resposta que nossa rota retornará.
type Response struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// Estrutura que representa a resposta da API do ViaCEP.
type ViaCepResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       string `json:"erro,omitempty"`
}

// Estrutura que representa a resposta da WeatherAPI.
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// Converte a temperatura de Celsius para Fahrenheit e Kelvin.
func converterTemperaturas(tempC float64) (float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK
}

// Carrega a API key a partir da variável de ambiente API_KEY".
func loadAPIKey() string {
	apiKey := os.Getenv("API_KEY")
	return apiKey
}

// ConsultaCEP realiza a requisição à API do ViaCEP e retorna a cidade.
func ConsultaCEP(cep string) (ViaCepResponse, error) {
	requestURL := fmt.Sprintf("%s%s/json/", viacepURL, cep)
	resp, err := http.Get(requestURL)
	if err != nil {
		return ViaCepResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ViaCepResponse{}, fmt.Errorf("erro na requisição do ViaCEP: status code %d", resp.StatusCode)
	}

	var viaCepResp ViaCepResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCepResp); err != nil {
		return ViaCepResponse{}, err
	}

	if viaCepResp.Erro == "true" {
		return ViaCepResponse{}, errors.New("can not find zipcode")
	}

	return viaCepResp, nil
}

// ConsultaTemperatura obtém a temperatura atual da cidade usando a WeatherAPI.
func ConsultaTemperatura(city, apiKey string) (float64, error) {
	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", weatherAPIURL, apiKey, encodedCity) // Alteração aqui

	resp, err := http.Get(requestURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("erro na requisição do WeatherAPI: status code %d", resp.StatusCode)
	}

	var weatherResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return 0, err
	}

	return weatherResp.Current.TempC, nil
}

// ValidaCEP verifica se o CEP tem 8 dígitos numéricos.
func ValidaCEP(cep string) (string, error) {
	cep = strings.ReplaceAll(cep, "-", "")
	if len(cep) != 8 {
		return "", errors.New("CEP inválido: deve conter 8 dígitos")
	}
	for _, r := range cep {
		if !unicode.IsDigit(r) {
			return "", errors.New("CEP inválido: deve conter apenas números")
		}
	}
	return cep, nil
}

func main() {
	apiKey := loadAPIKey()

	http.HandleFunc("/temperaturebycep", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request iniciada")
		defer log.Println("Request finalizada")

		cep := r.URL.Query().Get("cep")
		cep, err := ValidaCEP(cep)
		if err != nil {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		viaCep, err := ConsultaCEP(cep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if viaCep.Localidade == "" {
			http.Error(w, "Localidade não encontrada no ViaCEP", http.StatusInternalServerError)
			return
		}

		tempC, err := ConsultaTemperatura(viaCep.Localidade, apiKey)
		if err != nil {
			http.Error(w, "erro ao consultar a temperatura: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tempF, tempK := converterTemperaturas(tempC)

		respData := Response{
			TempC: tempC,
			TempF: tempF,
			TempK: tempK,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(respData); err != nil {
			http.Error(w, "Erro ao gerar resposta", http.StatusInternalServerError)
		}
	})

	fmt.Println("O Servidor está rodando na porta 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
