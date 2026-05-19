---
name: module-entity-go
description: Cria entidades de dominio Go com validacao explicita e testes unitarios completos.
---

# Module Entity Go

## Objetivo

Criar ou atualizar uma entidade de dominio em `modules/<modulo>/<aggregate>/model`, com validacao explicita e testes unitarios.

## Entradas obrigatorias

1. modulo
2. agregado
3. entidade
4. lista de atributos com tipos

## Workflow

1. Criar struct da entidade em `model/<entity>.go`.
2. Criar metodo `Validate()` que usa regras compartilhadas do `packages/shared`.
3. Criar testes em `model/<entity>_test.go` cobrindo:
   - criacao valida
   - invalidos por campo
4. Atualizar exports do pacote quando necessario.

## Guardrails

- Nao executar validacao no construtor.
- Nao criar persistencia ou handlers.

## Saida esperada

- Entidade criada com validacao.
- Testes unitarios completos.
