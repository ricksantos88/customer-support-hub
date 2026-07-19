# ADR-002: Roteamento de Frontend Administrativo via Fiber Static

## Contexto
O projeto necessita de um frontend simples e moderno para gerenciar agentes e sessões de suporte. Precisávamos decidir se criaríamos um repositório isolado ou integraríamos a entrega no mesmo servidor do backend Go.

## Decisão
Servir arquivos estáticos (HTML5, Vanilla CSS3, Javascript Puro) diretamente do backend em Go usando o recurso `app.Static` do Fiber, sob a rota `/admin-panel`.

## Consequências
* **Prós**:
  - Deploy em imagem Docker única.
  - Zero configuração de CORS.
  - Baixo consumo de recursos e latência mínima.
* **Contras**:
  - Limita a adoção imediata de frameworks frontend pesados, necessitando JavaScript puro para manipulação de DOM (o que é mais do que suficiente e performático para o escopo do MVP).
