# Estágio 1: Build
# Usa a imagem oficial do Go como base para compilar a aplicação.
FROM golang:latest AS builder

# Define o diretório de trabalho dentro do container.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências.
COPY go.mod go.sum ./
# Baixa as dependências do projeto.
RUN go mod download

# Copia todo o código-fonte da aplicação para o container.
COPY . .

# Compila a aplicação.
# -o /app/main: Especifica o nome e o local do arquivo de saída.
# CGO_ENABLED=0: Desabilita CGO para criar um binário estático.
# GOOS=linux GOARCH=amd64: Garante que o build seja para Linux.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .

# Estágio 2: Produção
# Usa uma imagem mínima para a versão final, reduzindo o tamanho e a superfície de ataque.
FROM alpine:latest

WORKDIR /app

# Copia apenas o binário compilado do estágio de build.
COPY --from=builder /app/main .

# Expõe a porta que a API (Gin) irá escutar.
EXPOSE 8080

# Comando para executar a aplicação quando o container iniciar.
CMD ["/app/main"]