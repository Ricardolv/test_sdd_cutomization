---
name: backend-go-config
description: Configura base compartilhada do backend Go com middlewares, erros, autenticacao JWT e helpers comuns.
---

# Backend Go Config

## Objetivo

Centralizar middlewares, tratamento de erro e autenticacao no backend Go (chi), evitando repeticao nos handlers.

## Workflow

1. Criar `internal/http/middleware` com:
   - `RequestID`, `RealIP`, `Logger`, `Recoverer` (via chi middleware).
   - middleware de `Content-Type` quando necessario.
2. Criar `internal/http/response` com helpers de erro padronizado.
3. Criar `internal/auth` com:
   - validacao de JWT (preferir `github.com/golang-jwt/jwt/v5`).
   - middleware `AuthRequired`.
4. Definir padrao de erro HTTP:
   - mapear `DomainError`, `ValidationError` e erros inesperados.
5. Registrar middlewares na composicao do router principal.

## Guardrails

- Nao duplicar `try/catch` manual nos handlers; use resposta centralizada.
- Nao expor stack trace em respostas.
- Nao guardar segredo JWT no repositorio.

## Saida esperada

- Middlewares e helpers compartilhados prontos.
- Padrao de erro consistente.
- Middleware de autenticacao reutilizavel.
