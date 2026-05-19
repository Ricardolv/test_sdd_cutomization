---
name: module-repository-go
description: Cria contrato de repositorio de dominio e fake/in-memory para testes em Go.
---

# Module Repository Go

## Objetivo

Criar a interface de repositorio do agregado e uma implementacao fake/in-memory para testes de casos de uso.

## Entradas obrigatorias

1. modulo
2. agregado
3. tipo de repositorio (`crud` ou `custom`)
4. entidade principal

## Workflow

1. Criar interface em `provider/<entity>_repository.go`.
2. Criar fake em `provider/fake_<entity>_repository_test.go` ou `test/fake_<entity>_repository.go`.
3. Implementar armazenamento em memoria com `map` ou `slice`.
4. Garantir exports do pacote.

## Guardrails

- Nao criar implementacao real de banco.
- Nao alterar entidades.
- Nao usar mocks de framework como base principal.

## Saida esperada

- Interface de repositorio criada.
- Fake/in-memory funcional para testes.
