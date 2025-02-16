package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type CepRequest struct {
	CEP string `json:"cep"`
}

func initTracer(serviceName string) func() {
	zipkinURL := os.Getenv("OTEL_EXPORTER_ZIPKIN_ENDPOINT")
	if zipkinURL == "" {
		zipkinURL = "http://localhost:9411/api/v2/spans"
	}
	exporter, err := zipkin.New(zipkinURL)
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background()) }
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

func encaminhaParaServicoB(ctx context.Context, cep string) (*http.Response, error) {
	tracer := otel.Tracer("servicoA")
	ctx, span := tracer.Start(ctx, "EncaminhaParaServicoB")
	defer span.End()

	requestBody, _ := json.Marshal(CepRequest{CEP: cep})
	urlServicoB := os.Getenv("SERVICO_B_URL")
	req, err := http.NewRequestWithContext(ctx, "POST", urlServicoB+"/temperaturebycep", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer("servicoA").Start(r.Context(), "HandlerServicoA")
	defer span.End()

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

	resp, err := encaminhaParaServicoB(ctx, cep)
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
	shutdown := initTracer("servicoA")
	defer shutdown()

	http.HandleFunc("/consulta", handler)
	fmt.Println("Serviço A rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
