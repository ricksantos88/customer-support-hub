# Plano Detalhado de Desenvolvimento e Documentação
## Customer Support Hub - WhatsApp API

**Data:** 15 de maio de 2026  
**Status:** Em planejamento  
**Versão:** 1.0

---

## 1. ANÁLISE DO PROJETO

### 1.1 Visão Geral
- **Objetivo:** Backend API em Go para integração com WhatsApp e suporte a frontend de customer support
- **Foco:** Receber/enviar mensagens, persistir dados, atualizações em tempo real via WebSocket
- **Escopo MVP:** Funcionalidades core sem automação ou IA

### 1.2 Stack Técnico
- **Backend:** Go 1.24+, Fiber
- **Banco:** PostgreSQL (persistência) + Redis (cache/sessões)
- **Infra:** Docker, Docker Compose
- **Integração:** WhatsApp Cloud API (Meta)
- **Comunicação:** WebSocket (tempo real)

### 1.3 Requisitos Não-Funcionais
- Latência < 300ms (operações internas)
- Escalabilidade horizontal
- Alta disponibilidade
- Gerenciamento seguro de secrets
- Logs estruturados

### 1.4 Arquitetura
- **Estilo:** Hexagonal + DDD lite com abstração de providers
- **Padrão:** Separação clara de camadas (Domain, Application, Infrastructure, Interfaces)
- **Provider:** Interface `MessagingProvider` para abstrair WhatsApp

---

## 2. FASES DE DESENVOLVIMENTO

### FASE 1: SETUP E INFRAESTRUTURA (Sprint 1)
**Duração:** 1 semana | **Entregáveis:** Projeto bootstrapped, infra pronta

#### 2.1.1 Tasks

**Task 1.1: Inicializar Projeto Go**
- [ ] Criar estrutura base do projeto
- [ ] Configurar `go.mod` com dependências: Fiber, GORM, Redis, JWT
- [ ] Definir folder structure conforme arquitetura
- [ ] Criar `.gitignore` e `Makefile` para comandos comuns

**Task 1.2: Configuração de Ambiente**
- [ ] Criar arquivo `.env.example` com todas as variáveis necessárias
- [ ] Implementar config management (viper ou similar)
- [ ] Setup de logging estruturado (slog ou zap)
- [ ] Configurar variáveis para dev/test/prod

**Task 1.3: Docker e Docker Compose**
- [ ] Criar Dockerfile para aplicação Go
- [ ] Criar docker-compose.yml com: API, PostgreSQL, Redis
- [ ] Configurar volumes para dados persistentes
- [ ] Testar local com docker-compose up

**Task 1.4: CI/CD Básico**
- [ ] Configurar GitHub Actions (lint, test, build)
- [ ] Setup de secrets no CI
- [ ] Build de imagem Docker no pipeline

---

### FASE 2: BANCO DE DADOS E MODELOS (Sprint 1-2)
**Duração:** 1 semana | **Entregáveis:** Schema definido, modelos Go prontos

#### 2.2.1 Tasks

**Task 2.1: Design do Schema PostgreSQL**
- [ ] Definir tabela `contacts` com campos: id (UUID), phone, name, created_at, updated_at, deleted_at
- [ ] Definir tabela `conversations` com: id, contact_id, status (open/pending/closed), assigned_agent_id, created_at, updated_at
- [ ] Definir tabela `messages` com: id, conversation_id, content, direction (inbound/outbound), sender_id, created_at
- [ ] Definir tabela `agents` com: id, name, email, jwt_secret, created_at, last_active
- [ ] Criar índices em: phone, conversation_id, assigned_agent_id, created_at

**Task 2.2: Migrations**
- [ ] Usar biblioteca de migrations (golang-migrate ou similar)
- [ ] Criar migration para v1 com todas as tabelas
- [ ] Documentar comando para rodar migrations
- [ ] Testar rollback

**Task 2.3: Modelos GORM**
- [ ] Criar structs Go para: Contact, Conversation, Message, Agent
- [ ] Implementar hooks GORM (BeforeCreate, BeforeSave)
- [ ] Definir relacionamentos (foreign keys, belongs to, has many)
- [ ] Adicionar validações básicas

