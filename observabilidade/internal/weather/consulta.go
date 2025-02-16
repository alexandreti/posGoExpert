package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
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
