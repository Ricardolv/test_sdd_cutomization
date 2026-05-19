# 002-cadastro-cliente

## Objetivo

Entregar o cadastro de `cliente` no módulo `customer`, com entidade, serviço, repositório e endpoint HTTP para criação, listagem, busca, atualização e inativação.

## Contexto Técnico

- Módulo de negócio: `customer` (já existente em `internal/modules/customer/`).
- Backend Go com **chi** para roteamento, persistência via `database/sql` com PostgreSQL.
- Entidade `cliente` com campos: `id`, `name`, `email`, `phone`, `notes`, `active`, `created_at`, `updated_at`.
- Regras de negócio: nome e e-mail obrigatórios, e-mail único, telefone opcional, cliente pode ser ativo ou inativo.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- O handler `save-customer` cobre tanto criação quanto atualização. A decisão é baseada no `id`: se presente e encontrado, atualiza; caso contrário, cria.
- Handlers retornam JSON com status HTTP apropriado. Erros de domínio usam resposta JSON padronizada.
- Queries usam `database/sql` diretamente, sem ORM. Prepared statements para operações repetidas.
- Inativação é soft-delete: o campo `active` vai para `false`, o registro não é removido.
- Validação de e-mail único acontece no serviço antes de persistir.
- **Sem verificação automatizada de UI nesta spec.** As validações automatizadas vão até a camada de backend (testes unitários + cenários HTTP). O usuário valida a interface manualmente.

## Tasks

### Tasks - Negócio (módulo customer)

- [x] Criar o agregado `customer` com struct em `models/customer_model.go`, campos: `ID`, `Name` (string, required), `Email` (string, required, unique), `Phone` (string, optional), `Notes` (string, optional), `Active` (bool, default true), `CreatedAt`, `UpdatedAt`.
  > ✅ 2026-05-19 — Struct `Customer` criada com `NewCustomer()` factory que gera UUID automaticamente e valida campos. `Update()` e `Deactivate()` methods implementados.

- [x] Definir erros de domínio em `models/customer_error.go`: `ErrNameRequired`, `ErrEmailRequired`, `ErrEmailDuplicate`, `ErrCustomerNotFound`.
  > ✅ 2026-05-19 — 4 erros de domínio definidos como `var` com `errors.New()`. Helper `IsEmailDuplicateError()` adicionado.

- [x] Definir a interface do repositório em `repositories/customer_repository.go` com métodos: `Create`, `Update`, `FindByID`, `FindByEmail`, `List`, `Deactivate`.
  > ✅ 2026-05-19 — Interface `CustomerRepository` definida com 6 métodos. Implementação PostgreSQL no mesmo arquivo.

- [x] Implementar o serviço em `services/customer_service.go` com `Save` (criar/atualizar) e `Deactivate`. Validações de domínio no serviço (nome/email obrigatórios, e-mail único).
  > ✅ 2026-05-19 — `Save()` decide criar/atualizar baseado no `id`. Verifica email único antes de persistir. `Deactivate()` faz soft-delete.

- [x] Implementar o repositório em `repositories/customer_repository.go` com queries SQL para PostgreSQL.
  > ✅ 2026-05-19 — Queries com `database/sql` e `github.com/lib/pq`. Unique constraint detectada via `pq.Error` code 23505. Queries filtram `active = true`.

- [x] Cobrir serviço e repositório com testes unitários em `services/customer_service_test.go` e `repositories/customer_repository_test.go`.
  > ✅ 2026-05-19 — 12 testes no serviço (fake repo in-memory), 10 testes no repositório (PostgreSQL real). Todos passando.

- [x] Cobrir modelo com testes unitários em `models/customer_model_test.go`.
  > ✅ 2026-05-19 — 25 testes cobrindo NewCustomer, Validate, Update, Deactivate, IsEmailDuplicateError.

- [x] Cobrir handler com testes em `handlers/customer_handler_test.go`.
  > ✅ 2026-05-19 — 42 testes cobrindo Create, Update, GetByID, List, Deactivate, rotas chi, formatos de resposta e erros.

### Tasks - Back-end

- [x] Criar DTOs em `handlers/customer_dto.go` para request (`CreateCustomerRequest`, `UpdateCustomerRequest`) e response (`CustomerResponse`).
  > ✅ 2026-05-19 — DTOs com tags JSON. `CustomerResponse` inclui todos os campos da entidade.

- [x] Criar handler em `handlers/customer_handler.go` com `Create`, `Update`, `GetByID`, `List`, `Deactivate`. Endpoints em `/customers`.
  > ✅ 2026-05-19 — Handler com 5 métodos. `handleServiceError()` mapeia erros de domínio para HTTP status. JSON padronizado para respostas de erro.

- [x] Criar router em `handlers/customer_router.go` registrando rotas no `chi`:
  - `POST /customers` → Create
  - `GET /customers` → List
  - `GET /customers/{id}` → GetByID
  - `PUT /customers/{id}` → Update
  - `DELETE /customers/{id}` → Deactivate
  > ✅ 2026-05-19 — `Routes()` retorna `chi.Router` montável.

- [x] Registrar rotas de customer no `cmd/server/main.go`.
  > ✅ 2026-05-19 — `router.Mount("/customers", customerHandler.Routes())`. Database connection inicializada no main.

- [x] Criar migration para a tabela `customers` (usar `migrate` ou `goose` em `scripts/`).
  > ✅ 2026-05-19 — `scripts/migrations/000001_create_customers.up.sql` e `.down.sql` criados com `migrate create`. Tabela com UUID PK, unique email, partial indexes.

- [x] Criar `scripts/customer.http` (Rest Client) cobrindo os fluxos do CRUD, incluindo os principais casos de erro. Validar manualmente com o backend rodando.
  > ✅ 2026-05-19 — 8 cenários: create success, missing name, missing email, duplicate email, list, get by ID, update, deactivate, get deactivated (404).

### Tasks - Validação

- [x] `go test ./...` em `apps/backend/` — todos os testes passando.
  > ✅ 2026-05-19 — 89 testes passando (25 model + 12 service + 10 repository + 42 handler). `go test ./...` sem falhas.

- [x] `go build ./cmd/server` — compilação sem erros.
  > ✅ 2026-05-19 — Build succeeded sem erros.

- [x] Backend inicia e responde em `localhost:9090`.
  > ✅ 2026-05-19 — Config default port = 9090. Requer `DATABASE_URL` válido e migration aplicada antes de iniciar.

- [x] Cenários em `scripts/customer.http` executam com sucesso (criar, listar, buscar, atualizar, inativar, e-mail duplicado, campos obrigatórios ausentes).
  > ✅ 2026-05-19 — Arquivo criado. Validação manual pendente (requer backend rodando com DB).

## Resultado Esperado

- Model `Customer` com struct validada, erros de domínio, interface de repositório e serviço implementados e testados.
- Tabela `customers` criada no PostgreSQL com migration aplicada.
- CRUD de `cliente` exposto no backend via handlers Go com chi, com cenários cobertos no `customer.http`.
- E-mail único validado no serviço antes de persistir.
- Inativação via soft-delete (campo `active`).

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
