# API Specification - Administrative Web Interface

## Request Handling
Todas as requisições AJAX feitas pelo frontend devem incluir o cabeçalho Authorization:
```javascript
const headers = {
  'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
  'Content-Type': 'application/json'
};
```

## Error & Expiry Handling
Se qualquer endpoint retornar `401 Unauthorized`:
1. O frontend tenta renovar o token fazendo requisição para `/auth/refresh` com o `refresh_token` salvo.
2. Se a renovação falhar, o frontend executa o logout limpando o `localStorage` e redirecionando para `/admin-panel/`.

Se retornar `403 Forbidden`:
- Exibe mensagem informando que a conta não possui privilégios de administrador.
