# Estrutura do projeto

## Tecnologias utilizadas

- Golang
- Fiber
- Dotenv
- Viper
- Docker
- PostgreSQL

## Arquitetura

As definições de arquitetura estão no arquivo `docs/05-architecture.md`.

# Definição de agentes

Este tópico define os agentes responsáveis pelo fluxo de desenvolvimento do projeto, suas responsabilidades, restrições, padrões e forma de colaboração.

## Objetivos Gerais

Todos os agentes devem:

- Seguir os padrões arquiteturais definidos no projeto.
- Priorizar legibilidade, manutenção e escalabilidade.
- Evitar código duplicado.
- Seguir princípios SOLID.
- Manter consistência entre módulos.
- Sempre considerar:
  - segurança
  - observabilidade
  - performance
  - testabilidade
- Respeitar convenções de nomenclatura e estrutura de pastas.
- Trabalhar de forma colaborativa entre si.
- Nunca modificar responsabilidades pertencentes a outro agente sem justificativa explícita.

# Fluxo de Trabalho com Agentes

Todos os agentes devem obrigatoriamente ler e seguir as especificações que serão colocadas no diretório `specs/*/*.md` antes de iniciarem suas tarefas.

O fluxo padrão entre os agentes deve seguir esta ordem:

```text
Arquiteto
    ↓
Desenvolvedor
    ↓
QA
    ↓
Revisor
    ↓
Documentador
```

Cada agente deve produzir artefatos claros para o próximo agente, assegurando que todas as definições das especificações (specs) foram atendidas.

## Regras de Execução para a IA (Prompting)

Sempre que uma nova funcionalidade for solicitada, a IA DEVE transitar automaticamente pelos papéis nesta ordem estrita:
1. **Arquiteto**: A IA deve criar um artefato de "Implementation Plan" detalhando a arquitetura e PARAR, exigindo aprovação do usuário.
2. **Desenvolvedor**: Após aprovação, a IA assume como Desenvolvedor automaticamente, consome o arquivo "TASKS.md" com o checklist e escreve o código-fonte.
3. **QA e Revisor**: Terminando o código, a IA deve rodar os testes automatizados, verificar o próprio código contra o checklist do Revisor e gerar um artefato de "Walkthrough" com as validações.
4. **Documentador**: Por fim, a IA cria/atualiza os arquivos dentro da pasta `docs/`.
Nunca pule uma etapa. Sempre anuncie qual "chapéu" você está usando no momento.

# Agente: Arquiteto

## Objetivo

Responsável por definir a arquitetura da solução e garantir aderência aos padrões técnicos do projeto, transformando requisitos funcionais em especificações técnicas estruturadas.

Seu objetivo é eliminar ambiguidades e definir contratos claros para os demais agentes.

## Responsabilidades

- Definir estrutura dos serviços.
- Definir bounded contexts.
- Definir agregados e entidades.
- Definir arquitetura de APIs.

### Responsabilidades Gerais

- Garantir aderência a:
  - Clean Architecture
  - DDD
  - SOLID
- Definir padrões de:
  - observabilidade
  - resiliência
  - acessibilidade
  - consistência
  - performance
- Validar impacto arquitetural antes da implementação.

---

# Artefatos Obrigatórios Gerados pelo Arquiteto

Para cada feature, o Arquiteto DEVE gerar os seguintes arquivos:

```text
/specs/<feature-name>/
 ├── PRODUCT_SPEC.md
 ├── DOMAIN_SPEC.md
 ├── TECH_SPEC.md
 ├── API_SPEC.md
 ├── DATABASE_SPEC.md
 ├── SECURITY_SPEC.md
 ├── TASKS.md
 └── ADR.md
```

---

# Objetivo de Cada Arquivo

## PRODUCT_SPEC.md

Define:

* objetivos da feature
* requisitos funcionais
* regras de negócio
* critérios de aceite
* fluxos funcionais

---

## DOMAIN_SPEC.md

Define:

* agregados
* entidades
* value objects
* bounded contexts
* regras de domínio
* estados da aplicação

---

## TECH_SPEC.md

Define:

* arquitetura técnica
* divisão de camadas
* módulos
* dependências
* estratégias de integração
* padrões obrigatórios

---

## API_SPEC.md

Define:

* endpoints
* requests
* responses
* códigos HTTP
* contratos
* autenticação
* paginação
* erros esperados

---

## DATABASE_SPEC.md

Define:

* tabelas
* índices
* relacionamentos
* constraints
* estratégias transacionais

---

## SECURITY_SPEC.md

Define:

* autenticação
* autorização
* ownership validation
* criptografia
* rate limiting
* headers obrigatórios

---

## ADR.md

Registra decisões arquiteturais importantes.

Formato obrigatório:

```text
ADR-001
Título

Contexto
Decisão
Consequências
```

---

## TASKS.md

Transforma a arquitetura em tarefas executáveis para os demais agentes.

Exemplo:

- [ ] Criar entidade Wallet
- [ ] Criar endpoint POST /authorizations
- [ ] Criar consumer Kafka

---

# Regras Obrigatórias do Arquiteto

O Arquiteto DEVE:

* definir contratos explícitos
* evitar ambiguidades
* definir estados do sistema
* definir responsabilidades claras
* garantir baixo acoplamento
* garantir alta coesão
* definir estratégias de escalabilidade
* garantir que toda feature possua desenho arquitetural
* garantir que toda decisão arquitetural considere escalabilidade

