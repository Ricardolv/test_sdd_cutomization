# Contexto Técnico Global

## Stack base

- back-end em **Go** com **chi** para roteamento
- front-end web em **Next.js** (React, TypeScript)
- UI com **shadcn/ui**
- banco relacional com **PostgreSQL** (sem ORM, queries com `database/sql`)
- API REST com JSON
- domínio principal: módulo `customer` em `internal/modules/`

## Decisões já tomadas

- arquitetura simples, legível e incremental (sem abstrações antecipadas)
- estrutura Go: `cmd/server`, `internal/modules`, `config/`, `scripts/`
- cada módulo segue: `models/`, `services/`, `repositories/`, `handlers/`
- tratamento centralizado de erro com middleware e respostas padronizadas
- evolução do produto por mudanças pequenas e rastreáveis (uma spec por entrega)

## Restrições

- cada mudança deve caber em uma spec objetiva
- manter o projeto pequeno o suficiente para ensino em aula

## Padrões de integração

- front-end consome API REST
- validações simples acontecem no client e no server
- erros de domínio e validação são tratados de forma padronizada no backend
- respostas de erro da API devem ser exibidas de forma compreensível no front
