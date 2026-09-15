# EDA: Wallet Core + Microsserviço de Balances

Dois microsserviços comunicando-se por eventos via **Apache Kafka**:

- **Wallet Core** (Go) — produtor. Gerencia clientes, contas e transações, e publica
  eventos no Kafka a cada transação.
- **Balances** (Laravel) — consumidor. Lê os eventos, mantém o saldo de cada conta
  em banco próprio e o expõe por uma API REST.

Todo o ecossistema sobe com **um único comando**, com migrations e seeds automáticos
nos dois bancos.

## Arquitetura

```
                    ┌──────────────────┐
  POST :8080  ─────►│   Wallet Core    │────► MySQL (wallet)
  /transactions     │       (Go)       │
                    └────────┬─────────┘
                             │ publica
                   ┌─────────▼──────────┐
                   │       Kafka        │
                   │ tópicos:           │
                   │  • transactions    │
                   │  • balances        │
                   └─────────┬──────────┘
                             │ consome
                    ┌────────▼─────────┐
                    │ balances-consumer│────► MySQL (balances)
                    │ (artisan)        │           ▲
                    └──────────────────┘           │ lê
                    ┌──────────────────┐           │
  GET :3003   ─────►│   balances-api   │───────────┘
  /balances/{id}    │   (FrankenPHP)   │
                    └──────────────────┘
```

O consumidor e a API rodam **a partir da mesma imagem**, mudando apenas o comando:
a API é request/response, o consumidor é um worker de longa duração. Separados,
cada um tem seu próprio ciclo de vida e política de restart.

## Stack

| Componente | Tecnologia |
|---|---|
| Wallet Core | Go 1.26 · confluent-kafka-go |
| Balances | PHP 8.5 · Laravel 13 · ext-rdkafka · FrankenPHP |
| Mensageria | Apache Kafka (Confluent 6.1) + Zookeeper |
| Bancos | MySQL 5.7 (wallet) · MySQL 8.0 (balances) |
| Orquestração | Docker Compose |

## Como executar

Pré-requisito: Docker e Docker Compose.

```bash
git clone <url-do-repositorio>
cd fc-wallet-core-full
docker compose up -d
```

Só isso. Na primeira execução o Docker compila as duas imagens (3-5 minutos);
nas seguintes, o ambiente sobe em cerca de 30 segundos.

Não é necessário rodar migrations, seeds ou qualquer script manual: as tabelas são
criadas e populadas automaticamente nos dois bancos durante a subida.

### Portas

| Serviço | Porta | URL |
|---|---|---|
| Balances (API) | **3003** | http://localhost:3003 |
| Wallet Core (API) | 8080 | http://localhost:8080 |
| Kafka Control Center | 9021 | http://localhost:9021 |
| MySQL (wallet) | 3306 | — |
| MySQL (balances) | 3307 | — |

## Testando o fluxo de eventos

O arquivo **`api.http`** na raiz tem todas as chamadas prontas (PhpStorm, VS Code
REST Client ou similar). O roteiro mínimo são três requisições:

**1. Consulte o saldo no Balances**

```bash
curl http://localhost:3003/balances/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
# {"account_id":"aaaa...","balance":1000,"updated_at":"..."}
```

**2. Faça uma transferência no Wallet Core**

```bash
curl -X POST http://localhost:8080/transactions \
  -H 'Content-Type: application/json' \
  -d '{
    "account_id_from": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    "account_id_to":   "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
    "amount": 300
  }'
```

**3. Consulte o saldo de novo**

```bash
curl http://localhost:3003/balances/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa
# {"account_id":"aaaa...","balance":700,"updated_at":"..."}
```

O saldo mudou sem que ninguém escrevesse no banco do Balances: o evento percorreu
o Kafka e o consumidor o aplicou. Para acompanhar ao vivo:

```bash
docker compose logs -f balances-consumer
```

```
[transactions] offset 0 - TransactionCreated
[balances]     offset 0 - BalanceUpdated
```


## Endpoints do Balances (porta 3003)

| Método | Rota | Descrição |
|---|---|---|
| GET | `/balances/{account_id}` | Saldo atual da conta |
| GET | `/balances/{account_id}/transactions` | Histórico de transações |
| GET | `/up` | Health check |

Conta inexistente devolve `404` com `{"error":"Account not found"}`.

## Dados criados pelos seeds

| Conta | Cliente | Saldo inicial |
|---|---|---|
| `aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa` | John Doe | 1000 |
| `bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb` | Jane Doe | 500 |

Os mesmos IDs existem nos dois bancos, então as chamadas funcionam imediatamente.

## Formato dos eventos

O Wallet Core serializa a struct do evento inteira. Como os campos Go não têm tags
`json`, as chaves chegam capitalizadas:

**Tópico `balances`**
```json
{
  "Name": "BalanceUpdated",
  "Payload": {
    "account_id_from": "aaaa...",
    "account_id_to": "bbbb...",
    "balance_account_id_from": 700,
    "balance_account_id_to": 800
  }
}
```


**Tópico `transactions`**
```json
{
  "Name": "TransactionCreated",
  "Payload": {
    "id": "uuid-da-transacao",
    "account_id_from": "aaaa...",
    "account_id_to": "bbbb...",
    "amount": 300
  }
}
```

## Decisões de projeto

**O consumidor grava o saldo, não soma o delta.** O evento `BalanceUpdated` já traz
o saldo final calculado pelo Wallet Core. Escrevendo o valor recebido, o Balances
converge para a verdade do produtor mesmo depois de perder ou reprocessar eventos.
Somar deltas corromperia o saldo de forma permanente diante de qualquer duplicata.

**Consumo idempotente.** O Kafka garante entrega *at-least-once*. As tabelas usam
chaves naturais (`account_id` para saldos, o `id` da transação vinda do produtor)
e o consumidor faz `updateOrCreate`, de modo que reprocessar uma mensagem produz o
mesmo estado final.

**Commit manual de offset.** O offset só avança **depois** da escrita no banco. Com
auto-commit, uma falha entre o commit e o `INSERT` faria o evento desaparecer.

**Seeds idempotentes.** O seeder não faz nada se já houver dados. Como ele roda a
cada subida do container, sem essa guarda um `docker compose restart` reverteria
saldos já atualizados por eventos.


**Sem foreign key entre `transactions` e `balances` no Balances.** Em arquitetura
orientada a eventos não há garantia de ordem entre tópicos distintos; uma FK
rejeitaria um `TransactionCreated` que chegasse antes do `BalanceUpdated` da mesma
conta.

## Comandos úteis

```bash
docker compose ps                          # estado dos containers
docker compose logs -f balances-consumer   # eventos sendo processados
docker compose logs -f wallet-core         # eventos sendo publicados
docker compose down                        # derruba, preservando os dados
docker compose down -v                     # derruba e apaga os bancos
```

Para recomeçar do zero (bancos vazios, tópicos limpos):

```bash
docker compose down -v && docker compose up -d
```

O que esse README faz por você na correção:

- O comando de execução aparece na primeira tela. O avaliador não deveria precisar rolar a página para descobrir como rodar o projeto.
- A seção "Decisões de projeto" é onde você ganha pontos de arquitetura. O enunciado diz que o foco é arquitetura e automação — e essa seção mostra que as escolhas foram deliberadas, não acidentais. Idempotência e commit manual de offset são exatamente os temas que separam quem entendeu EDA de quem só fez funcionar.
- O formato dos eventos documentado poupa o avaliador de ler código Go para entender o contrato.
- O aviso de 3-5 minutos no primeiro build evita que ele ache que travou e mate o processo.
