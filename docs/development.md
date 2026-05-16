# Guia de Desenvolvimento

## Pré-requisitos

- Go `1.26+`
- Node.js `20+`
- Docker e Docker Compose
- Make

## Subir o Ambiente Completo

```bash
cp .env.example .env
make up
make logs
```

Serviços padrão:

- frontend em `http://localhost:3000`
- backend em `http://localhost:8080`
- PostgreSQL exposto em `localhost:5433`

## Desenvolvimento Local sem Docker

### 1. Banco de dados

Você pode usar o PostgreSQL do `docker compose`:

```bash
docker compose up -d postgres
```

Nesse caso, para rodar o backend localmente a URL precisa apontar para a porta publicada no host:

```env
DATABASE_URL=postgres://treinozap:treinozap@localhost:5433/treinozap?sslmode=disable
```

### 2. Backend

```bash
cd backend
go mod download
go run ./cmd/api/main.go
```

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

## Variáveis de Ambiente

O backend usa `godotenv.Load()`, então o carregamento automático considera o diretório em que o processo for iniciado.

Na prática:

- se você iniciar pela raiz com variáveis exportadas, o `.env` da raiz funciona bem;
- se você rodar `cd backend && go run ./cmd/api/main.go`, prefira ter também um `backend/.env` ou exportar as variáveis antes.

Para desenvolvimento local, a forma mais simples é manter um `.env` na raiz do projeto e usar os alvos do `Makefile`, ou duplicar o arquivo para `backend/.env` quando for rodar o backend manualmente.

Exemplo:

```env
APP_ENV=development
DATABASE_URL=postgres://treinozap:treinozap@localhost:5433/treinozap?sslmode=disable
HTTP_PORT=8080
JWT_SECRET=dev-secret
JWT_EXPIRES_IN=24h
WHATSAPP_PROVIDER=mock
WHATSAPP_ADMIN_ENABLED=true
ADMIN_EMAIL=admin@treinozap.com
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Migrations

Aplicar todas:

```bash
make migrate-up
```

Reverter a última:

```bash
make migrate-down
```

As migrations seguem o padrão:

```text
000001_create_trainers.up.sql
000001_create_trainers.down.sql
```

## Testes

### Backend

```bash
make test
```

ou

```bash
cd backend
go test ./...
```

### Frontend

Atualmente o fluxo mais útil de validação é o build:

```bash
cd frontend
npm run build
```

## Formatação

```bash
make fmt
```

Esse comando:

- roda `gofmt -w .` no backend;
- tenta aplicar `eslint --fix` no frontend.

## Comandos Úteis

```bash
make help
make up
make down
make logs
make backend
make frontend
make dev-backend
make dev-frontend
```

## WhatsApp no Desenvolvimento

### Provider `mock`

Use para desenvolver sem um número real:

```env
WHATSAPP_PROVIDER=mock
```

O envio é apenas registrado em log.

### Provider `whatsmeow`

Use quando quiser validar o fluxo real:

```env
WHATSAPP_PROVIDER=whatsmeow
```

Fluxo básico:

1. subir backend com acesso ao PostgreSQL;
2. acessar o painel com um usuário admin;
3. abrir `Configurações > WhatsApp`;
4. clicar em conectar;
5. escanear o QR Code.

## Convenções Práticas

- prefira manter endpoints e serviços organizados por domínio;
- preserve o envelope JSON padrão das respostas;
- normalize telefone e e-mail no backend, não só no frontend;
- ao mexer em docs, valide contra o código real antes de descrever comportamentos.

## Sugestão de Mensagens de Commit

```text
feat: adiciona CRUD de exercícios
fix: corrige ativação de treino anterior
docs: atualiza documentação da API e WhatsApp
```
