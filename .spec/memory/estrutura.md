# Estrutura do Projeto

## Estrutura alvo do repositório

```text
.spec/
  changes/
  memory/
  shared/
  templates/
apps/
  frontend/
    src/
      app/
      shared/
        components/ui/  # shadcn/ui
  backend/
    cmd/
      server/
    config/
    scripts/
    internal/
      modules/
        customer/
          models/
          services/
          repositories/
          handlers/
```

## Responsabilidades

- `.spec/changes` — specs de mudanças específicas
- `.spec/memory` — contexto global do projeto
- `.spec/shared` — convenções reutilizáveis entre specs
- `.spec/templates` — modelos para criar novas mudanças
- `apps/frontend` — telas, formulários e integração com a API (Next.js)
- `apps/backend/cmd/server` — entry point da aplicação Go
- `apps/backend/internal/modules/<dominio>` — regras de negócio por domínio

## Organização de módulos

- um módulo por área de negócio relevante
- regras de negócio primeiro, detalhes técnicos depois
- cada módulo segue: `models/`, `services/`, `repositories/`, `handlers/`

## Limites entre camadas

- o front-end (Next.js) não conhece banco de dados
- o back-end (Go/chi) expõe casos de uso via API REST
- regras de negócio não dependem diretamente da interface web
- a spec descreve a mudança **antes** da implementação

## Convenções para `apps/frontend/src/shared`

- `shared` contém apenas código reutilizável e sem acoplamento ao estado ou configuração da aplicação atual
- configurações, navegação, chaves de storage e dados específicos do projeto ficam na camada da aplicação (`app`), fora de `shared`
- contextos e providers reutilizáveis e agnósticos ao projeto podem viver em `shared`
- stores, contextos e providers ligados a autenticação, sessão, regras de negócio, rotas, tenant ou permissões devem ficar em `app`
- a pasta `shared/components/ui` é exceção: componentes originados do **shadcn/ui** podem manter a convenção original da biblioteca
- estas convenções podem evoluir quando houver necessidade explícita do time, desde que permaneçam consistentes dentro do contexto alterado

## Convenções para `apps/backend`

- `cmd/server` — entry point da aplicação
- `internal/` — código privado da aplicação (não importável externamente)
- handlers usam `chi` para roteamento
- services contêm regras de negócio e orquestram repositories
- repositories encapsulam acesso a dados (PostgreSQL, Redis, etc.)
- middleware para autenticação, logging, recovery e CORS
