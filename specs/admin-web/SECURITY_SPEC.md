# Security Specification - Administrative Web Interface

## Token Security
- O `access_token` e `refresh_token` são salvos localmente no navegador via `localStorage`.
- Como a segurança é validada no backend por token JWT assinado, invasões ou modificações no JS local não darão acesso a dados confidenciais se a sessão expirar ou for inválida.

## Self-demotion & Self-deletion Safety in UI
- No painel administrativo, a opção de deletar e alterar a role de um agente (de `admin` para `agent`) estará visivelmente desabilitada ou oculta para a linha correspondente ao administrador logado no momento.

## Cross-Origin Resource Sharing (CORS)
- Sendo servido na mesma origem (`http://localhost:8080/admin-panel/`), não há risco de ataques de CORS ou necessidade de habilitar origens adicionais para o painel administrativo.
