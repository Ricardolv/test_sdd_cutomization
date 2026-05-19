---
name: module-aggregate-go
description: Cria a estrutura base de um agregado no modulo Go, com pacotes de model, provider e usecase.
---

# Module Aggregate Go

## Objetivo

Criar a estrutura base de um agregado em `modules/<modulo>/<aggregate>/` com pacotes `model`, `provider` e `usecase`.

## Entradas obrigatorias

1. modulo
2. agregado
3. modo de usecase: `crud` ou `example`

## Workflow

1. Validar que `modules/<modulo>` existe.
2. Criar pacotes:
   - `modules/<modulo>/<aggregate>/model`
   - `modules/<modulo>/<aggregate>/provider`
   - `modules/<modulo>/<aggregate>/usecase`
3. Criar arquivos base conforme o modo.
4. Atualizar `doc.go` ou `package` quando o modulo tiver convencoes locais.

## Guardrails

- Nao criar handlers ou persistencia real.
- Nao inventar regras de negocio.

## Saida esperada

- Estrutura do agregado criada com pacotes basicos.
