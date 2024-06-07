# Configuração por Flags

- maxConcurrency: Número máximo de execuções simultâneas.
- maxDepth: Profundidade máxima de execução.
- config: Arquivo de configuração.
- proxy: Proxy para requisições.
- proxyURL: URL do proxy.
- url: URL do site.
- mem: Salvar cache apenas na memória.
- tlds: Lista de TLDs para serem usadas.
- postgresURI: URI de conexão com o banco de dados PostgreSQL.
- userAgent: User-Agent para requisições.

## Exemplo de uso

```bash
./crawler -maxConcurrency 10 -maxDepth 5 \
-proxy -proxyURL http://localhost:4444 \
-url https://example.com -tlds com,org,net \
-postgresURI postgres://user:password@localhost:5432/db \
-userAgent "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"
```