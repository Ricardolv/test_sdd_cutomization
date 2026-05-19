---
name: backend-go-provider-implementation
description: Implementa providers tecnicos do dominio no backend Go, mantendo contratos intactos.
---

# Backend Go Provider Implementation

## Objetivo

Criar implementacoes concretas de providers tecnicos definidos nos modulos de dominio, registrando-os na composicao do backend.

## Entradas obrigatorias

- interface do provider (nome ou path).

## Workflow

1. Ler a interface do provider no modulo.
2. Criar implementacao em `apps/backend/internal/provider/<modulo>/`.
3. Injetar configuracao via env quando necessario.
4. Registrar a implementacao no bootstrap do backend.
5. Criar testes unitarios quando houver logica observavel.

## Guardrails

- Nao modificar o contrato do dominio.
- Nao adicionar dependencias sem necessidade real.
- Preferir bibliotecas simples e estaveis.

## Saida esperada

- Implementacao concreta pronta para injecao.
- Dependencias ajustadas somente quando necessario.
