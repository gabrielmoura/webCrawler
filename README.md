# Um Webcrawler Simples

## Em Desenvolvimento

## Objetivo
Este projeto tem como objetivo criar um webcrawler simples, escrito em Go, que possa ser usado para coletar informações de sites na clearnet, I2P e Tor. O crawler oferece suporte a proxies, controle de concorrência e definição de profundidade. Os dados coletados podem ser salvos em um banco de dados PostgreSQL, com planos futuros para uma interface gráfica.

## Uso
Você pode configurar o crawler usando flags de linha de comando ou um arquivo de configuração em YAML. Um exemplo de configuração pode ser encontrado em [example_config.yml](example_config.yml).

### Uso com Arquivo de Configuração

Exemplo de configuração (`config.yml`):

```yml
maxConcurrency: 10
maxDepth: 50
proxy:
  proxyURL: http://localhost:4444
url: http://i2pforum.i2p
tlds:
  - i2p
```

Para executar o crawler usando um arquivo de configuração, utilize o comando abaixo. O crawler procurará o arquivo `config.yml` nos seguintes diretórios:
- `/etc/crw`
- `/opt/crw`
- Diretório atual (`.`)

```bash
./crawler -config
```

### Uso com Flags

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
## Consumo de Recursos
O Crawler pode consumir mais ou menos recursos conforme as configurações de concorrência e profundidade.
Recomenda-se ajustar essas configurações conforme a capacidade do servidor e a quantidade de dados que deseja coletar.
Alguns comandos úteis para monitorar o consumo de recursos são:

#### Ver Memória
```bash
pmap -x $(pidof crawler)| awk '/kB total/ {print $5, $2}'
```

#### Ver CPU
```bash
ps aux | grep crawler | awk '{print $11, $3}'
```

## Contribuição
Sinta-se à vontade para contribuir com o desenvolvimento deste projeto. Para relatar problemas ou sugerir melhorias, abra uma issue ou envie um pull request.