**Task 2.4: Repositório Pattern**
- [ ] Criar interfaces de repositório para cada entidade
- [ ] Implementar PostgreSQL repositories
- [ ] Implementar mock repositories para testes
- [ ] Testar operações CRUD

---

### FASE 3: AUTENTICAÇÃO E SEGURANÇA (Sprint 2)
**Duração:** 1 semana | **Entregáveis:** JWT setup, middleware de auth, secrets management

#### 2.3.1 Tasks

**Task 3.1: JWT Implementation**
- [ ] Gerar JWT_SECRET securo
- [ ] Implementar geração de tokens para agents
- [ ] Implementar validação de tokens
- [ ] Adicionar refresh token logic
- [ ] Definir TTL para tokens (ex: 24h)

**Task 3.2: Middleware de Autenticação**
- [ ] Criar middleware para validar Bearer tokens
- [ ] Implementar middleware para extrair agent_id do token
- [ ] Criar middleware de autorização baseada em roles (se necessário)
- [ ] Testar com requisições mock

**Task 3.3: Secrets Management**
- [ ] Configurar suporte a .env (dotenv)
- [ ] Documentar processo de geração de secrets
- [ ] Implementar validação de secrets obrigatórios
- [ ] Setup de secrets no Docker Compose

**Task 3.4: Rate Limiting e Proteção**
- [ ] Implementar rate limiting por IP
- [ ] Implementar rate limiting por agent (JWT claims)
- [ ] Configurar CORS apropriadamente
- [ ] Adicionar validação de input

---

### FASE 4: INTEGRAÇÃO COM WHATSAPP (Sprint 2-3)
**Duração:** 1.5 semanas | **Entregáveis:** Provider implementado, webhooks funcionando

#### 2.4.1 Tasks

**Task 4.1: Setup Meta Cloud API**
- [ ] Criar/configurar app no Meta Developers
- [ ] Gerar access token
- [ ] Capturar phone_number_id e business_account_id
- [ ] Documentar o processo (incluir no 04-setup-meta-cloud.md)
- [ ] Testar autenticação

**Task 4.2: Provider Interface**
- [ ] Implementar interface `MessagingProvider`
- [ ] Criar struct `MetaProvider` com client HTTP
- [ ] Implementar método `SendText()`
- [ ] Implementar método `HandleWebhook()`
- [ ] Adicionar mock provider para testes

**Task 4.3: Webhook Receiver**
- [ ] Criar handler HTTP POST `/webhooks/whatsapp`
- [ ] Implementar validação de webhook (signature verification)
- [ ] Criar handler HTTP GET `/webhooks/whatsapp` (challenge)
- [ ] Parsear payload do WhatsApp
- [ ] Testar com Postman/curl

**Task 4.4: Webhook Processor**
- [ ] Criar use case `ReceiveMessage`
- [ ] Parsear tipo de mensagem (texto, media, etc)
- [ ] Extrair sender, conteúdo, timestamp
- [ ] Persistir em `messages` table
- [ ] Testar validação de payload

---

### FASE 5: MENSAGENS - RECEBIMENTO (Sprint 3)
**Duração:** 1 semana | **Entregáveis:** Recebimento e persistência de mensagens

#### 2.5.1 Tasks

**Task 5.1: Use Case - Receive Message**
- [ ] Criar `ReceiveMessage` use case na application layer
- [ ] Implementar lógica: extrair sender, conteúdo, timestamp
- [ ] Se novo contact: criar em `contacts`
- [ ] Se nova conversa: criar em `conversations` (status=open)
- [ ] Persistir mensagem em `messages` (direction=inbound)

**Task 5.2: Validação de Webhook**
- [ ] Implementar validação de signature (HMAC)
- [ ] Implementar deduplicação (evitar reprocessamento)
- [ ] Adicionar logging estruturado
- [ ] Testar casos: válido, inválido, duplicado

**Task 5.3: Testes**
- [ ] Unit tests para ReceiveMessage use case
- [ ] Integration tests com banco fake
- [ ] E2E tests simulando webhook do WhatsApp
- [ ] Coverage > 80%

