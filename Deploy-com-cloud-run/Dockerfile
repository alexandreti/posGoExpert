# Etapa de build
FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o app

# Etapa de execução com Alpine para suporte a TLS
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app /app

# Instala os certificados SSL para evitar erro de TLS
RUN apk --no-cache add ca-certificates

EXPOSE 8080
CMD ["/app/app"]


#docker build -t alexandreti/deploy-com-cloud-run -f Dockerfile .