# Banco de Dados

## Visão Geral

O TreinoZap usa PostgreSQL `16` como banco principal.

As migrations ficam em `backend/migrations/` e cobrem a estrutura base do domínio, o histórico de mensagens e as tabelas relacionadas ao WhatsApp.

## Tabelas Principais

### `trainers`

Treinadores cadastrados na plataforma.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `name` | `TEXT` | Nome do treinador |
| `email` | `TEXT` | E-mail único |
| `password_hash` | `TEXT` | Senha com bcrypt |
| `phone` | `TEXT` | Telefone |
| `role` | `TEXT` | `trainer` ou `admin` |
| `status` | `TEXT` | `active` ou `inactive` |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

### `clients`

Clientes vinculados a treinadores.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `trainer_id` | `UUID` | Referência para `trainers` |
| `name` | `TEXT` | Nome do cliente |
| `phone` | `TEXT` | Telefone normalizado |
| `status` | `TEXT` | `active`, `inactive` ou `blocked` |
| `goal` | `TEXT` | Objetivo do aluno |
| `notes` | `TEXT` | Observações livres |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

Regra importante:

- `phone` é `UNIQUE` no banco inteiro, o que combina com o modelo atual de número central único no WhatsApp.

### `exercises`

Biblioteca de exercícios criada por treinador.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `trainer_id` | `UUID` | Referência para `trainers` |
| `name` | `TEXT` | Nome do exercício |
| `muscle_group` | `TEXT` | Grupo muscular |
| `equipment` | `TEXT` | Equipamento |
| `video_url` | `TEXT` | Link opcional de demonstração |
| `notes` | `TEXT` | Observações |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

### `workouts`

Treinos atribuídos aos clientes.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `trainer_id` | `UUID` | Referência para `trainers` |
| `client_id` | `UUID` | Referência para `clients` |
| `name` | `TEXT` | Nome do treino |
| `status` | `TEXT` | `draft`, `active` ou `archived` no fluxo atual |
| `starts_at` | `DATE` | Data inicial opcional |
| `ends_at` | `DATE` | Data final opcional |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

Regra operacional:

- a ativação de um treino arquiva o treino ativo anterior do mesmo cliente.

### `workout_sections`

Seções de um treino, como `Treino A`, `Treino B` e `Treino C`.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `workout_id` | `UUID` | Referência para `workouts` |
| `name` | `TEXT` | Nome da seção |
| `description` | `TEXT` | Descrição opcional |
| `order_index` | `INT` | Ordem de exibição |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

### `workout_exercises`

Itens de exercício dentro de cada seção.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `section_id` | `UUID` | Referência para `workout_sections` |
| `exercise_id` | `UUID` | Referência opcional para `exercises` |
| `exercise_name` | `TEXT` | Nome congelado no momento da montagem |
| `sets` | `TEXT` | Quantidade de séries |
| `reps` | `TEXT` | Quantidade de repetições |
| `rest_seconds` | `INT` | Descanso em segundos |
| `load_note` | `TEXT` | Observação de carga |
| `technique_note` | `TEXT` | Observação técnica |
| `video_url` | `TEXT` | Link opcional |
| `order_index` | `INT` | Ordem dentro da seção |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

### `whatsapp_channels`

Representa o canal central de WhatsApp.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `name` | `TEXT` | Nome do canal |
| `phone` | `TEXT` | Número associado |
| `jid` | `TEXT` | JID do WhatsApp |
| `status` | `TEXT` | `disconnected`, `connecting` ou `connected` |
| `is_default` | `BOOLEAN` | Identifica o canal principal |
| `last_connected_at` | `TIMESTAMPTZ` | Última conexão conhecida |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

Observação:

- a operação atual do `whatsmeow` usa a store própria da biblioteca no PostgreSQL, além desta modelagem de domínio.

### `whatsapp_messages`

Histórico de mensagens recebidas e enviadas.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `channel_id` | `UUID` | Referência opcional para `whatsapp_channels` |
| `trainer_id` | `UUID` | Referência opcional para `trainers` |
| `client_id` | `UUID` | Referência opcional para `clients` |
| `direction` | `TEXT` | `inbound` ou `outbound` |
| `phone` | `TEXT` | Número envolvido na conversa |
| `message` | `TEXT` | Conteúdo textual |
| `command` | `TEXT` | Comando detectado |
| `status` | `TEXT` | Ex.: `received` ou `sent` |
| `provider_message_id` | `TEXT` | ID externo, quando aplicável |
| `created_at` | `TIMESTAMPTZ` | Momento do registro |

### `automation_rules`

Regras de automação por treinador.

| Coluna | Tipo | Descrição |
|---|---|---|
| `id` | `UUID` | Chave primária |
| `trainer_id` | `UUID` | Referência para `trainers` |
| `keyword` | `TEXT` | Palavra-chave |
| `action` | `TEXT` | Ação configurada |
| `is_active` | `BOOLEAN` | Indica se a regra está ativa |
| `created_at` | `TIMESTAMPTZ` | Data de criação |
| `updated_at` | `TIMESTAMPTZ` | Data de atualização |

Observação:

- a tabela já existe e está pronta para evolução, mas o fluxo principal atual de mensagens é dirigido pelo handler de automação embutido no backend.

## Relações Relevantes

```text
trainers 1 --- N clients
trainers 1 --- N exercises
trainers 1 --- N workouts
clients  1 --- N workouts
workouts 1 --- N workout_sections
workout_sections 1 --- N workout_exercises
trainers 1 --- N whatsapp_messages
clients  1 --- N whatsapp_messages
```

## Índices e Regras Operacionais

- `clients.phone` é único.
- Há índices por `trainer_id` nas tabelas de domínio principais.
- `whatsapp_messages` possui índices por `trainer_id`, `client_id`, `phone` e `direction`.
- `workout_sections` e `workout_exercises` usam `ON DELETE CASCADE` nas relações internas.

## Estratégia de Persistência

- O backend usa SQL explícito via `pgx`.
- Criação e atualização de treinos com seções e exercícios são transacionais.
- As listagens de clientes, exercícios e mensagens suportam paginação.
- A store do `whatsmeow` também usa o PostgreSQL da aplicação.
