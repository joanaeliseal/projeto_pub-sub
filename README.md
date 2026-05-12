# Middleware Publish/Subscribe (Pub/Sub)

## Objetivo do Projeto

Implementação de um middleware Pub/Sub para a disciplina de Programação Distribuída. O sistema permite comunicação assíncrona entre publicadores e assinantes através de tópicos, usando sockets TCP e protocolo JSON newline-delimited.

## Funcionalidades

- Múltiplas conexões TCP simultâneas via goroutines
- Criação e remoção dinâmica de tópicos
- Operações de `publish`, `subscribe` e `unsubscribe`
- Bufferização de mensagens por tópico (tamanho configurável, padrão 100)
- Dispatch assíncrono: recebimento e encaminhamento são independentes
- Descarte de mensagens sem assinantes com notificação ao publicador
- Load balancer com roteamento determinístico por hash de tópico
- Biblioteca cliente (`lib`) que abstrai toda a comunicação TCP

## Estrutura do Projeto

```
projeto_pub-sub/
├── shared/
│   └── protocol.go              # Protocolo JSON: structs, encode/decode, constantes
├── broker/
│   ├── main.go                  # Servidor TCP, accept loop, handlers
│   ├── broker.go                # Gerenciador central de tópicos
│   ├── topico.go                # Tópico com buffer e goroutine dispatcher
│   └── cliente.go               # Representação de conexão TCP de cliente
├── loadbalancer/
│   ├── main.go                  # Servidor TCP do load balancer
│   └── balancer.go              # Roteamento determinístico por hash de tópico
├── lib/
│   └── client.go                # Biblioteca cliente para aplicações
└── exemplos/
    ├── publisher-cambio/        # Publica: dolar, euro (AwesomeAPI)
    ├── publisher-cripto/        # Publica: bitcoin (AwesomeAPI), selic (BCB)
    ├── subscriber-importadores/ # Assina: dolar, euro
    └── subscriber-investidores/ # Assina: bitcoin, selic
```

## Protocolo de Comunicação

Cada mensagem é um JSON completo terminado com `\n` (newline-delimited):

```json
{"type":"subscribe","topic":"dolar"}
{"type":"publish","topic":"dolar","payload":{"compra":"5.72","venda":"5.73"}}
{"type":"unsubscribe","topic":"dolar"}
```

O broker responde com:
```json
{"type":"ack","success":true,"topic":"dolar"}
{"type":"error","success":false,"code":"NO_SUBSCRIBERS","message":"..."}
{"type":"error","success":false,"code":"BUFFER_FULL","message":"..."}
```

Mensagens entregues a subscribers chegam com:
```json
{"type":"message","topic":"dolar","payload":{...},"timestamp":"2026-..."}
```

## Balanceamento de Carga (Topic-Hash)

O load balancer usa hash determinístico do nome do tópico para escolher o broker:

```
broker = fnv32(topic) % número_de_brokers
```

Isso garante que publisher e subscriber do mesmo tópico sempre se encontrem no mesmo broker, sem necessidade de sincronização entre brokers. A biblioteca cliente (`lib`) abstrai esse detalhe: a aplicação só conhece o endereço do load balancer.

## Cenário de Teste

### Descrição

Simulação de um sistema de informações financeiras em tempo real, consumindo dados públicos do mercado brasileiro via APIs abertas e sem autenticação.

