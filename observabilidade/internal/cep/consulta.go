package cep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ViaCepResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	Erro       string `json:"erro,omitempty"`
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
		return ViaCepResponse{}, fmt.Errorf("Erro na requisição do ViaCEP: status code %d", resp.StatusCode)
	}

	var viaCepResp ViaCepResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCepResp); err != nil {
		return ViaCepResponse{}, err
	}

	if viaCepResp.Erro == "true" {
		return ViaCepResponse{}, errors.New("Can not find zipcode")
	}

	return viaCepResp, nil
}
