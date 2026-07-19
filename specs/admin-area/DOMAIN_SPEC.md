# Domain Specification - Administrative Area

## Bounded Context
* **Access Control & Management (Admin Context)**: Lida com a gestão do ciclo de vida dos agentes e gerenciamento de suas sessões ativas.

## Aggregate Roots & Entities
### 1. Agent (Aggregate Root - Já existente no domínio)
* **Entidade**: `Agent`
* **Campos principais**: `ID`, `Name`, `Email`, `PasswordHash`, `Role` (admin/agent), `CreatedAt`, `LastActive`.
* **Regras de Validação do Domínio**:
  - `Email` formatado e único.
  - `Role` válido (restrito a `admin` ou `agent`).
  - `PasswordHash` gerado via Bcrypt.

### 2. Session (Aggregate - Já existente)
* **Entidade**: `Session`
* **Campos principais**: `ID`, `AgentID`, `RefreshTokenHash`, `UserAgent`, `IPAddress`, `CreatedAt`, `LastUsedAt`, `ExpiresAt`, `RevokedAt`.

## Value Objects
* **Role**: Representa a função/permissão do agente no sistema. Contém validação dos valores permitidos (`admin`, `agent`).

## Domain Rules (Regras de Negócio do Domínio)
1. **Uniqueness**: O e-mail de um agente deve ser único em toda a base de dados.
2. **Immutability of Admin Safety**: Não é permitido desativar a própria conta de administrador nem alterar a própria role para `agent`.
3. **Session Revocation**: A revogação de uma sessão define o campo `RevokedAt` com o timestamp atual. Uma sessão que possui `RevokedAt` preenchido é considerada inativa de forma permanente.

## Application States (Estados)
* **Agent State**:
  - `Active`: Agente cadastrado e liberado para login.
  - `Inactive/Deleted`: Agente com credenciais invalidadas ou removido.
* **Session State**:
  - `Active`: `ExpiresAt` no futuro E `RevokedAt` nulo.
  - `Expired`: `ExpiresAt` no passado.
  - `Revoked`: `RevokedAt` preenchido.
