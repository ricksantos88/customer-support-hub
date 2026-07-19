# API Specification - Administrative Area

## Authentication & Base Path
Todos os endpoints listados abaixo requerem autenticação via JWT Bearer Token e permissão de administrador:
```http
Authorization: Bearer <token_jwt_do_admin>
```
Prefix base path: `/admin`

---

## 1. POST /admin/agents
Cria um novo agente no sistema.

* **Request Body**:
```json
{
  "name": "Carlos Souza",
  "email": "carlos@example.com",
  "password": "strong-password-123",
  "role": "agent"
}
```

* **Response (201 Created)**:
```json
{
  "id": "e492b450-482a-4db3-96b0-b98a1c97a552",
  "name": "Carlos Souza",
  "email": "carlos@example.com",
  "role": "agent",
  "created_at": "2026-07-19T00:00:00Z"
}
```

* **Response (400 Bad Request)**: E-mail inválido, senha vazia, role incorreto.
* **Response (409 Conflict)**: E-mail já cadastrado.

---

## 2. GET /admin/agents
Lista todos os agentes cadastrados no sistema (com paginação).

* **Query Parameters**:
  * `limit` (default: 10, max: 100)
  * `offset` (default: 0)

* **Response (200 OK)**:
```json
[
  {
    "id": "64205176-37c6-4d51-bdd3-e5e086e0bb7a",
    "name": "Alice Silva",
    "email": "alice@example.com",
    "role": "admin",
    "created_at": "2026-07-19T00:00:00Z"
  },
  {
    "id": "e492b450-482a-4db3-96b0-b98a1c97a552",
    "name": "Carlos Souza",
    "email": "carlos@example.com",
    "role": "agent",
    "created_at": "2026-07-19T00:00:00Z"
  }
]
```

---

## 3. PUT /admin/agents/:id
Atualiza dados cadastrais de um agente.

* **Request Body**:
```json
{
  "name": "Carlos Souza Atualizado",
  "email": "carlos.novo@example.com",
  "role": "admin"
}
```

* **Response (200 OK)**:
```json
{
  "id": "e492b450-482a-4db3-96b0-b98a1c97a552",
  "name": "Carlos Souza Atualizado",
  "email": "carlos.novo@example.com",
  "role": "admin",
  "updated_at": "2026-07-19T00:05:00Z"
}
```

* **Response (400 Bad Request)**: Tentativa de demotivar a si mesmo (admin para agent).
* **Response (404 Not Found)**: Agente não localizado.

---

## 4. DELETE /admin/agents/:id
Inativa/Remove um agente do sistema e revoga imediatamente todas as suas sessões ativas.

* **Response (204 No Content)**: Sucesso.
* **Response (400 Bad Request)**: Tentativa de remover a si mesmo (auto-bloqueio).
* **Response (404 Not Found)**: Agente não localizado.

---

## 5. GET /admin/sessions
Lista todas as sessões ativas no sistema.

* **Response (200 OK)**:
```json
[
  {
    "id": "78205176-37c6-4d51-bdd3-e5e086e0bb7b",
    "agent_id": "64205176-37c6-4d51-bdd3-e5e086e0bb7a",
    "agent_name": "Alice Silva",
    "ip_address": "127.0.0.1",
    "user_agent": "Mozilla/5.0...",
    "created_at": "2026-07-19T00:00:00Z",
    "expires_at": "2026-07-20T00:00:00Z"
  }
]
```

---

## 6. DELETE /admin/sessions/:id
Revoga uma sessão ativa remotamente.

* **Response (204 No Content)**: Sucesso.
* **Response (404 Not Found)**: Sessão não encontrada.

---

## 7. GET /admin/status
Verifica integridade do banco de dados, conexão com o Redis e presença das chaves de integração do WhatsApp.

* **Response (200 OK)**:
```json
{
  "database": "connected",
  "redis": "connected",
  "whatsapp_integration": "configured"
}
```
