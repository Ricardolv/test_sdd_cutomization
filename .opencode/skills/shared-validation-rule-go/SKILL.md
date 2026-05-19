---
name: shared-validation-rule-go
description: Cria regras de validacao reutilizaveis no pacote shared para entidades do dominio em Go.
---

# Shared Validation Rule Go

## Objetivo

Adicionar uma regra de validacao reutilizavel em `packages/shared/validation`.

## Entradas obrigatorias

1. nome da regra
2. objetivo da validacao
3. tipo de valor validado
4. mensagem ou codigo de erro

## Workflow

1. Criar arquivo em `packages/shared/validation/rules/`.
2. Implementar interface `Rule` do shared.
3. Atualizar export do pacote de validacao.
4. Criar teste unitario para a regra.

## Guardrails

- Nao criar regras hiper especificas de um unico modulo.
- Reutilizar helpers do shared quando existirem.

## Saida esperada

- Regra criada com teste.
- Exports atualizados no shared.
