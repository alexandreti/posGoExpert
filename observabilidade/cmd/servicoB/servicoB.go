package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type CepRequest struct {
	CEP string `json:"cep"`
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ViaCepResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       string `json:"erro,omitempty"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func converterTemperaturas(tempC float64) (float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15
	return tempF, tempK
}

func ConsultaCEP(cep string) (ViaCepResponse, error) {
	viacepURL := "https://viacep.com.br/ws/"
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

func ConsultaTemperatura(city string) (float64, error) {
	weatherAPIURL := "http://api.weatherapi.com/v1/current.json"
	apiKey := os.Getenv("API_KEY")
	encodedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("%s?key=%s&q=%s&aqi=no", weatherAPIURL, apiKey, encodedCity)

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

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req CepRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	json.Unmarshal(body, &req)

	viaCep, err := ConsultaCEP(req.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tempC, err := ConsultaTemperatura(viaCep.Localidade)
	if err != nil {
		http.Error(w, "Erro ao consultar temperatura", http.StatusInternalServerError)
		return
	}

	tempF, tempK := converterTemperaturas(tempC)
	respData := Response{
		City:  viaCep.Localidade,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respData)
}

func main() {
	http.HandleFunc("/temperaturebycep", handler)
	fmt.Println("Serviço B rodando na porta 8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
