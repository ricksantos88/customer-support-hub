# Technical Specification - Administrative Web Interface

## Architecture Style
Sendo uma Single Page Application (SPA) minimalista servida pelo backend (Opção A), a interface administrativa residirá no diretório raiz:
```text
web/
  admin/
    assets/
    index.html
    styles.css
    app.js
```

## Front-End Design System & Palette
- **Palette**: Slate-dark background (`#0d1117`), dark-card backgrounds (`#161b22`), glass borders (`#30363d`), vibrant accents (Neon Green `#39d353` and Neon Blue `#00b4d8`).
- **Typography**: Google Font `Inter`.
- **Interactive UI**: CSS transitions, hover scales, loading animations, status indicator pulses, and responsive grids.

## Backend Serving Strategy
- Mapear a rota `/admin-panel` no Fiber:
  ```go
  app.Static("/admin-panel", "./web/admin")
  ```
- No `Dockerfile`, copiar a pasta `./web` no estágio final:
  ```dockerfile
  COPY --from=builder /app/web ./web
  ```

## API Interaction Layer
O arquivo `app.js` encapsula todas as requisições AJAX (`fetch`) no backend e gerencia os tokens salvos no `localStorage`:
- `POST /auth/login`
- `GET /admin/status`
- `GET /admin/agents`
- `POST /admin/agents`
- `PUT /admin/agents/:id`
- `DELETE /admin/agents/:id`
- `GET /admin/sessions`
- `DELETE /admin/sessions/:id`
