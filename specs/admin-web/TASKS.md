# Tasks - Administrative Web Interface

- [ ] **Configuração e Roteamento no Backend**
  - [ ] Mapear rota estática `/admin-panel` no roteador do Fiber (`router.go`) para servir `./web/admin`
  - [ ] Modificar o `Dockerfile` para copiar a pasta `./web` na imagem final de produção

- [ ] **Desenvolvimento Frontend (HTML/CSS)**
  - [ ] Criar arquivo HTML básico (`index.html`) com contêineres para Login, Dashboard, Tabelas de Agentes e Sessões, e Modal de CRUD
  - [ ] Criar arquivo CSS moderno (`styles.css`) aplicando design dark-mode, glassmorphism, tipografia Inter, transições suaves e layout responsivo

- [ ] **Programação da Interface (JavaScript)**
  - [ ] Desenvolver `app.js` encapsulando as chamadas AJAX (`fetch`) no backend
  - [ ] Implementar fluxo de Login e salvamento de tokens
  - [ ] Implementar fluxo de carregamento dinâmico de agentes, sessões e diagnósticos do sistema
  - [ ] Implementar abertura/fechamento de modais e chamadas de escrita (criação, edição e exclusão de agentes; revogação de sessões)
  - [ ] Ocultar ações de auto-exclusão e auto-demissão na linha do admin logado

- [ ] **Verificação e Testes**
  - [ ] Validar carregamento correto no Docker e pelo localhost
  - [ ] Testar cenários de login, acessos negados com tokens inválidos, e expiração de sessões
