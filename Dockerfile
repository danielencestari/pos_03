# Build stage - Compilação da aplicação
FROM golang:1.24-alpine AS builder

# Instalar ferramentas necessárias
RUN apk add --no-cache git

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
# CGO_ENABLED=0: Compilação estática
# GOOS=linux: Target Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/ordersystem/main.go

# Production stage - Imagem final minimalista
FROM alpine:latest

# Instalar certificados SSL para HTTPS
RUN apk --no-cache add ca-certificates

# Criar usuário não-root para segurança
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /root/

# Copiar binário da aplicação do estágio anterior
COPY --from=builder /app/main .

# Copiar arquivo de configuração
COPY --from=builder /app/.env .

# Alterar propriedade dos arquivos para o usuário não-root
RUN chown -R appuser:appgroup /root/

# Mudar para usuário não-root
USER appuser

# Expor portas que a aplicação utiliza
EXPOSE 8000 8080 50051

# Comando para executar a aplicação
CMD ["./main"] 