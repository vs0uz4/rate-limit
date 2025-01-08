# Rate Limiter

> [!IMPORTANT]
> Para poder executar o projeto contido neste repositório é necessário que se tenha o Docker instalado no computador. Para maiores informações siga o site <https://www.docker.com/>

- [Rate Limiter](#rate-limiter)
  - [Desafio GoLang Pós GoExpert - Rate Limiter](#desafio-golang-pós-goexpert---rate-limiter)
    - [Requisitos do Desafio](#requisitos-do-desafio)
    - [Exemplos](#exemplos)
    - [Entrega](#entrega)
  - [O que é o Rate Limiter](#o-que-é-o-rate-limiter)
  - [Funcionalidades](#funcionalidades)
  - [Requisitos](#requisitos)
    - [Tecnologias](#tecnologias)
    - [Configurações](#configurações)
    - [Exemplo de `TOKEN_LIMITS`](#exemplo-de-token_limits)
    - [Endpoints](#endpoints)
  - [Testes Automatizados](#testes-automatizados)
  - [Executando o Projeto](#executando-o-projeto)
    - [Passos](#passos)
  - [Exemplos de CURL](#exemplos-de-curl)
    - [Sem token](#sem-token)
    - [Com token](#com-token)
    - [Com token inválido/inexistente](#com-token-inválidoinexistente)

## Desafio GoLang Pós GoExpert - Rate Limiter

Este projeto faz parte da Pós GoExpert como desafio, nele precisaríamos criar um `rate limiter` que pudesse ser utilizado para controlar o tráfego de requisições para um serviço web. O `rate limiter` deveria ser capaz de limitar o número de requisições com base em dois critérios:

1 - **Endereço IP**: O rate limiter deve restringir o número de requisições recebidas de um único endereço de IP dentro de um intervalo de tempo de 1 segundo.

2 - **Token de Acesso**: O rate limiter deve também limitar as requisições baseadas em um TOKEN de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O token deve ser informado no header da requisição, no seguinte formato:
    a - API_KEY: <TOKEN>

> [!IMPORTANT]
> As configurações de limite do token de acesso devem sobrepor as do IP. Por exemplo, se o limite do IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações de limite do token.

### Requisitos do Desafio

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
  
> **Código HTTP**: 429 \
> **Mensagem**: you have reached the maximum number of requests or actions allowed within a certain time frame

- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

### Exemplos

1 - **Limitação por IP**: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.

2 - **Limitação por Token**: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.

Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

> [!TIP]
> Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

### Entrega

- O código-fonte completo da implementação.
- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- O servidor web deve responder na porta 8080.

## O que é o Rate Limiter

Este projeto implementa um **Rate Limiter** robusto para controle de requisições por IP e token, utilizando o Redis como mecanismo de persistência. A solução foi desenvolvida com foco em extensibilidade, eficiência e simplicidade de configuração.

## Funcionalidades

- Limitação de requisições por **IP** e/ou **TOKEN** de acesso único.
- Configuração granular de limites por segundo para cada token.
- Configuração do tempo de bloqueio por IP/TOKEN após exceder o limite.
- Middleware flexível para injeção no servidor web.
- Troca fácil do mecanismo de persistência através da estratégia **PersistenceProvider**.
- Todas as configurações realizadas via variáveis de ambiente ou arquivo **.env**.

## Requisitos

### Tecnologias

- **Golang**: Linguagem principal;
- **Redis**: Base de dados para persistência;
- **Docker**: Containerização da aplicação e redis.

### Configurações

As configurações do sistema podem ser realizadas via **variáveis de ambiente** ou no arquivo `.env` na raiz do projeto:

| Variável          | Descrição                              | Valor Padrão     |
|-------------------|----------------------------------------|------------------|
| REDIS_HOST        | Host do Redis                          | **Obrigatório**  |
| REDIS_PORT        | Porta do Redis                         | `6379`           |
| LIMITER_IP_LIMIT  | Limite global por IP (req/seg)         | `5`              |
| BLOCK_DURATION    | Duração do bloqueio em segundos        | `300`            |
| TOKEN_LIMITS      | Limites customizados por token (JSON)  | Exemplo abaixo   |

### Exemplo de `TOKEN_LIMITS`

```json
{"token1":10,"token2":15,"token3":20,"token4":30}
```

### Endpoints

- **GET /ping:** Endpoint de teste que responde com pong.

## Testes Automatizados

A aplicação possui **100%** de cobertura de testes nas seguintes áreas:

- Lógica principal do rate limiter.
- Middleware para controle de requisições.
- Integração com o Redis, , incluindo inicialização e validação do cliente.
- Cenários de erro, incluindo persistência e TTL.

## Executando o Projeto

O projeto foi desenvolvido utilizando Docker/Docker Compose de modo a facilitar a execução do mesmo, entõ para poder executar o projeto basta seguir os passos descritos abaixo.

### Passos

1. Clone o repositório:

```bash
git clone https://github.com/vs0uz4/rate-limit.git
cd rate-limit
```

2. Configure o `.env` na raiz com base no arquivo `.env.example`.

3. Suba o ambiente com Docker Compose:

```bash
docker-compose up --build
```

> [!TIP]
> Certifique-se que não tenha nada rodando em sua máquina nas portas `8080`, `8081` e `6379` que são as portas utilizadas pelo projeto e suas dependências.
>
> 8080 - Aplicação \
> 8081 - UI Redis Commander \
> 6379 - Redis Server

4. Faça requisições para o endpoint **GET /ping**.

## Exemplos de CURL

### Sem token

```bash
curl --request GET \
  --url <http://localhost:8080/ping>
```

### Com token

```bash
curl --request GET \
  --url <http://localhost:8080/ping> \
  --header 'API_KEY: token1'
```

### Com token inválido/inexistente

```bash
curl --request GET \
  --url <http://localhost:8080/ping> \
  --header 'API_KEY: invalid_token'
```

> [!IMPORTANT]
> No caso de requisições com **TOKENS** inválidos, ou seja, que não constem na ENV `TOKEN_LIMITS` a request será tratada com o limite de bloqueio padrão para IP, que no momento atual, está configurado para 5 req/s.
