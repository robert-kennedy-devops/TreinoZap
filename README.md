# TreinoZap

Plataforma web para personal trainers gerenciarem alunos, montarem treinos e entregarem os planos diretamente pelo WhatsApp.

O projeto foi organizado como um monorepo com backend em Go, frontend em Next.js e PostgreSQL, com suporte a envio simulado (`mock`) ou integração real com WhatsApp via `whatsmeow`.

## Visão Geral

No fluxo principal, o treinador acessa o painel web para:

- cadastrar clientes;
- manter uma biblioteca de exercícios;
- montar treinos com seções como A, B e C;
- ativar um treino para um cliente;
- enviar o treino manualmente ou deixar o cliente solicitar pelo WhatsApp.

O cliente não precisa de login nem aplicativo próprio: a interação acontece por mensagens enviadas ao número central da plataforma.

## Principais Capacidades

- autenticação de treinadores com JWT;
- cadastro e gestão de clientes;
- CRUD de exercícios;
- criação, edição, ativação e arquivamento de treinos;
- envio manual de treino por WhatsApp;
- automação para respostas por comando no WhatsApp;
- histórico de mensagens por cliente;
- área administrativa para conexão e monitoramento do número central;
- execução completa com Docker Compose.

## Arquitetura

```text
Navegador
   |
   +--> Frontend (Next.js / React)
           |
           +--> Backend API (Go + chi)
                   |
                   +--> PostgreSQL
                   |
                   +--> WhatsApp Provider
                         |- mock
                         \- whatsmeow
```

### Backend

- Go `1.26`
- `chi` para roteamento HTTP
- `pgx` para acesso ao PostgreSQL
- `golang-migrate` para migrations
- `golang-jwt/jwt` para autenticação
- `whatsmeow` para integração real com WhatsApp

### Frontend

- Next.js `16`
- React `19`
- TypeScript
- Tailwind CSS `4`

### Infra

- PostgreSQL `16`
- Docker e Docker Compose

## Estrutura do Repositório

```text
.
├── backend
│   ├── cmd
│   │   ├── api       # servidor HTTP principal
│   │   ├── migrate   # runner de migrations
│   │   └── worker    # worker opcional para WhatsApp
│   ├── internal      # domínios, serviços, handlers e integrações
│   └── migrations    # versionamento do banco
├── frontend          # painel web em Next.js
├── docs              # documentação complementar
├── docker-compose.yml
└── Makefile
```

## Requisitos

- Docker e Docker Compose
- Make

Para desenvolvimento sem Docker:

- Go `1.26+`
- Node.js `20+`
- PostgreSQL

## Como Executar com Docker

### 1. Configurar variáveis

```bash
cp .env.example .env
```

Edite o `.env` se quiser alterar portas, segredo JWT ou provider do WhatsApp.

### 2. Subir o ambiente

```bash
make up
```

### 3. Acompanhar os logs

```bash
make logs
```

### Endereços padrão

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Health check: `http://localhost:8080/health`
- PostgreSQL exposto no host: `localhost:5433`

## Variáveis de Ambiente

As principais variáveis usadas pelo projeto são:

| Variável | Exemplo | Descrição |
|---|---|---|
| `APP_ENV` | `development` | Ambiente da aplicação |
| `DATABASE_URL` | `postgres://treinozap:treinozap@postgres:5432/treinozap?sslmode=disable` | String de conexão do backend |
| `HTTP_PORT` | `8080` | Porta HTTP do backend |
| `JWT_SECRET` | `change-me` | Chave usada para assinar JWT |
| `JWT_EXPIRES_IN` | `24h` | Duração do token |
| `WHATSAPP_PROVIDER` | `mock` | Provider do WhatsApp: `mock` ou `whatsmeow` |
| `WHATSAPP_ADMIN_ENABLED` | `true` | Habilita recursos administrativos do WhatsApp |
| `ADMIN_EMAIL` | `admin@treinozap.com` | E-mail tratado como administrador |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080/api/v1` | URL base consumida pelo frontend |

## Fluxo de Uso

1. Acesse `http://localhost:3000/register` e crie um treinador.
2. Faça login no painel.
3. Cadastre clientes com nome e telefone.
4. Cadastre exercícios.
5. Monte um treino para um cliente.
6. Ative o treino.
7. Envie o treino manualmente pelo painel ou aguarde o cliente solicitar via WhatsApp.

