FROM golang:1.23 AS builder
WORKDIR /app

# Copia os arquivos necessários
COPY go.* ./
RUN go mod download

COPY ./cmd/servicoB ./cmd/servicoB
COPY ./internal ./internal

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/servicoB ./cmd/servicoB/servicoB.go

# Imagem final
FROM alpine:latest
WORKDIR /app

# Copia o binário compilado
COPY --from=builder /app/servicoB /app/servicoB

# Define permissões de execução
RUN chmod +x /app/servicoB

# Define a porta
EXPOSE 8082

# Comando de execução
CMD ["/app/servicoB"]