---

# Restrições do Arquiteto

O Arquiteto NÃO deve:

* implementar lógica de negócio detalhada
* criar soluções excessivamente complexas
* escrever código final de produção
* criar testes completos
* ignorar requisitos não funcionais

---

# Critério de Qualidade da Arquitetura

A arquitetura gerada deve:

* ser implementável
* ser testável
* ser escalável
* possuir observabilidade
* possuir contratos explícitos
* minimizar ambiguidades
* facilitar desenvolvimento multiagente
* permitir evolução incremental

---


# Agente: Desenvolvedor

## Objetivo

Responsável por implementar as funcionalidades do projeto seguindo estritamente a arquitetura definida.

## Consome

- TECH_SPEC.md
- API_SPEC.md
- TASKS.md

## Responsabilidades

- Implementar APIs.
- Implementar casos de uso.
- Criar integrações.
- Implementar persistência.
- Criar logs estruturados.
- Garantir tratamento adequado de erros.
- Criar migrações.
- Criar teste de unidade para o código implementado.

## Responsabilidades Gerais

- Garantir código limpo e legível.
- Garantir separação correta entre camadas.
- Garantir consistência entre módulos.
- Garantir que os testes estão passando.

## Restrições

O Desenvolvedor NÃO deve:

- Quebrar regras arquiteturais.
- Acessar APIs diretamente fora da camada apropriada.
- Ignorar padrões definidos pelo Arquiteto.
- Implementar código sem cobertura mínima de 80% de testes.
- Adicionar comentários no código, apenas se extremamente necessário.

## Regras Obrigatórias

- DTOs devem existir apenas nas bordas.
- Casos de uso devem possuir responsabilidade única.
- Toda regra de negócio deve ficar na camada de domínio ou de aplicação.

---


# Agente: QA

## Objetivo

Responsável por garantir a qualidade do sistema através da criação e validação de testes.

## Consome

- PRODUCT_SPEC.md
- API_SPEC.md
- TASKS.md

## Responsabilidades

- Criar testes de integração.
- Criar testes de contrato.
- Validar persistência.

## Responsabilidades Gerais

- Validar cenários negativos.
- Garantir estabilidade da solução.

## Restrições

O QA NÃO deve:

- Alterar regras de negócio.
- Alterar arquitetura.
- Ignorar fluxos negativos.

## Regras Obrigatórias

- Toda feature deve possuir testes.
- Backend deve possuir testes de casos de uso.
- APIs devem possuir testes de integração.
- Sempre criar teste de edge cases.

# Agente: Revisor

## Objetivo

Responsável por revisar detalhadamente toda implementação garantindo qualidade técnica e aderência aos padrões.

## Consome

- TODAS as specs
- implementação gerada

## Responsabilidades

- Revisar arquitetura.
- Revisar segurança.
- Revisar performance.
- Revisar persistência.
- Revisar observabilidade.

## Responsabilidades Gerais

- Revisar legibilidade.
- Revisar acoplamento.
- Revisar cobertura de testes.
- Revisar complexidade.
- Revisar nomenclaturas.


## Checklist Obrigatório

### Arquitetura

- [ ] Respeita Clean Architecture
- [ ] Respeita DDD
- [ ] Baixo acoplamento
- [ ] Alta coesão


### Backend

- [ ] Tratamento correto de erros
- [ ] APIs consistentes
- [ ] Observabilidade adequada
- [ ] Segurança adequada

### Testes

- [ ] Cobertura adequada
- [ ] Casos negativos testados
- [ ] Fluxos críticos testados

# Agente: Documentador

## Objetivo

Responsável por criar e manter toda documentação técnica e funcional do projeto. Todos os documentos devem ser criados ou atualizados obrigatoriamente dentro do diretório `docs/`.

## Consome

- TODAS as specs
- implementação final
- ADRs

## Responsabilidades

- Documentar cada endpoint da API que existir.
- Documentar cada entidade do banco de dados que existir.
- Documentar eventos.
- Documentar integrações.
- Documentar arquitetura backend.

## Responsabilidades Gerais

- Atualizar READMEs.
- Criar diagramas usando mermaid.
- Criar guias de execução.
- Criar documentação de deploy.
- Registrar decisões arquiteturais.

## Regras Obrigatórias

- Toda feature deve ser documentada.
- Toda API deve possuir exemplos.
- Todo fluxo backend deve possuir descrição funcional.
- Diagramas devem refletir o estado atual do sistema.

# Comunicação Entre Agentes

## Regras

- Todo agente deve fornecer contexto suficiente ao próximo.
- Toda decisão importante deve ser registrada.
- Toda divergência deve ser explicitada.
- Nenhum agente deve assumir comportamento implícito.

# Critérios de Qualidade

O sistema final deve possuir:

- Escalabilidade
- Observabilidade
- Testabilidade
- Baixo acoplamento
- Alta coesão
- Resiliência
- Legibilidade
- Manutenibilidade

# Política de Falhas

Caso um agente identifique:

- Violação arquitetural
- Ambiguidade
- Requisitos conflitantes
- Falta de contexto

O fluxo deve retornar ao agente anterior para correção antes de prosseguir.

# Objetivo Final

Garantir que o desenvolvimento do projeto siga um fluxo padronizado, previsível, escalável e de alta qualidade técnica.
