---
name: module-use-case-go
description: Cria casos de uso Go com contratos de entrada/saida e testes unitarios completos.
---

# Module Use Case Go

## Objetivo

Criar um caso de uso no dominio com contratos de entrada/saida, implementacao simples e testes com fakes concretas.

## Entradas obrigatorias

1. modulo
2. agregado
3. nome do caso de uso
4. tipo de cenario (`crud` ou `custom`)
5. se retorna saida ou `void`

## Workflow

1. Criar arquivo em `usecase/<case>.go`.
2. Definir structs `In` e `Out` quando necessario.
3. Implementar a funcao `Execute(ctx, in)` em uma struct `UseCase`.
4. Criar teste em `usecase/<case>_test.go` com fakes do modulo.

## Guardrails

- Nao criar handlers ou persistencia real.
- Nao inventar regra de negocio nao pedida.

## Saida esperada

- Caso de uso criado com contratos claros.
- Testes cobrindo sucesso e falhas.
