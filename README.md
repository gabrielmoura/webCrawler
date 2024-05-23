# Um webcrawler Simples
## Em desenvolvimento

## Objetivo
Um webcrawler que possa ser usado para coletar informações de sites na clearnet, i2p e tor.
Simples, escrito em Go, com suporte a proxies, controle de concorrência e profundidade.
Os dados coletados podem ser salvos em um banco de dados mongodb, futuramente poderá ser implementado suporte a outros bancos de dados.

## Uso:

### I2P
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -proxy -proxyURL=http://localhost:4444 \
  -url=http://i2pforum.i2p \
  -tlds=i2p
```

### Tor
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -proxy -proxyURL=http://localhost:9050 \
  -url=http://zqktlwi4fecvo6ri.onion \
  -tlds=onion
```

### Clearnet
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -url=https://www.uol.com.br
```
