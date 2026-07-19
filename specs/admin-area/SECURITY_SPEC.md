# Security Specification - Administrative Area

## Authentication & Authorization
* **Bearer Token Validation**: A API administrativa utiliza exclusivamente o middleware de Bearer Token existente, validando as assinaturas dos tokens de acesso gerados para os agentes.
* **Role Enforcement Middleware**: Apenas requisições que contenham a claim de JWT `"agent_role": "admin"` no token serão permitidas. O middleware `RequireRole("admin")` impede qualquer requisição com claim de `"agent"` de prosseguir, retornando imediatamente **403 Forbidden**.

## Rate Limiting
* Aplicar o middleware de rate limit padrão:
  * **IP-based limit**: Bloqueia requisições em massa por IP para os endpoints `/admin/*`.
  * **Agent-based limit**: Limita requisições por Token JWT do agente.

## Specific Safety Controls
1. **No Self-Revocation/Self-Demotion**:
   - Um administrador não pode usar o endpoint `PUT /admin/agents/:id` para atualizar a sua própria `role` de `admin` para `agent`.
   - Um administrador não pode chamar `DELETE /admin/agents/:id` informando o seu próprio ID para evitar auto-exclusão e travamento do sistema administrativo.
2. **Immediate Invalidation**:
   - Ao desativar ou deletar um agente, todas as chaves de sessões ativas do respectivo agente devem ser limpas do Redis em lote. Isso garante que requisições subsequentes com tokens de acesso antigos sejam rejeitadas mesmo antes da expiração padrão de 15 minutos do JWT.

## Required HTTP Headers
* `X-Content-Type-Options: nosniff`
* `X-Frame-Options: DENY`
* `X-XSS-Protection: 1; mode=block`
