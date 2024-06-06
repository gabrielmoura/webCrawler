# Nome do binário a ser gerado
BINARY=crawler

# Diretório do arquivo main.go
MAIN_DIR=cmd/crawler

# Opções de compilação
LDFLAGS=-ldflags="-s -w"

ALL:
	@echo "Comandos disponíveis:"
	@echo "  make build - Compila o binário"
	@echo "  make upx - Compila e compacta o binário com UPX"

# Regra padrão: build
build:
	go build $(LDFLAGS) -o $(BINARY) ./$(MAIN_DIR)

# Regra para build com UPX
upx: build
	upx --best $(BINARY)
