# Arquitetura do TreinoZap

## Visão Geral

O TreinoZap é um monorepo com frontend em Next.js, backend em Go e PostgreSQL como banco principal.

O fluxo de negócio parte do painel do treinador e chega ao cliente pelo WhatsApp:

```text
Treinador no navegador
        |
        v
Frontend (Next.js :3000)
        |
        v
Backend API (Go :8080)
        |
        +--> PostgreSQL
        |
        \--> Provider de WhatsApp
             |- mock
             \- whatsmeow
```

Em produção local com Docker Compose, a API já é capaz de operar o `whatsmeow` diretamente. O `worker` existe como processo opcional para cenários específicos, mas não é necessário no fluxo padrão do projeto.

## Backend

### Stack

- Go `1.26`
- `chi` para roteamento HTTP
- `pgx` para acesso ao PostgreSQL
- `golang-migrate` para migrations
- `golang-jwt/jwt` para autenticação
- `bcrypt` para hash de senha
- `whatsmeow` para integração com WhatsApp

### Entrypoints

```text
/backend/cmd/api
  servidor HTTP principal

/backend/cmd/migrate
  runner de migrations

/backend/cmd/worker
  worker opcional de processamento WhatsApp
```

### Estrutura de pacotes

```text
/backend
  /internal
    /admin       endpoints administrativos
    /auth        geração e validação de JWT
    /automation  tratamento de comandos recebidos no WhatsApp
    /client      domínio de clientes
    /config      leitura de variáveis de ambiente
    /database    conexão com PostgreSQL
    /exercise    domínio de exercícios
    /http
      /middleware autenticação e proteção de rotas
      /response   envelope padrão de resposta JSON
      /routes     montagem das rotas da API
    /message     histórico de mensagens
    /trainer     domínio de treinadores
    /whatsapp    abstrações e providers
    /workout     domínio de treinos, seções e exercícios
```

### Padrão por domínio

Em geral, cada domínio segue este arranjo:

```text
model.go
dto.go
repository.go
service.go
handler.go
```

Responsabilidades:

- `model.go`: estruturas de domínio;
- `dto.go`: contratos de entrada e saída;
- `repository.go`: persistência com SQL via `pgx`;
- `service.go`: regras de negócio;
- `handler.go`: camada HTTP.

## Frontend

### Stack

- Next.js `16`
- React `19`
- TypeScript
- Tailwind CSS `4`

### Estrutura atual

O frontend usa App Router diretamente em `frontend/app`, sem a pasta `src`.

```text
/frontend
  /app
    /(dashboard)   área autenticada
    /login
    /register
  /components
    /layout
    /ui
  /lib
    api.ts
    auth.ts
  /types
```

### Autenticação no frontend

O token JWT é armazenado em `localStorage` com a chave `treinozap_token`.

Fluxo resumido:

1. login chama `POST /api/v1/auth/login`;
2. o token retornado é salvo no navegador;
3. requisições subsequentes enviam `Authorization: Bearer <token>`;
4. o layout autenticado valida a presença do token no cliente.

## Banco de Dados

O projeto usa PostgreSQL `16`.

Os dados principais estão distribuídos entre:

- `trainers`
- `clients`
- `exercises`
- `workouts`
- `workout_sections`
- `workout_exercises`
- `whatsapp_channels`
- `whatsapp_messages`
- `automation_rules`

Mais detalhes em [database.md](database.md).

## Integração com WhatsApp

O backend depende de duas abstrações:

```go
type Sender interface {
    SendText(ctx context.Context, phone string, message string) error
}

type ConnectionManager interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Status(ctx context.Context) (Status, error)
    QRCode(ctx context.Context) (string, error)
}
```

Isso permite alternar entre:

- `mock`, usado para desenvolvimento;
- `whatsmeow`, usado para conexão real.

No provider real:

- a sessão é armazenada no banco usado pelo projeto;
- o QR Code é gerado sob demanda;
- a API tenta reconectar sessões já pareadas;
- mensagens recebidas acionam a automação de comandos.

Mais detalhes em [whatsapp.md](whatsapp.md).

## Ciclo de Requisição

Fluxo típico da API autenticada:

```text
HTTP request
   |
   v
chi router
   |
   v
Auth middleware
   |
   v
Handler
   |
   v
Service
   |
   v
Repository
   |
   v
PostgreSQL
```

## Decisões de Projeto

- Monorepo para manter frontend, backend e documentação próximos.
- SQL explícito em repositórios para previsibilidade e controle.
- Um único número central de WhatsApp no MVP para simplificar operação.
- Isolamento por `trainer_id` na maior parte dos recursos de negócio.
- Envelope JSON consistente para facilitar consumo no frontend.
