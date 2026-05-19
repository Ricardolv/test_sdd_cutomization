---
name: backend-go-sync-module
description: Sincroniza entidades do dominio com o schema usando migrations SQL, sem ORM.
---

# Backend Go Sync Module

## Objetivo

Manter o schema do banco alinhado com as entidades do dominio de um modulo, usando migrations SQL versionadas.

## Entradas obrigatorias

- modulo alvo (nome ou path).

## Workflow

1. Ler entidades do modulo e identificar campos persistidos.
2. Criar migrations SQL em `apps/backend/migrations/`:
   - `YYYYMMDDHHMM_<modulo>_create.sql`
3. Evitar alterar migrations antigas; criar novas migrations incrementais.
4. Validar com a ferramenta de migrations adotada no projeto.

## Guardrails

- Nao executar sem modulo explicitamente informado.
- Nao criar schema automaticamente em producao.
- Tratar renames e drops como mudanca de risco.

## Saida esperada

- Migration criada para o modulo.
- Schema alinhado com o dominio.