**Task 5.4: Documentação de API**
- [ ] Documentar webhook em Swagger/OpenAPI
- [ ] Documentar formato de payload esperado
- [ ] Documentar códigos de erro
- [ ] Criar exemplos de requisição/resposta

---

### FASE 6: MENSAGENS - ENVIO (Sprint 3)
**Duração:** 1 semana | **Entregáveis:** Endpoint de envio, integração com Meta

#### 2.6.1 Tasks

**Task 6.1: Endpoint POST /messages/send**
- [ ] Criar handler HTTP POST `/messages/send`
- [ ] Validar authentication (JWT middleware)
- [ ] Validar payload (conversation_id, text)
- [ ] Implementar autorização (agent pode mandar em sua conversa?)

**Task 6.2: Use Case - Send Message**
- [ ] Criar `SendMessage` use case
- [ ] Validar que conversa existe e está aberta
- [ ] Chamar `MessagingProvider.SendText()`
- [ ] Persistir em `messages` (direction=outbound)
- [ ] Retornar message_id e status

**Task 6.3: Error Handling**
- [ ] Tratar erros da Meta API
- [ ] Retry logic (exponential backoff)
- [ ] Fallback em caso de falha
- [ ] Logging detalhado de erros

**Task 6.4: Testes**
- [ ] Unit tests com mock provider
- [ ] Integration tests com PostgreSQL fake
- [ ] E2E com Docker Compose
- [ ] Coverage > 80%

---

### FASE 7: WEBSOCKET E REAL-TIME (Sprint 4)
**Duração:** 1.5 semanas | **Entregáveis:** WebSocket manager, broadcast de eventos

#### 2.7.1 Tasks

**Task 7.1: WebSocket Manager**
- [ ] Criar `WebSocketManager` na infrastructure layer
- [ ] Implementar pool de conexões por conversation_id
- [ ] Implementar broadcast para todos os agents de uma conversa
- [ ] Implementar graceful disconnect

**Task 7.2: Eventos de Mensagem**
- [ ] Definir estrutura de evento JSON (type, data, timestamp)
- [ ] Quando mensagem inbound: enviar via WS aos agents
- [ ] Quando mensagem outbound enviada: confirmar ao agent
- [ ] Testar com múltiplas conexões simultâneas

**Task 7.3: Endpoints WebSocket**
- [ ] Criar handler WS `/ws/conversations/:id`
- [ ] Validar JWT para upgrade
- [ ] Registrar conexão no manager
- [ ] Enviar histórico inicial de mensagens

**Task 7.4: Testes**
- [ ] Unit tests para WebSocketManager
- [ ] Testes de múltiplas conexões
- [ ] Testes de disconnect/reconnect
- [ ] Testes de broadcast

---

### FASE 8: ENDPOINTS COMPLEMENTARES (Sprint 4)
**Duração:** 1 semana | **Entregáveis:** CRUD de conversas, listagem, health

#### 2.8.1 Tasks

**Task 8.1: GET /conversations**
- [ ] Listar conversas abertas (filtrar por status)
- [ ] Paginação (limit, offset)
- [ ] Ordenar por created_at descending
- [ ] Retornar: id, contact.name, status, last_message, assigned_agent

**Task 8.2: GET /conversations/:id**
- [ ] Retornar detalhes da conversa
- [ ] Validar autorização
- [ ] Incluir contact info
- [ ] Incluir agent info

**Task 8.3: GET /conversations/:id/messages**
- [ ] Retornar histórico de mensagens
- [ ] Paginação
- [ ] Ordenar por created_at ascending
- [ ] Retornar: id, content, direction, sender, timestamp

**Task 8.4: PUT /conversations/:id (Status)**
- [ ] Permitir agent mudar status (open→closed)
- [ ] Validar autorização
- [ ] Log da mudança
- [ ] Emitir evento via WS

**Task 8.5: GET /health**
- [ ] Retornar status: ok/degraded/error
- [ ] Verificar DB connectivity
- [ ] Verificar Redis connectivity
- [ ] Retornar versão da app

---