## WhatsApp

O projeto trabalha com um único número central de WhatsApp administrado pela plataforma.

### Modo `mock`

Use em desenvolvimento quando não quiser integrar com um número real:

```env
WHATSAPP_PROVIDER=mock
```

Nesse modo, os envios são simulados e registrados em log.

### Modo `whatsmeow`

Use para conectar um número real:

```env
WHATSAPP_PROVIDER=whatsmeow
```

Fluxo esperado:

1. entrar com um usuário administrador;
2. acessar `Configurações > WhatsApp`;
3. iniciar a conexão;
4. escanear o QR Code com o número central;
5. manter a sessão persistida para reconexão automática.

O backend também trata casos operacionais importantes, como:

- reconexão automática após queda;
- ignorar mensagens de grupo;
- ignorar mensagens sem texto;
- persistência da sessão em volume Docker.

### Comandos aceitos pelo cliente

| Comando | Ação |
|---|---|
| `status` | Informa se existe treino ativo |
| `treino` | Retorna o treino ativo completo |
| `menu` | Mostra opções disponíveis |
| `a` ou `treino a` | Retorna a seção A |
| `b` ou `treino b` | Retorna a seção B |
| `c` ou `treino c` | Retorna a seção C |
| `ajuda` | Mostra instruções de uso |

## Endpoints Principais

### Públicos

```text
GET  /health
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Autenticados

```text
GET    /api/v1/me

GET    /api/v1/clients
POST   /api/v1/clients
GET    /api/v1/clients/{id}
PUT    /api/v1/clients/{id}
DELETE /api/v1/clients/{id}

GET    /api/v1/exercises
POST   /api/v1/exercises
GET    /api/v1/exercises/{id}
PUT    /api/v1/exercises/{id}
DELETE /api/v1/exercises/{id}

GET    /api/v1/clients/{clientId}/workouts
POST   /api/v1/clients/{clientId}/workouts
GET    /api/v1/workouts/{id}
PUT    /api/v1/workouts/{id}
DELETE /api/v1/workouts/{id}
POST   /api/v1/workouts/{id}/activate
POST   /api/v1/workouts/{id}/send-whatsapp

GET    /api/v1/messages
GET    /api/v1/clients/{clientId}/messages
```

### Administrativos

```text
GET  /api/v1/admin/whatsapp/status
GET  /api/v1/admin/whatsapp/qr
POST /api/v1/admin/whatsapp/connect
POST /api/v1/admin/whatsapp/disconnect
GET  /api/v1/admin/trainers
GET  /api/v1/admin/clients
```

## Comandos Úteis

```bash
make help
make up
make down
make logs
make migrate-up
make migrate-down
make test
make fmt
make dev-backend
make dev-frontend
```

## Desenvolvimento Local

### Backend

```bash
cd backend
go mod download
go run ./cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Observações

- em ambiente local sem Docker, ajuste o `DATABASE_URL` para sua instância PostgreSQL;
- o backend carrega variáveis via `.env`;
- o frontend consome a API definida em `NEXT_PUBLIC_API_URL`.

## Documentação Complementar

- [Arquitetura](docs/architecture.md)
- [API](docs/api.md)
- [Banco de Dados](docs/database.md)
- [Desenvolvimento](docs/development.md)
- [Segurança](docs/security.md)
- [WhatsApp](docs/whatsapp.md)

## Limitações Atuais

- um único número central de WhatsApp para toda a plataforma;
- sem aplicativo ou portal dedicado para o aluno;
- sem módulo financeiro;
- sem suporte a múltiplas sessões por treinador;
- foco em operação MVP.

## Próximos Passos Sugeridos

- multi-sessão de WhatsApp por treinador;
- portal do aluno;
- evolução das regras de automação;
- relatórios operacionais;
- integrações de pagamento;
- recursos de avaliação física e acompanhamento.
