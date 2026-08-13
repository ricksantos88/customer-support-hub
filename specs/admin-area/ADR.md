# ADR-001: Gestão de Agentes e Sessões Administrativas na Mesma Base de Código

## Contexto
O projeto precisa de uma área administrativa para gerenciar agentes (CRUD) e revogar sessões ativas. Atualmente, o ecossistema é baseado em Go com Fiber e PostgreSQL + Redis, estruturado com Arquitetura Hexagonal. Havia uma dúvida se a lógica administrativa deveria ser um microsserviço isolado ou incorporada ao backend existente.

## Decisão
Implementar a área administrativa como um conjunto de endpoints protegidos por rota (`/admin/*`) dentro do mesmo backend atual, utilizando o middleware de autenticação existente e aplicando validações específicas de privilégio (`RequireRole("admin")`) a nível de middleware HTTP.

## Consequências
* **Prós**:
  - Reutilização completa da infraestrutura de banco de dados (tabelas `agents` e `sessions`).
  - Sem necessidade de criar um novo deploy ou de gerenciar autenticação inter-serviços.
  - Facilidade de desenvolvimento rápida compartilhando modelos GORM e repositórios existentes.
* **Contras**:
  - O backend passa a conter lógica de negócio tanto para agentes de suporte quanto para administradores do sistema, aumentando o tamanho do binário.
