---
name: config-go-database
description: Configura persistencia no backend Go usando database/sql e driver PostgreSQL, sem ORM.
---

# Config Go Database

## Objetivo

Padronizar a infraestrutura de persistencia do backend Go usando `database/sql`, sem ORM, com configuracao via variaveis de ambiente.

## Workflow

1. Garantir que `apps/backend` existe e possui `go.mod`.
2. Adicionar driver PostgreSQL:
   - `go get github.com/jackc/pgx/v5/stdlib`
3. Criar pacote `internal/db` com:
   - `Open(ctx, dsn) (*sql.DB, error)`
   - `PingContext` e configuracao de pool (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`).
4. Configurar variaveis:
   - `DATABASE_URL` (ou `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`).
5. Expor healthcheck simples do banco no backend.
6. Se houver migrations, manter em `apps/backend/migrations/` e exigir ferramenta definida pelo projeto.

## Guardrails

- Nao usar ORM ou query builder pesado.
- Nao embutir segredos no repositorio.
- Nao abrir conexoes sem fechamento correto no shutdown.

## Saida esperada

- Driver instalado.
- Pacote `internal/db` pronto para injecao.
- Configuracao de banco via env estabelecida.
