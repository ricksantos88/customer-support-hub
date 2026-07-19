# Database Specification - Administrative Area

## Schema Analysis
As tabelas necessárias para suportar a área administrativa já existem no banco de dados Postgres e não exigem alterações no schema principal:

### 1. Tabela `agents`
* Mapeamento GORM no struct `Agent` ([agent.go](file:///Users/gedalias.caldas/Documents/customer-support-hub/internal/models/agent.go)).
* A coluna `role` possui restrição check (`chk_agents_role`) para os valores `('admin', 'agent')`.

### 2. Tabela `sessions`
* Mapeamento GORM no struct `Session` ([session.go](file:///Users/gedalias.caldas/Documents/customer-support-hub/internal/models/session.go)).

## Transaction Strategies
1. **Deactivate Agent**:
   - Quando um agente é deletado ou inativado (soft delete ou deativação lógica), todas as suas sessões devem ser revogadas em uma transação única do banco de dados:
     ```sql
     BEGIN;
     UPDATE agents SET deleted_at = NOW() WHERE id = agent_id;
     UPDATE sessions SET revoked_at = NOW() WHERE agent_id = agent_id AND revoked_at IS NULL;
     COMMIT;
     ```
   - Em seguida, o cache das sessões associadas a esse agente deve ser invalidado no Redis.

## Indexes Used
* Índice único em `agents(email)` para garantir a unicidade de e-mails.
* Índice em `sessions(agent_id)` para rápida busca e revogação das sessões de um agente.
* Índice em `sessions(refresh_token_hash)` para controle de revogação/validação.