**APIs utilizadas:**
- [AwesomeAPI](https://docs.awesomeapi.com.br/) — cotações de câmbio e criptomoedas
- [API do Banco Central do Brasil](https://dadosabertos.bcb.gov.br/) — taxa SELIC

| Tópico    | Publisher              | Subscriber               | Dados                            |
|-----------|------------------------|--------------------------|----------------------------------|
| `dolar`   | publisher-cambio       | subscriber-importadores  | Cotação USD/BRL em tempo real    |
| `euro`    | publisher-cambio       | subscriber-importadores  | Cotação EUR/BRL em tempo real    |
| `bitcoin` | publisher-cripto       | subscriber-investidores  | Preço do BTC em reais            |
| `selic`   | publisher-cripto       | subscriber-investidores  | Taxa SELIC (Banco Central)       |

- **subscriber-importadores**: empresa importadora monitorando câmbio para precificação
- **subscriber-investidores**: painel de investidores acompanhando cripto e juros

### Variáveis de Ambiente

Copie `.env.example` para `.env` e ajuste conforme necessário:

```bash
cp .env.example .env
```

| Variável          | Padrão                           | Descrição                        |
|-------------------|----------------------------------|----------------------------------|
| `PORT`            | `:9000`                          | Porta do broker                  |
| `LB_PORT`         | `:8080`                          | Porta do load balancer           |
| `BROKERS`         | `localhost:9000,localhost:9001`  | Lista de brokers (vírgula)       |
| `LB_ADDR`         | `localhost:8080`                 | Endereço do LB nos exemplos      |
| `BUFFER_SIZE`     | `100`                            | Tamanho do buffer por tópico     |
| `INTERVALO_CAMBIO`| `10`                             | Segundos entre buscas de câmbio  |
| `INTERVALO_BITCOIN`| `15`                            | Segundos entre buscas de BTC     |
| `INTERVALO_SELIC` | `60`                             | Segundos entre buscas da SELIC   |

### Como Executar

#### Opção 1 — Docker Compose (recomendado)

Sobe todos os serviços com um único comando:

```bash
docker compose up --build
```

Para acompanhar apenas os logs dos subscribers (mais limpo para apresentação):

```bash
docker compose logs -f subscriber-importadores subscriber-investidores
```

Para encerrar tudo:

```bash
docker compose down
```

---

#### Opção 2 — Terminais separados (sem Docker)

**1. Brokers** (dois terminais):
```bash
# Terminal 1 — Broker A
go run ./broker

# Terminal 2 — Broker B
PORT=:9001 go run ./broker
```

**2. Load Balancer** (um terminal):
```bash
# Terminal 3
# Porta padrão: 8080 (a 8000 é usada pelo Docker Desktop)
BROKERS=localhost:9000,localhost:9001 go run ./loadbalancer
```

**3. Subscribers** (antes dos publishers — precisam estar inscritos para receber):
```bash
# Terminal 4
go run ./exemplos/subscriber-importadores

# Terminal 5
go run ./exemplos/subscriber-investidores
```

**4. Publishers** (iniciam busca nas APIs e publicam em loop):
```bash
# Terminal 6
go run ./exemplos/publisher-cambio

# Terminal 7
go run ./exemplos/publisher-cripto
```

### Saída Esperada

**subscriber-importadores:**
```
[INFO] Importadora aguardando cotações de: dolar, euro
[DÓLAR]  compra R$5.7192 | venda R$5.7202 | var +0.01% | máx R$5.7362 | mín R$5.7048
[EURO]   compra R$6.2100 | venda R$6.2130 | var -0.05% | máx R$6.2500 | mín R$6.1900
```

**subscriber-investidores:**
```
[INFO] Painel de Investidores aguardando: bitcoin, selic
[BITCOIN] R$561234.00 | var +2.3% | máx R$570000.00 | mín R$550000.00
[SELIC]   13.75% a.a. (referência: 09/05/2026)
```

**Load Balancer:**
```
[INFO] Load Balancer iniciado em :8000
[INFO] Brokers disponíveis: [localhost:9000 localhost:9001]
[LB] 127.0.0.1:XXXXX → tópico 'dolar'   → localhost:9000
[LB] 127.0.0.1:XXXXX → tópico 'euro'    → localhost:9001
[LB] 127.0.0.1:XXXXX → tópico 'bitcoin' → localhost:9000
[LB] 127.0.0.1:XXXXX → tópico 'selic'   → localhost:9001
```

### Teste sem Load Balancer

Para testar apenas o broker sem o LB:
```bash
# Terminal 1 — broker direto na porta padrão
go run ./broker

# Nos exemplos, aponte para o broker:
LB_ADDR=localhost:9000 go run ./exemplos/publisher-cambio
LB_ADDR=localhost:9000 go run ./exemplos/subscriber-importadores
```

## Informações

- **Disciplina:** Programação Distribuída 2026.1
- **Professor:** Ruan Delgado Gomes
- **Alunas:** Joana Elise e Maria Eduarda Vitorino
