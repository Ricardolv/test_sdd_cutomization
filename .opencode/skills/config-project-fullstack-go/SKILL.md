---
name: config-project-fullstack-go
description: Cria um monorepo fullstack com apps/frontend (Next.js, Angular ou SvelteKit) e apps/backend (Go + chi), usando a versao mais recente do Go instalada, sem ORM.
---

# Config Project Fullstack Go

## Objetivo

Criar a estrutura base de um monorepo com `apps/frontend` (Next.js, Angular ou SvelteKit) e `apps/backend` (Go + chi), com `go.work` na raiz, configuracao HTTP basica, CORS e variaveis de ambiente.

## Workflow

1. Executar sempre a partir da pasta que sera a raiz do projeto.
2. Garantir as pastas base:
   - `apps/`
   - `modules/`
   - `packages/`
3. Frontend (escolha obrigatoria):
   - Se o usuario nao especificar, interromper e perguntar explicitamente qual frontend deseja.
   - Opcao Next.js:
     - `npx create-next-app@latest apps/frontend --yes --src-dir --use-npm`
   - Opcao Angular:
     - `mkdir -p apps && (cd apps && npx -y @angular/cli@latest new frontend)`
   - Opcao SvelteKit:
     - `npx sv create apps/frontend`
4. Backend Go:
   - criar `apps/backend/` com `cmd/api/main.go` e `internal/` para handlers.
   - inicializar modulo:
     - `go mod init <modulePath>/apps/backend`
   - adicionar chi:
     - `go get github.com/go-chi/chi/v5`
5. Criar `go.work` na raiz com:
   - `apps/backend`
   - `modules/*`
   - `packages/*`
6. Router chi base em `cmd/api/main.go`:
   - `chi.NewRouter()`
   - middlewares: `RequestID`, `RealIP`, `Logger`, `Recoverer`
   - `http.ListenAndServe(":4000", router)`
7. Configurar o frontend com a URL do backend:
   - Next.js:
     - criar `.env.example` e `.env` com `NEXT_PUBLIC_API_URL=http://localhost:4000`
   - Angular:
     - ajustar `apps/frontend/src/environments/environment.ts` com `apiUrl: "http://localhost:4000"`
     - quando existir `environment.development.ts`, espelhar o mesmo valor
   - SvelteKit:
     - criar `.env.example` e `.env` com `PUBLIC_API_URL=http://localhost:4000`
8. Validar:
   - `go test ./...` em `apps/backend`
   - `npm --prefix apps/frontend run build`

## Guardrails

- Usar a versao mais recente do Go instalada (nao fixar versao antiga no `go.mod`).
- Usar `chi` como router principal.
- Nao usar ORM; somente `database/sql` quando houver persistencia.
- Nao escolher framework de frontend sem confirmacao do usuario.

## Saida esperada

- `apps/frontend` criado com o framework escolhido.
- `apps/backend` criado com Go + chi.
- `go.work` configurado com apps, modules e packages.
- Variaveis de ambiente do frontend prontas para consumir o backend.