### FASE 9: OBSERVABILIDADE (Sprint 5)
**Duração:** 1 semana | **Entregáveis:** Logs, métricas, tracing básico

#### 2.9.1 Tasks

**Task 9.1: Structured Logging**
- [ ] Configurar slog com output JSON
- [ ] Adicionar request ID em toda requisição (middleware)
- [ ] Logar: método, path, status, latência, user_id
- [ ] Logar erros com stack trace
- [ ] Definir log levels por componente

**Task 9.2: Metrics**
- [ ] Adicionar Prometheus client
- [ ] Métricas: request count, latência, erros por tipo
- [ ] Métricas: DB queries, Redis ops
- [ ] Métricas: WebSocket conexões ativas
- [ ] Endpoint GET /metrics

**Task 9.3: Tracing**
- [ ] Setup básico de tracing (jaeger ou similar)
- [ ] Adicionar trace ID em logs
- [ ] Logar spans importantes (DB, Redis, HTTP calls)

**Task 9.4: Alertas**
- [ ] Definir alertas: erro rate > 1%, latência > 500ms
- [ ] Definir alertas: DB conexões > 80%, Redis memory > 80%
- [ ] Configuração de alertas em docker-compose

---

### FASE 10: TESTES E QUALIDADE (Sprint 5)
**Duração:** 1 semana | **Entregáveis:** Coverage > 80%, testes E2E

#### 2.10.1 Tasks

**Task 10.1: Unit Tests**
- [ ] Cobertura de todos os use cases
- [ ] Cobertura de handlers HTTP
- [ ] Cobertura de providers
- [ ] Uso de mocks e stubs
- [ ] Target: > 80% coverage

**Task 10.2: Integration Tests**
- [ ] Testes com PostgreSQL real (testcontainers)
- [ ] Testes com Redis real
- [ ] Testes de fluxo completo (receber→processar→enviar)

**Task 10.3: E2E Tests**
- [ ] Docker Compose setup para testes
- [ ] Simular WebSocket client
- [ ] Simular webhook do WhatsApp
- [ ] Testar múltiplos cenários

**Task 10.4: Code Quality**
- [ ] Setup de linting (golangci-lint)
- [ ] Setup de formatter (gofmt)
- [ ] Code review checklist
- [ ] Documentação de código (godoc)

---

### FASE 11: DOCUMENTAÇÃO (Sprint 5)
**Duração:** 1 semana | **Entregáveis:** Documentação completa

#### 2.11.1 Tasks

**Task 11.1: API Documentation**
- [ ] Swagger/OpenAPI spec completo
- [ ] Todos os endpoints documentados
- [ ] Modelos de request/response
- [ ] Códigos de erro documentados
- [ ] Setup de Swagger UI em /api/docs

**Task 11.2: README técnico**
- [ ] Setup local com docker-compose
- [ ] Rodando testes
- [ ] Variáveis de ambiente
- [ ] Estrutura de pastas
- [ ] Comandos importantes (Makefile)

**Task 11.3: Architecture Decision Records (ADRs)**
- [ ] ADR-002: Padrão Hexagonal
- [ ] ADR-003: WebSocket vs polling
- [ ] ADR-004: Cache strategy (Redis)

**Task 11.4: Guides**
- [ ] Guide: Como adicionar novo endpoint
- [ ] Guide: Como testar localmente
- [ ] Guide: Como fazer deploy
- [ ] Guide: Troubleshooting comum

---

### FASE 12: DEPLOYMENT (Sprint 5)
**Duração:** 1 semana | **Entregáveis:** MVP em produção

#### 2.12.1 Tasks

**Task 12.1: Hardening**
- [ ] Validação rigorosa de input
- [ ] SQL injection prevention (GORM already does)
- [ ] CSRF protection (se aplicável)
- [ ] HTTPS enforced
- [ ] Secrets nunca em logs

**Task 12.2: Performance**
- [ ] Profiling com pprof
- [ ] Otimização de queries (índices, N+1)
- [ ] Caching de conversas ativas em Redis
- [ ] Connection pooling DB

