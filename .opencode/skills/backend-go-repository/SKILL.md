---
name: backend-go-repository
description: Implementa repositorio de dominio no backend Go usando database/sql, sem ORM.
---

# Backend Go Repository

## Objetivo

Implementar uma interface de repositorio do dominio usando `database/sql`, com consultas SQL explicitas.

## Entradas obrigatorias

- interface do repositorio (nome ou path explicito).

## Workflow

1. Ler a interface alvo no modulo de dominio.
2. Criar implementacao em `apps/backend/internal/repository/<modulo>/`.
3. Usar `*sql.DB` ou `*sql.Tx` via construtor.
4. Implementar metodos com SQL explicito:
   - `create`, `update`, `delete`, `findById`, `findPage` quando existirem.
5. Mapear rows manualmente para entidades.

## Guardrails

- Nao alterar a interface do dominio.
- Nao usar ORM ou query builder pesado.
- Nao hardcode de SQL sem placeholders.

## Saida esperada

- Implementacao concreta criada.
- Repositorio pronto para injecao no backend.
