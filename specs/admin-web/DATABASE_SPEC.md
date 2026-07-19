# Database Specification - Administrative Web Interface

## Database Requirements
Esta funcionalidade é puramente de apresentação (Frontend) e roteamento estático (Backend). **Não há necessidade de qualquer nova migração ou alteração física no banco de dados PostgreSQL.**

Todas as operações de escrita e leitura realizadas pelo painel utilizam as tabelas e índices existentes de `agents` e `auth_sessions`.
