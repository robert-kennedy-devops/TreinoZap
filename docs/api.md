# API Reference

Base URL padrão em desenvolvimento: `http://localhost:8080`

## Convenção de Resposta

O backend responde usando envelopes JSON.

### Sucesso

```json
{
  "data": {}
}
```

### Erro

```json
{
  "error": {
    "message": "mensagem de erro",
    "code": "ERROR_CODE"
  }
}
```

Observação importante:

- endpoints simples retornam o recurso diretamente dentro de `data`;
- endpoints paginados retornam um objeto dentro de `data`, com lista e metadados de paginação.

Exemplo de listagem paginada:

```json
{
  "data": {
    "data": [],
    "total": 0,
    "page": 1,
    "page_size": 20,
    "total_pages": 0
  }
}
```

## Health Check

### `GET /health`

Verifica se a API está no ar.

Resposta:

```json
{
  "data": {
    "status": "ok"
  }
}
```

## Autenticação

### `POST /api/v1/auth/register`

Cadastra um novo treinador.

Body:

```json
{
  "name": "Carlos Personal",
  "email": "carlos@email.com",
  "password": "123456",
  "phone": "5592999999999"
}
```

Regras atuais:

- `name`, `email` e `password` são obrigatórios;
- a senha deve ter pelo menos 6 caracteres;
- `email` é normalizado para lowercase;
- se o e-mail for igual a `ADMIN_EMAIL`, o usuário é criado com role `admin`.

Resposta:

```json
{
  "data": {
    "id": "uuid",
    "name": "Carlos Personal",
    "email": "carlos@email.com",
    "phone": "5592999999999",
    "role": "trainer",
    "status": "active",
    "created_at": "2026-05-15T00:00:00Z",
    "updated_at": "2026-05-15T00:00:00Z"
  }
}
```

### `POST /api/v1/auth/login`

Autentica um treinador e retorna JWT.

Body:

```json
{
  "email": "carlos@email.com",
  "password": "123456"
}
```

Resposta:

```json
{
  "data": {
    "token": "eyJ...",
    "trainer": {
      "id": "uuid",
      "name": "Carlos Personal",
      "email": "carlos@email.com",
      "phone": "5592999999999",
      "role": "trainer",
      "status": "active"
    }
  }
}
```

### `GET /api/v1/me`

Retorna os dados do treinador autenticado.

Header obrigatório:

```text
Authorization: Bearer <token>
```

## Clientes

Todos os endpoints abaixo exigem autenticação.

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/clients` | Lista clientes com paginação |
| `POST` | `/api/v1/clients` | Cria cliente |
| `GET` | `/api/v1/clients/{id}` | Busca cliente por ID |
| `PUT` | `/api/v1/clients/{id}` | Atualiza cliente |
| `DELETE` | `/api/v1/clients/{id}` | Inativa cliente |

### Query params de `GET /api/v1/clients`

- `page`
- `page_size`
- `search`

### Exemplo de criação

```json
{
  "name": "João Silva",
  "phone": "(92) 99999-9999",
  "goal": "Hipertrofia",
  "notes": "Treina 4x por semana"
}
```

Comportamentos relevantes:

- o telefone é normalizado para apenas dígitos;
- `name` e `phone` são obrigatórios;
- o `DELETE` atual faz inativação lógica, alterando `status` para `inactive`.

## Exercícios

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/exercises` | Lista exercícios com paginação |
| `POST` | `/api/v1/exercises` | Cria exercício |
| `GET` | `/api/v1/exercises/{id}` | Busca exercício |
| `PUT` | `/api/v1/exercises/{id}` | Atualiza exercício |
| `DELETE` | `/api/v1/exercises/{id}` | Remove exercício |

### Query params de `GET /api/v1/exercises`

- `page`
- `page_size`
- `search`

### Exemplo de payload

```json
{
  "name": "Supino reto",
  "muscle_group": "Peito",
  "equipment": "Barra",
  "video_url": "https://example.com/video",
  "notes": "Priorizar controle na descida"
}
```

## Treinos

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/clients/{clientId}/workouts` | Lista treinos de um cliente |
| `POST` | `/api/v1/clients/{clientId}/workouts` | Cria treino completo |
| `GET` | `/api/v1/workouts/{id}` | Busca treino com seções e exercícios |
| `PUT` | `/api/v1/workouts/{id}` | Atualiza treino |
| `DELETE` | `/api/v1/workouts/{id}` | Arquiva treino |
| `POST` | `/api/v1/workouts/{id}/activate` | Ativa treino e arquiva o ativo anterior |
| `POST` | `/api/v1/workouts/{id}/send-whatsapp` | Envia treino manualmente por WhatsApp |

### Estrutura resumida do treino

```json
{
  "name": "Plano Hipertrofia - Semana 1",
  "status": "draft",
  "starts_at": "2026-05-15",
  "ends_at": "2026-05-22",
  "sections": [
    {
      "name": "Treino A - Peito e Tríceps",
      "description": "Foco em volume",
      "order_index": 1,
      "exercises": [
        {
          "exercise_id": "uuid-opcional",
          "exercise_name": "Supino reto",
          "sets": "4",
          "reps": "10",
          "rest_seconds": 60,
          "load_note": "Carga moderada",
          "technique_note": "Controlar a descida",
          "video_url": "https://example.com/video",
          "order_index": 1
        }
      ]
    }
  ]
}
```

Observações:

- a API persiste o treino e suas seções em transação;
- `GET /api/v1/workouts/{id}` retorna seções e exercícios embutidos;
- `POST /activate` garante um único treino `active` por cliente.

## Mensagens

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/messages` | Lista histórico geral do treinador |
| `GET` | `/api/v1/clients/{clientId}/messages` | Lista histórico de um cliente |

### Query params de `GET /api/v1/messages`

- `page`
- `page_size`
- `direction`

Valores típicos de `direction`:

- `inbound`
- `outbound`

## Administração

Todos os endpoints abaixo exigem role `admin`.

Quando `WHATSAPP_ADMIN_ENABLED=false`, os endpoints administrativos de WhatsApp retornam `403`.

### WhatsApp

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/admin/whatsapp/status` | Retorna status de conexão |
| `GET` | `/api/v1/admin/whatsapp/qr` | Retorna o QR Code atual |
| `POST` | `/api/v1/admin/whatsapp/connect` | Inicia conexão ou reconexão |
| `POST` | `/api/v1/admin/whatsapp/disconnect` | Desconecta a sessão atual |

Formato atual de status:

```json
{
  "data": {
    "connected": true,
    "phone": "5592999999999",
    "jid": "5592999999999@s.whatsapp.net",
    "last_connected": "2026-05-15T21:00:00Z"
  }
}
```

### Administração global

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/admin/trainers` | Lista treinadores |
| `GET` | `/api/v1/admin/clients` | Lista clientes globalmente |

`GET /api/v1/admin/clients` aceita `search`.

## Códigos de Erro Comuns

| HTTP | Code | Quando aparece |
|---|---|---|
| `400` | `BAD_REQUEST` | Body inválido ou validação básica |
| `401` | `UNAUTHORIZED` | Token ausente ou inválido |
| `401` | `INVALID_CREDENTIALS` | Login com credenciais inválidas |
| `403` | `FORBIDDEN` | Usuário sem role adequada |
| `404` | `NOT_FOUND` | Recurso inexistente |
| `409` | `CONFLICT` | E-mail ou telefone em uso, conforme caso |
| `500` | `INTERNAL_ERROR` | Erro interno inesperado |
