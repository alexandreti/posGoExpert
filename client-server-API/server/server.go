package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite3", "./currency.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao abrir o banco de dados: %v\n", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS currency (id INTEGER PRIMARY KEY, timestamp INTEGER, create_date TEXT, bid TEXT)")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar a tabela currency: %v\n", err)
		return
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()
		log.Println("Request iniciada")
		defer log.Println("Request finalizada")

		bid, err := fetchAndStoreCurrency(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{"bid": bid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("O Servidor está rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}

func fetchAndStoreCurrency(db *sql.DB) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout da requisição ao site economia.awesomeapi.com.br, alcançado!")
		}
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var currencyResponse CurrencyResponse
	err = json.Unmarshal(body, &currencyResponse)
	if err != nil {
		return "", err
	}

	bid := currencyResponse.USDBRL.Bid

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = db.ExecContext(
		ctx,
		"INSERT INTO currency (timestamp, create_date, bid) VALUES (?, ?, ?)",
		currencyResponse.USDBRL.Timestamp,
		currencyResponse.USDBRL.CreateDate,
		bid,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout para gravação dos dados no banco de dados alcançado!")
		}
		return "", err
	}

	return bid, nil
}
