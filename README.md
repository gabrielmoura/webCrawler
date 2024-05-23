# Um webcrawler Simples
## Em desenvolvimento

## Objetivo
Um webcrawler que possa ser usado para coletar informações de sites na clearnet, i2p e tor.
Simples, escrito em Go, com suporte a proxies, controle de concorrência e profundidade.
Os dados coletados podem ser salvos em um banco de dados mongodb, futuramente poderá ser implementado suporte a outros bancos de dados.

## Uso:
É possível usar as flags para configurar o crawler, ou usar um arquivo de configuração em yaml.
Um exemplo de configuração pode ser encontrado em [config.yml](example_config.yml).

### Uso com arquivo de configuração:

```yml
maxConcurrency: 10
maxDepth: 50
proxy:
  proxyURL: http://localhost:4444
url: http://i2pforum.i2p
tlds:
  - i2p
```

O comando abaixo irá executar o crawler procurando as configurações do arquivo `config.yml` nos diretórios:
- `/etc/crw`
- `/opt/crw`
- `.` (diretório atual)
```bash
./crawler -config
```


### Uso com Flags:
#### I2P
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -proxy -proxyURL=http://localhost:4444 \
  -url=http://i2pforum.i2p \
  -tlds=i2p
```

#### Tor
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -proxy -proxyURL=http://localhost:9050 \
  -url=http://zqktlwi4fecvo6ri.onion \
  -tlds=onion
```

#### Clearnet
```bash
./crawler -maxConcurrency=10 \
  -maxDepth=50 \
  -url=https://www.uol.com.br
```