**Task 12.3: High Availability**
- [ ] Health checks configurados
- [ ] Graceful shutdown
- [ ] Zero-downtime deployment strategy
- [ ] Database backup strategy

**Task 12.4: Deployment**
- [ ] Build image Docker otimizada
- [ ] Push para registry
- [ ] Deploy em staging
- [ ] Deploy em produção
- [ ] Monitoring em produção

---

## 3. DOCUMENTAÇÃO NECESSÁRIA

### 3.1 Estrutura de Documentação

```
docs/
├── 00-quick-start.md          [NOVO] Como rodar localmente
├── 01-overview.md             [ATUAL]
├── 02-provider-comparison.md  [ATUAL]
├── 03-cost-benefit-decision.md [ATUAL]
├── 04-setup-meta-cloud.md     [ATUALIZAR] Com detalhes operacionais
├── 05-architecture.md         [ATUALIZAR] Com diagramas
├── 06-api-contracts.md        [EXPANDIR] Swagger ref
├── 07-database-schema.md      [EXPANDIR] Com índices, migrations
├── 08-mvp-roadmap.md          [ATUAL]
├── 09-deployment.md           [EXPANDIR] Docker, k8s, infra
├── 10-development-guide.md    [NOVO] Como trabalhar no projeto
├── 11-testing-guide.md        [NOVO] Estratégia de testes
├── 12-observability.md        [NOVO] Logging, métricas, tracing
├── 13-api-reference.md        [NOVO] OpenAPI/Swagger completo
└── adr/
    ├── adr-001-whatsapp-provider.md [ATUAL]
    ├── adr-002-hexagonal-architecture.md [NOVO]
    └── adr-003-websocket-strategy.md [NOVO]
```

### 3.2 Documentos Principais a Criar/Atualizar

**00-quick-start.md**
- Pré-requisitos (Go 1.24+, Docker, Make)
- Clone e setup: `make setup`
- Rodar localmente: `make dev`
- Rodar testes: `make test`
- Links para outros docs

**10-development-guide.md**
- Estrutura do projeto passo a passo
- Pattern hexagonal explicado
- Como adicionar novo endpoint (checklist)
- Convenções de código (Go idiomático)
- Nomes de functions, variáveis, tipos

**11-testing-guide.md**
- Estratégia de testes (unit, integration, E2E)
- Como rodar testes
- Como mockar dependências
- Coverage targets
- Como testar WebSocket

**12-observability.md**
- Structured logging
- Prometheus métricas
- Jaeger tracing setup
- Como debugar issues em produção

**13-api-reference.md**
- OpenAPI/Swagger spec
- Ou gerar automaticamente do código

---

## 4. MILESTONES E TIMELINE

### Timeline Total: 12 Semanas (Sprints de 1 semana)

| Sprint | Fase | Duração | Entregável |
|--------|------|---------|-----------|
| 1 | Setup + DB Modelos | 2 sem | Projeto + DB schema prontos |
| 2 | Auth + WhatsApp Setup | 2 sem | JWT + Provider abstrato |
| 3 | Mensagens (in/out) | 2 sem | Receber e enviar mensagens |
| 4 | WebSocket + CRUD | 2 sem | Real-time + endpoints auxiliares |
| 5 | Observabilidade + Deploy | 2 sem | MVP em produção |
| 6 | Buffer/Ajustes | 2 sem | Refinements, documentação final |

### Checkpoints Críticos
- **Semana 2:** Docker Compose rodando, migrations OK
- **Semana 4:** WhatsApp webhook recebendo mensagens
- **Semana 6:** Envio de mensagens funcionando
- **Semana 8:** WebSocket real-time ativo
- **Semana 10:** Testes automatizados > 80% coverage
- **Semana 12:** MVP em produção com observabilidade

---

## 5. CRITÉRIOS DE ACEIÇÃO POR FASE

### Fase 1: Setup ✓
- [ ] Repositório configurado com estrutura correta
- [ ] Docker Compose com 4 serviços rodando
- [ ] CI/CD pipeline executando
- [ ] README com instruções de setup

