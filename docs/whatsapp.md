# WhatsApp

## Visão Geral

O TreinoZap opera com um único número central de WhatsApp no modelo atual.

Os treinadores não conectam contas próprias. A conexão é gerenciada por um usuário administrador por meio do painel web.

## Providers Disponíveis

### `mock`

Usado para desenvolvimento e testes de fluxo sem integração real.

Comportamento:

- `SendText` apenas registra a mensagem em log;
- `Connect` e `Disconnect` alteram um estado simulado;
- `QRCode` retorna um valor fictício.

### `whatsmeow`

Usado para integração real com WhatsApp Web.

Comportamentos implementados:

- conexão com store persistida no PostgreSQL;
- geração de QR Code sob demanda;
- reconexão automática após desconexão;
- reaproveitamento de sessão pareada quando já existir;
- recebimento de mensagens privadas com texto.

## Fluxo de Conexão

Quando o provider é `whatsmeow`, o fluxo normal é:

1. o admin acessa a tela de configurações do WhatsApp;
2. aciona o endpoint de conexão;
3. a API inicializa o fluxo de QR Code;
4. o admin escaneia o código com o número central;
5. a sessão fica persistida para reconexões futuras.

Se a sessão já existir salva, a API tenta reconectar sem exigir novo QR Code.

## Fluxo de Mensagem Recebida

```text
Cliente envia mensagem
        |
        v
Provider recebe evento
        |
        v
Backend extrai telefone e texto
        |
        v
Mensagem inbound é registrada
        |
        v
Cliente é buscado por telefone
        |
        +--> não encontrado
        |      responde com orientação de cadastro
        |
        +--> inactive
        |      responde que o acesso está inativo
        |
        +--> blocked
        |      não responde
        |
        \--> cliente válido
               comando é normalizado e roteado
```

## Comandos Reconhecidos

Os comandos atualmente tratados pelo fluxo de automação são:

| Comando | Resultado |
|---|---|
| `oi`, `ola`, `bom dia`, `boa tarde`, `boa noite`, `hello`, `hi` | Envia menu |
| `status` | Informa se há treino ativo |
| `treino`, `meu treino`, `quero meu treino`, `enviar treino` | Envia treino completo |
| `menu` | Envia opções disponíveis |
| `a`, `treino a` | Envia seção A |
| `b`, `treino b` | Envia seção B |
| `c`, `treino c` | Envia seção C |
| `ajuda` | Envia instruções |
| qualquer outro | Envia fallback orientando a usar `menu` |

## Formatação das Mensagens

O formatter atual usa texto puro com marcações compatíveis com WhatsApp.

Exemplo simplificado:

```text
Olá, João! Aqui está seu treino ativo:

*Plano Hipertrofia - Semana 1*

*Treino A - Peito e Tríceps*
1. Supino reto — 4x10 — descanso 60s
   Carga: moderada
   Técnica: controlar a descida
```

Também há suporte a:

- `load_note`
- `technique_note`
- `video_url`

## Status e Administração

Os endpoints administrativos expõem:

- status da conexão;
- QR Code atual;
- acionamento de conexão;
- desconexão manual.

Formato atual do status:

```json
{
  "connected": true,
  "phone": "5592999999999",
  "jid": "5592999999999@s.whatsapp.net",
  "last_connected": "2026-05-15T21:00:00Z"
}
```

## Decisões de Implementação

- mensagens de grupo são ignoradas;
- mensagens do próprio dispositivo são ignoradas;
- mensagens sem texto são ignoradas;
- o QR Code fica disponível apenas após `Connect`;
- quando a conexão já está autenticada, o QR deixa de ser necessário.

## Persistência

Há duas camadas de persistência relevantes:

- a store interna do `whatsmeow`, usada para manter a sessão autenticada;
- a tabela `whatsapp_messages`, usada para histórico funcional do produto.

Mensagens inbound e outbound são registradas com:

- direção (`inbound` ou `outbound`);
- telefone;
- comando identificado;
- status;
- relacionamento com cliente e treinador quando possível.

## Segurança

- apenas usuários com role `admin` operam os endpoints de conexão;
- não existe endpoint público para envio arbitrário de mensagens;
- o MVP não oferece broadcast;
- a identificação do cliente acontece pelo telefone normalizado.

## Limitações Atuais

- número único central para toda a plataforma;
- ausência de multi-sessão por treinador;
- automação baseada em comandos fixos do backend;
- sem mídia, áudio ou anexos no fluxo principal atual.
