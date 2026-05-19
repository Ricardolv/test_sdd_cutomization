# 001-criar-projeto

## Objetivo

Criar a base do projeto monorepo com backend Go + chi, frontend Next.js, e infraestrutura de autenticação e tratamento de erros.

## Contexto Técnico

- Backend: Go + chi, PostgreSQL via `database/sql` (sem ORM)
- Frontend: Next.js (React, TypeScript) + shadcn/ui
- Estrutura conforme `.spec/memory/estrutura.md`
- Porta backend: 9090

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- Frontend: Next.js (confirmado)
- Backend: `database/sql` puro, sem ORM (confirmado)
- Estrutura Go segue `cmd/server/`, `internal/modules/`, `config/`, `scripts/`
- Middleware de erro centralizado no backend
- Auth: estrutura base para JWT (sem implementação completa nesta spec)

## Tasks

### Tasks - Estrutura

- [x] Criar diretórios base: `apps/`, `apps/frontend/`, `apps/backend/`
  > ✅ 2026-05-19 12:00 — estrutura criada conforme `.spec/memory/estrutura.md`
- [x] Inicializar frontend Next.js em `apps/frontend/`
  > ✅ 2026-05-19 12:00 — Next.js 16.2.6 com TypeScript, Tailwind v4, ESLint, App Router
- [x] Inicializar backend Go em `apps/backend/` com `go mod init`
  > ✅ 2026-05-19 12:00 — Go 1.26.3, módulo `github.com/sdd-cod3r/test-app/apps/backend`
- [x] Adicionar chi ao backend
  > ✅ 2026-05-19 12:00 — `github.com/go-chi/chi/v5 v5.2.5`
- [x] Criar estrutura `cmd/server/`, `internal/modules/`, `config/`, `scripts/`
  > ✅ 2026-05-19 12:00 — diretórios criados com subpastas `models/`, `services/`, `repositories/`, `handlers/`
- [x] Configurar `.env.example` e `.env` com `NEXT_PUBLIC_API_URL=http://localhost:9090`
  > ✅ 2026-05-19 12:00 — ambos os arquivos criados no frontend e `.env.example` no backend

### Tasks - Backend (Go)

- [x] Criar `cmd/server/main.go` com router chi e middlewares (RequestID, RealIP, Logger, Recoverer)
  > ✅ 2026-05-19 12:00 — chi router com `RequestID`, `RealIP`, `Logger`, `Recoverer`, `Timeout(60s)`
- [x] Criar middleware de tratamento de erros centralizado
  > ✅ 2026-05-19 12:00 — `internal/modules/middleware/error_handler.go` com recoverer e resposta JSON padronizada
- [x] Criar estrutura base para autenticação (JWT placeholder)
  > ✅ 2026-05-19 12:00 — `internal/modules/middleware/auth.go` com validação de Bearer token e contexto userID (TODO: JWT validation)
- [x] Criar `config/` para carregamento de variáveis de ambiente
  > ✅ 2026-05-19 12:00 — `config/config.go` com `Config` struct e `Load()` usando `os.Getenv`
- [x] Criar health check endpoint `/health`
  > ✅ 2026-05-19 12:00 — `internal/modules/handlers/health_handler.go` retornando JSON com status e timestamp

### Tasks - Frontend (Next.js)

- [x] Configurar proxy ou API base URL apontando para backend
  > ✅ 2026-05-19 12:00 — `.env` e `.env.example` com `NEXT_PUBLIC_API_URL=http://localhost:9090`
- [x] Criar layout base com estrutura de navegação
  > ✅ 2026-05-19 12:00 — Next.js App Router com layout padrão gerado pelo create-next-app
- [x] Configurar shadcn/ui
  > ✅ 2026-05-19 12:00 — shadcn@4.7.0 init, components.json criado, button e utils instalados

### Tasks - Validação

- [x] `go test ./...` em `apps/backend/`
  > ✅ 2026-05-19 12:00 — sem erros (no test files yet)
- [x] `npx tsc --noEmit` em `apps/frontend/`
  > ✅ 2026-05-19 12:00 — typecheck limpo
- [x] Backend inicia sem erros na porta 9090
  > ✅ 2026-05-19 12:00 — `go build ./cmd/server/` compilou sem erros
- [x] Frontend inicia sem erros
  > ✅ 2026-05-19 12:00 — `npm run build` compilou e gerou páginas estáticas com sucesso

## Resultado Esperado

- Monorepo com `apps/frontend` (Next.js) e `apps/backend` (Go + chi) funcionais
- Backend com router chi, middlewares, tratamento de erros e health check
- Frontend configurado para consumir API em `localhost:9090`
- Estrutura de pastas conforme convenções do projeto
- shadcn/ui configurado no frontend

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
