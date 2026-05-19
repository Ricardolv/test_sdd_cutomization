---
name: config-package-shared-go
description: Prepara packages/shared como modulo Go com erros, validacao e base de usecases reutilizaveis.
---

# Config Package Shared Go

## Objetivo

Criar o pacote compartilhado em `packages/shared` como modulo Go reutilizavel por `modules/*` e `apps/backend`.

## Workflow

1. Criar `packages/shared` e iniciar o modulo:
   - `go mod init <modulePath>/packages/shared`
2. Estrutura base sugerida:
   - `error/` (`DomainError`, `ValidationError`)
   - `model/` (`Entity`, `EntityState`, `PageResult`)
   - `usecase/` (`UseCase` interface)
   - `validation/` (`Rule`, `Validator`, helpers)
3. Adicionar `packages/shared` ao `go.work`.
4. Verificar `go test ./...` no pacote shared.

## Guardrails

- Nao adicionar dependencias do backend no shared.
- Manter contratos reutilizaveis e leves.

## Saida esperada

- `packages/shared` criado como modulo Go.
- Contratos basicos de erro, entidade, usecase e validacao prontos.