### Fase 2: DB ✓
- [ ] Migrations rodando sem erro
- [ ] Modelos GORM compilam
- [ ] Repositórios testados (unit tests)
- [ ] Schema documentado

### Fase 3: Auth ✓
- [ ] JWT sendo gerado e validado
- [ ] Middleware protegendo endpoints
- [ ] Secrets seguros no .env
- [ ] Testes de auth passando

### Fase 4: WhatsApp Setup ✓
- [ ] Meta App criado e configurado
- [ ] Credentials no .env
- [ ] Provider interface implementada
- [ ] Testes com mock provider

### Fase 5: Receber Mensagens ✓
- [ ] Webhook recebendo do WhatsApp
- [ ] Mensagens persistidas no DB
- [ ] Contatos criados automaticamente
- [ ] Testes E2E passando

### Fase 6: Enviar Mensagens ✓
- [ ] Endpoint POST /messages/send
- [ ] Mensagens sendo enviadas via Meta
- [ ] Persistência confirmada
- [ ] Erros tratados com retry

### Fase 7: WebSocket ✓
- [ ] Conexões WS aceitas e validadas
- [ ] Broadcast funcionando
- [ ] Múltiplas conexões simultâneas
- [ ] Load tests OK (100+ conexões)

### Fase 8: CRUD Endpoints ✓
- [ ] GET /conversations listando
- [ ] GET /conversations/:id com detalhes
- [ ] GET /conversations/:id/messages com histórico
- [ ] PUT para mudar status
- [ ] GET /health OK

### Fase 9: Observabilidade ✓
- [ ] Logs estruturados em JSON
- [ ] Métricas Prometheus coletando
- [ ] Dashboard Grafana configurado
- [ ] Alertas disparando corretamente

### Fase 10: Testes ✓
- [ ] Coverage > 80%
- [ ] CI passando (lint + tests)
- [ ] E2E tests em Docker
- [ ] Performance OK (< 300ms)

### Fase 11: Documentação ✓
- [ ] Swagger/OpenAPI completo
- [ ] README técnico atualizado
- [ ] Guides de desenvolvimento
- [ ] ADRs documentados

### Fase 12: Deployment ✓
- [ ] Imagem Docker otimizada
- [ ] Rodando em staging
- [ ] Rodando em produção
- [ ] Monitoring ativo

---

## 6. STACK DE DEPENDÊNCIAS GO

```
Core:
- github.com/gofiber/fiber/v3
- gorm.io/gorm
- gorm.io/driver/postgres

Auth:
- github.com/golang-jwt/jwt/v5

Cache/Messaging:
- github.com/redis/go-redis/v9

HTTP Client:
- (net/http stdlib)

Logging:
- log/slog (stdlib)

Testing:
- testing (stdlib)
- github.com/stretchr/testify

DB Migrations:
- github.com/golang-migrate/migrate/v4

Config:
- github.com/spf13/viper

Optional (futuro):
- OpenTelemetry para tracing
- Prometheus client para métricas
```

---

## 7. ARQUIVOS PRINCIPAIS A CRIAR

```
customer-support-hub/
├── README.md
├── PLANO-DESENVOLVIMENTO-API.md (este arquivo)
├── Makefile
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── .github/workflows/
│   ├── test.yml
│   ├── build.yml
│   └── deploy.yml
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── contact/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── contact_test.go
│   │   ├── conversation/
│   │   ├── message/
│   │   └── agent/
│   ├── application/
│   │   ├── receive_message/
│   │   │   ├── use_case.go
│   │   │   └── use_case_test.go
│   │   ├── send_message/
│   │   └── assign_conversation/
│   ├── infrastructure/
│   │   ├── db/
│   │   │   ├── connection.go
│   │   │   ├── postgres_contact_repo.go
│   │   │   └── migrations/
│   │   ├── cache/
│   │   │   └── redis_manager.go
│   │   ├── websocket/
│   │   │   └── manager.go
│   │   └── whatsapp/
│   │       ├── provider.go
│   │       └── meta_provider.go
│   ├── interfaces/
│   │   ├── http/
│   │   │   ├── router.go
│   │   │   ├── handlers/
│   │   │   │   ├── messages.go
│   │   │   │   ├── conversations.go
│   │   │   │   ├── webhooks.go
│   │   │   │   ├── health.go
│   │   │   │   └── handlers_test.go
│   │   │   └── middleware/
│   │   │       └── auth.go
│   │   ├── ws/
│   │   │   ├── handler.go
│   │   │   └── handler_test.go
│   │   └── dto/
│   │       └── message.go
│   ├── config/
│   │   └── config.go
│   └── logger/
│       └── logger.go
├── migrations/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── tests/
│   ├── integration/
│   └── e2e/
└── docs/
    ├── 00-quick-start.md
    ├── 10-development-guide.md
    ├── 11-testing-guide.md
    ├── 12-observability.md
    ├── 13-api-reference.md
    └── adr/
```

