package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidaCEP(t *testing.T) {
	tests := []struct {
		name     string
		cepInput string
		expected string
		wantErr  bool
	}{
		{"CEP válido com hífen", "12345-678", "12345678", false},
		{"CEP válido sem hífen", "12345678", "12345678", false},
		{"CEP válido com hífen no meio", "1234-5678", "12345678", false},
		{"CEP inválido - caractere não numérico", "12a45-678", "", true},
		{"CEP com espaço", "12345 678", "", true},
		{"CEP com caracteres especiais", "12#45-678", "", true},
		{"CEP longo demais", "123456789", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidaCEP(tt.cepInput)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Esperava erro para o CEP '%s', mas não houve erro", tt.cepInput)
				}
			} else {
				if err != nil {
					t.Errorf("Não esperava erro para o CEP '%s', mas ocorreu: %v", tt.cepInput, err)
				}
				if result != tt.expected {
					t.Errorf("Esperava resultado '%s', mas obteve '%s'", tt.expected, result)
				}
			}
		})
	}
}

func TestConverterTemperaturas(t *testing.T) {
	tests := []struct {
		name      string
		tempC     float64
		expectedF float64
		expectedK float64
	}{
		{"Zero graus Celsius", 0, 32, 273.15},
		{"100 graus Celsius", 100, 212, 373.15},
		{"Temperatura negativa", -10, 14, 263.15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempF, tempK := converterTemperaturas(tt.tempC)
			if tempF != tt.expectedF {
				t.Errorf("Esperava %f Fahrenheit, mas obteve %f", tt.expectedF, tempF)
			}
			if tempK != tt.expectedK {
				t.Errorf("Esperava %f Kelvin, mas obteve %f", tt.expectedK, tempK)
			}
		})
	}
}

func TestConsultaCEP(t *testing.T) {
	mockResponse := `{"cep":"12345678","localidade":"São Paulo","uf":"SP"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "99999999") {
			http.Error(w, `{"erro": "true"}`, http.StatusNotFound)
			return
		}
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	oldURL := viacepURL
	viacepURL = server.URL + "/"
	defer func() { viacepURL = oldURL }()

	tests := []struct {
		name    string
		cep     string
		wantErr bool
	}{
		{"CEP encontrado", "12345678", false},
		{"CEP inexistente", "99999999", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConsultaCEP(tt.cep)
			if tt.wantErr && err == nil {
				t.Errorf("Esperava erro, mas não ocorreu")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Não esperava erro, mas ocorreu: %v", err)
			}
		})
	}
}

func TestConsultaTemperatura(t *testing.T) {
	mockResponse := `{"current":{"temp_c":25.5}}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	oldURL := weatherAPIURL
	weatherAPIURL = server.URL
	defer func() { weatherAPIURL = oldURL }()

	temp, err := ConsultaTemperatura("São Paulo", "fake-key")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if temp != 25.5 {
		t.Errorf("Esperava 25.5°C, mas obteve %f", temp)
	}
}

func TestTemperatureByCEPHandler(t *testing.T) {
	mockCEPResponse := `{"cep":"12345678","localidade":"São Paulo","uf":"SP"}`
	mockWeatherResponse := `{"current":{"temp_c":22.0}}`

	viacepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockCEPResponse))
	}))
	defer viacepServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockWeatherResponse))
	}))
	defer weatherServer.Close()

	oldViaCEP := viacepURL
	oldWeatherAPI := weatherAPIURL
	viacepURL = viacepServer.URL
	weatherAPIURL = weatherServer.URL
	defer func() {
		viacepURL = oldViaCEP
		weatherAPIURL = oldWeatherAPI
	}()

	req, _ := http.NewRequest("GET", "/temperaturebycep?cep=12345678", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := "fake-key"
		cep := "12345678"

		viaCep, _ := ConsultaCEP(cep)
		tempC, _ := ConsultaTemperatura(viaCep.Localidade, apiKey)

		tempF, tempK := converterTemperaturas(tempC)

		respData := Response{TempC: tempC, TempF: tempF, TempK: tempK}
		json.NewEncoder(w).Encode(respData)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Esperava código 200, mas recebeu %d", rr.Code)
	}

	expectedBody := `{"temp_C":22,"temp_F":71.6,"temp_K":295.15}`
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("Esperava %s, mas obteve %s", expectedBody, rr.Body.String())
	}
}
