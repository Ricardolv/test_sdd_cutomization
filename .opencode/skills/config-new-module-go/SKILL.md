---
name: config-new-module-go
description: Cria um novo modulo Go em modules/<modulo>, registra no go.work e conecta o backend.
---

# Config New Module Go

## Objetivo

Criar de forma deterministica um novo modulo de negocio em `modules/<modulo>` como modulo Go independente e registra-lo no `go.work`.

## Entradas obrigatorias

1. `nome do modulo` (kebab-case).
2. `modulePath` base do repositorio (ex.: `github.com/org/projeto`).

## Workflow

1. Validar que o root possui `go.work` ou cria-lo.
2. Criar `modules/<modulo>/` com:
   - `go mod init <modulePath>/modules/<modulo>`
   - estrutura base: `<aggregate>/model`, `<aggregate>/provider`, `<aggregate>/usecase` quando solicitado.
3. Adicionar o modulo ao `go.work` com `go work use modules/<modulo>`.
4. Conectar o backend:
   - garantir que `apps/backend` usa `go work` para enxergar o novo modulo.

## Guardrails

- Nao criar o modulo se `modules/<modulo>` ja existir.
- Nao prosseguir sem `modulePath`.
- Manter o modulo isolado e sem dependencia de infra do backend.

## Saida esperada

- `modules/<modulo>` criado com `go.mod`.
- `go.work` atualizado para incluir o novo modulo.