---

## 8. PRIORIDADES E DEPENDÊNCIAS

### Ordem Lógica (sem skip possível)
1. **Setup + Docker** → sem isto, nada roda
2. **DB + Modelos** → sem isto, não há persistência
3. **Auth (JWT)** → sem isto, não há segurança
4. **WhatsApp Provider** → sem isto, não há integração
5. **Receber Mensagens** → primeira funcionalidade
6. **Enviar Mensagens** → segunda funcionalidade
7. **WebSocket** → terceira funcionalidade (real-time)
8. **CRUD** → endpoints auxiliares
9. **Observabilidade** → necessário para produção
10. **Testes Completos** → validação antes de deploy
11. **Documentação** → conhecimento do projeto
12. **Deployment** → go-live

### Possíveis Paralelizações
- Tasks 1.2, 1.3, 1.4 podem ser paralelas (após 1.1)
- Tasks 2.1, 2.3 podem ser paralelas
- Tasks 3.1, 3.2, 3.3 podem ser paralelas
- Tasks 8.1-8.5 podem ser paralelas (após fase 7)

---

## 9. RISCOS E MITIGAÇÕES

| Risco | Impacto | Probabilidade | Mitigação |
|-------|---------|--------------|-----------|
| Atraso na aprovação Meta App | Alto | Média | Começar desenvolvimento com mock provider antes |
| Rate limiting do WhatsApp | Médio | Baixa | Implementar circuit breaker, exponential backoff |
| Performance em WebSocket | Alto | Média | Load tests frequentes, profiling |
| Indisponibilidade PostgreSQL | Alto | Baixa | Backups automáticos, replicação (futuro) |
| Indisponibilidade Redis | Médio | Baixa | Fallback para sessões em DB (degradação) |
| Secret exposure em logs | Alto | Baixa | Code review, scanning automático |
| Schema instável no início | Médio | Alta | Migrations reversíveis, versionamento |

---

## 10. SUCESSO DO MVP

### Critérios de Sucesso
- ✅ Receber mensagens do WhatsApp em < 2s
- ✅ Enviar mensagens em < 2s
- ✅ Atualizar UI em tempo real via WebSocket (< 500ms)
- ✅ Handle 100+ conversas simultâneas
- ✅ Latência p99 < 300ms
- ✅ Uptime > 99%
- ✅ Error rate < 0.1%
- ✅ Cobertura de testes > 80%

### Métricas de Monitoramento
- Requisições por segundo
- Latência (p50, p95, p99)
- Taxa de erro por endpoint
- Conexões WebSocket ativas
- Pool de DB utilizado
- Memória Redis utilizado
- Disk usage PostgreSQL

---

## 11. PRÓXIMOS PASSOS IMEDIATOS

### Semana 1 - Actions:
1. [ ] Criar repositório GitHub
2. [ ] Estrutura de pastas conforme ADN
3. [ ] go.mod com dependências base
4. [ ] Dockerfile + docker-compose.yml
5. [ ] .env.example
6. [ ] Makefile com targets: setup, dev, test, build
7. [ ] README com quick-start
8. [ ] Primeiro commit

### Verificação:
- [ ] `make dev` inicia todos os serviços
- [ ] `make test` passa sem erro
- [ ] `make build` gera imagem Docker

---

**Próxima Review:** Após Semana 2 (Setup + DB modelos prontos)

---
