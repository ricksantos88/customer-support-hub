# Tasks - Administrative Area Implementation

- [ ] **Setup de DTOs e Contratos**
  - [ ] Criar DTOs para criação de agente (`CreateAgentRequest`)
  - [ ] Criar DTOs para atualização de agente (`UpdateAgentRequest`)
  - [ ] Criar DTOs para resposta do status do sistema (`SystemStatusResponse`)

- [ ] **Camada de Aplicação (Application Service)**
  - [ ] Criar interface e estrutura para o `AdminService` em `internal/application/admin`
  - [ ] Implementar caso de uso de criação de agente
  - [ ] Implementar caso de uso de listagem de agentes (com paginação)
  - [ ] Implementar caso de uso de edição de agente (com regras de negócio de não auto-demissão)
  - [ ] Implementar caso de uso de deleção de agente (com invalidação de sessões e cache Redis)
  - [ ] Implementar caso de uso de listagem de sessões ativas
  - [ ] Implementar caso de uso de revogação de sessão remota
  - [ ] Implementar diagnóstico de status da aplicação

- [ ] **Camada de Interface HTTP (HTTP Handlers & Router)**
  - [ ] Criar `AdminHandler` em `internal/interfaces/http/handlers/admin_handler.go`
  - [ ] Vincular endpoints de `/admin/*` às rotas protegidas pelo middleware `RequireRole("admin")` em `internal/interfaces/http/router.go`
  - [ ] Registrar rotas administrativas na documentação interativa (Swagger) do `swaggor`

- [ ] **Testes de Qualidade (QA & Unit Tests)**
  - [ ] Escrever testes unitários cobrindo todos os casos de uso no `AdminService` (com mocks de repositório)
  - [ ] Validar cobertura de testes acima de 80%
  - [ ] Escrever testes de integração/E2E validando as restrições de permissão do middleware `RequireRole("admin")`
