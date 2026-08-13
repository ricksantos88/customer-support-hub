# Technical Specification - Administrative Area

## Technical Architecture Overview
Seguindo a **Clean Architecture** e **Hexagonal Architecture** adotada no projeto, o gerenciamento administrativo será implementado isoladamente nas camadas adequadas:

```text
internal/
  application/
    admin/                <- [NEW] Admin Use Cases
      service.go
  interfaces/
    http/
      handlers/
        admin_handler.go  <- [NEW] HTTP REST Endpoints for Admin
```

## Layer Distribution

### 1. Interfaces Layer (http/handlers)
* **`AdminHandler`**: Responsável por receber as requisições HTTP REST, mapear os DTOs de entrada, chamar o caso de uso e retornar as respostas JSON.
* **Middlewares**: As rotas administrativas usarão obrigatoriamente:
  - `BearerAuthMiddleware`: Para autenticar o agente e preencher o contexto.
  - `RequireRole("admin")`: Para proibir agentes comuns de chamarem os endpoints.

### 2. Application Layer (application/admin)
* **`AdminService`**: Orquestra os Casos de Uso Administrativos.
  - `CreateAgent(ctx, name, email, password, role)`
  - `ListAgents(ctx, limit, offset)`
  - `UpdateAgent(ctx, id, fields)`
  - `DeleteAgent(ctx, id)`
  - `ListActiveSessions(ctx)`
  - `RevokeSession(ctx, sessionID)`
  - `GetSystemStatus(ctx)`

### 3. Infrastructure & Repository Layer
* Reutiliza as interfaces de repositórios existentes (`AgentRepository`, `SessionRepository`) estendidas ou injetadas diretamente no `AdminService`.

## Dependencies
* `github.com/gofiber/fiber/v3` (HTTP Framework)
* `github.com/google/uuid` (IDs únicos)
* `gorm.io/gorm` (Persistência no DB)
* `github.com/redis/go-redis/v9` (Para invalidação do cache de sessões)

## Design Patterns & Standards
* **Dependency Injection (DI)**: Injetar repositórios e serviços via construtores padrão (`NewService`, `NewAdminHandler`).
* **Input Validation**: Validar DTOs utilizando pacotes padrão de validação ou hooks de domínio.
* **Separation of Concerns**: DTOs devem ser mapeados apenas nas bordas (Interfaces). Casos de Uso no Application Layer recebem parâmetros primitivos ou entidades limpas.
