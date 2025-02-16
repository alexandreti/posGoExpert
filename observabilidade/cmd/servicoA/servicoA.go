package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type CepRequest struct {
	CEP string `json:"cep"`
}

func ValidaCEP(cep string) (string, error) {
	cep = strings.ReplaceAll(cep, "-", "")
	if len(cep) != 8 {
		return "", fmt.Errorf("CEP inválido: deve conter 8 dígitos")
	}
	for _, r := range cep {
		if !unicode.IsDigit(r) {
			return "", fmt.Errorf("CEP inválido: deve conter apenas números")
		}
	}
	return cep, nil
}

func encaminhaParaServicoB(cep string) (*http.Response, error) {
	requestBody, _ := json.Marshal(CepRequest{CEP: cep})
	urlServicoB := os.Getenv("SERVICO_B_URL")
	return http.Post(urlServicoB+"/temperaturebycep", "application/json", bytes.NewBuffer(requestBody))
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

	cep, err := ValidaCEP(req.CEP)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	resp, err := encaminhaParaServicoB(cep)
	if err != nil {
		http.Error(w, "Erro ao chamar o Serviço B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler a resposta do Serviço B", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func main() {
	http.HandleFunc("/consulta", handler)
	fmt.Println("Serviço A rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
