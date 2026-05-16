# Segurança

## Isolamento por `trainer_id`

Os recursos principais do painel são escopados por treinador autenticado.

Isso vale especialmente para:

- `clients`
- `exercises`
- `workouts`
- `whatsapp_messages`

Na prática, o backend extrai o `trainer_id` do JWT e o propaga para services e repositories, reduzindo o risco de acesso cruzado por manipulação de URL.

## JWT

- tokens são assinados com `JWT_SECRET`;
- o token contém `trainer_id` e `role`;
- a expiração padrão é controlada por `JWT_EXPIRES_IN`;
- endpoints autenticados exigem `Authorization: Bearer <token>`.

No frontend atual, o JWT é armazenado em `localStorage`.

## Roles e Administração

- endpoints em `/api/v1/admin/*` exigem role `admin`;
- a role é atribuída no cadastro quando o e-mail coincide com `ADMIN_EMAIL`;
- os endpoints administrativos de WhatsApp também respeitam `WHATSAPP_ADMIN_ENABLED`;
- treinadores comuns não podem operar endpoints administrativos do WhatsApp.

## Senhas

- senhas são armazenadas apenas como hash `bcrypt`;
- `password_hash` não é retornado nas respostas da API;
- falhas de login retornam mensagem genérica de credenciais inválidas.

## WhatsApp

Medidas já implementadas no fluxo atual:

- apenas um número central é utilizado no MVP;
- mensagens de grupo são ignoradas;
- mensagens enviadas pelo próprio dispositivo são ignoradas;
- mensagens sem texto são ignoradas;
- telefones recebidos são normalizados antes da busca;
- clientes `inactive` recebem mensagem de acesso inativo;
- clientes `blocked` não recebem resposta automática.

## CORS

As origens atualmente aceitas no backend são:

- `http://localhost:3000`
- `http://frontend:3000`

Isso cobre desenvolvimento local e o ambiente padrão do Docker Compose.

## Logs

Boas práticas observadas na implementação atual:

- senhas não são logadas;
- tokens não são logados;
- respostas de erro ao cliente não incluem stack trace;
- falhas internas retornam envelope de erro padronizado.

## Normalização de Dados

- e-mails são normalizados para lowercase no cadastro e login;
- telefones são reduzidos a dígitos no backend;
- comandos do WhatsApp são normalizados com lowercase, trim e remoção de acentos.

## SQL Injection

O acesso ao banco usa SQL parametrizado com placeholders posicionais, como `$1`, `$2` e assim por diante.

Não há construção de query concatenando entrada bruta do usuário para filtros principais.

## Riscos e Limitações Atuais

- o JWT em `localStorage` é simples e funcional, mas expõe trade-offs de segurança conhecidos no frontend;
- não há suporte atual a múltiplas sessões de WhatsApp por treinador;
- a proteção administrativa do WhatsApp depende de role `admin` e da flag `WHATSAPP_ADMIN_ENABLED`.
