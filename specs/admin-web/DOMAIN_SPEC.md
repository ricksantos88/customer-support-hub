# Domain Specification - Administrative Web Interface

## Bounded Context
* **Admin Web Presentation Context**: Apresenta e visualiza os agregados e entidades gerenciados no contexto administrativo de maneira amigável.

## Domain Presentation States
* **Frontend Authentication State**:
  - `Authenticated`: Contém um `access_token` válido e dados do admin logado.
  - `Unauthenticated`: Redireciona para o formulário de login.
* **System Status Indicator**:
  - `Healthy`: Todos os serviços conectados.
  - `Degraded`: Redis ou WhatsApp com falha de conexão/configuração.
  - `Critical`: Banco de dados offline.

## Mapped Aggregate Views
* **Agent List View**: Representa a coleção de entidades `Agent` ativas.
* **Active Sessions List View**: Representa a coleção de entidades `Session` ativas.
