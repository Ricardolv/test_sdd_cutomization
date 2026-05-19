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

- [ ] Criar o agregado `customer` com struct em `models/customer_model.go`, campos: `ID`, `Name` (string, required), `Email` (string, required, unique), `Phone` (string, optional), `Notes` (string, optional), `Active` (bool, default true), `CreatedAt`, `UpdatedAt`.

- [ ] Definir erros de domínio em `models/customer_error.go`: `ErrNameRequired`, `ErrEmailRequired`, `ErrEmailDuplicate`, `ErrCustomerNotFound`.

- [ ] Definir a interface do repositório em `repositories/customer_repository.go` com métodos: `Create`, `Update`, `FindByID`, `FindByEmail`, `List`, `Deactivate`.

- [ ] Implementar o serviço em `services/customer_service.go` com `Save` (criar/atualizar) e `Deactivate`. Validações de domínio no serviço (nome/email obrigatórios, e-mail único).

- [ ] Implementar o repositório em `repositories/customer_repository.go` com queries SQL para PostgreSQL.

- [ ] Cobrir serviço e repositório com testes unitários em `services/customer_service_test.go` e `repositories/customer_repository_test.go`.

### Tasks - Back-end

- [ ] Criar DTOs em `handlers/customer_dto.go` para request (`CreateCustomerRequest`, `UpdateCustomerRequest`) e response (`CustomerResponse`).

- [ ] Criar handler em `handlers/customer_handler.go` com `Create`, `Update`, `GetByID`, `List`, `Deactivate`. Endpoints em `/customers`.

- [ ] Criar router em `handlers/customer_router.go` registrando rotas no `chi`:
  - `POST /customers` → Create
  - `GET /customers` → List
  - `GET /customers/{id}` → GetByID
  - `PUT /customers/{id}` → Update
  - `DELETE /customers/{id}` → Deactivate

- [ ] Registrar rotas de customer no `cmd/server/main.go`.

- [ ] Criar migration para a tabela `customers` (usar `migrate` ou `goose` em `scripts/`).

- [ ] Criar `scripts/customer.http` (Rest Client) cobrindo os fluxos do CRUD, incluindo os principais casos de erro. Validar manualmente com o backend rodando.

### Tasks - Validação

- [ ] `go test ./...` em `apps/backend/` — todos os testes passando.

- [ ] `go build ./cmd/server` — compilação sem erros.

- [ ] Backend inicia e responde em `localhost:9090`.

- [ ] Cenários em `scripts/customer.http` executam com sucesso (criar, listar, buscar, atualizar, inativar, e-mail duplicado, campos obrigatórios ausentes).

## Resultado Esperado

- Model `Customer` com struct validada, erros de domínio, interface de repositório e serviço implementados e testados.
- Tabela `customers` criada no PostgreSQL com migration aplicada.
- CRUD de `cliente` exposto no backend via handlers Go com chi, com cenários cobertos no `customer.http`.
- E-mail único validado no serviço antes de persistir.
- Inativação via soft-delete (campo `active`).

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
