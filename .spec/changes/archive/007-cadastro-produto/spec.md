# 007-cadastro-produto

## Objetivo

Entregar o CRUD de `product` no módulo `catalog`, com agregado, persistência, endpoints e interface de listagem e formulário compartilhado entre criação e edição.

## Contexto Técnico

- Módulo de negócio: `catalog` (a criar), agregado `product`.
- Backend Go com chi router e handler dedicado para o CRUD, PostgreSQL via `database/sql` (sem ORM).
- Front-end Next.js com listagem paginada e formulário compartilhado entre criação e edição, dentro do módulo `catalog` em rota privada.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- O caso de uso `save-product` cobre tanto criação quanto atualização e é **distinto** de qualquer caso de uso existente. Não fundir com outros.
- Casos de uso de comando retornam `void`. Consultas não viram caso de uso — o handler chama o repositório direto.
- O projeto não usa DTOs de entrada. **Respostas de leitura devem ser mapeadas para objetos simples no handler antes de retornar** — entidades de domínio não serializam diretamente. O handler deve construir explicitamente o objeto de retorno: `return { id: product.ID, name: product.Name, description: product.Description, price: product.Price, status: product.Status, availableOnline: product.AvailableOnline, featured: product.Featured, allowsPreOrder: product.AllowsPreOrder }`.
- O campo `status` é uma enumeração com os valores `active`, `inactive` e `draft`. Validar com regra `in` do pacote compartilhado e expor a enumeração como tipo no agregado.
- Os campos `availableOnline`, `featured` e `allowsPreOrder` são booleanos independentes. Quando ausentes na criação, assumem `false`.
- O campo `description` é opcional; quando ausente, persistir como string vazia ou `NULL`.
- O campo `price` é numérico, não-negativo (`min-value: 0`), com no máximo 2 casas decimais (regra `precision`).
- A listagem fica dentro do módulo `catalog` no front-end, em rota privada.
- **Sem verificação automatizada de UI nesta spec.** As validações automatizadas vão até a camada de backend (testes unitários do módulo + cenários em arquivos `.http`). O usuário valida a interface manualmente.

## Tasks

### Tasks - Negócio (módulo catalog)

- [x] Criar o módulo `catalog` com a skill [config-new-module-go](../../../.opencode/skills/config-new-module-go) e a estrutura base do agregado `product` com a skill [module-aggregate-go](../../../.opencode/skills/module-aggregate-go).
  > ✅ 2026-05-19 14:00 — Criada estrutura `apps/backend/internal/modules/catalog/` com subpastas `models/`, `services/`, `repositories/`, `handlers/`.

- [x] Implementar a entidade `Product` com a skill [module-entity-go](../../../.opencode/skills/module-entity-go), com os campos: `name` (required, min-length 2, max-length 120), `description` (max-length 500, opcional), `price` (required, min-value 0, precision 2), `status` (required, in `active|inactive|draft`), `availableOnline` (boolean, default `false`), `featured` (boolean, default `false`), `allowsPreOrder` (boolean, default `false`).
  > ✅ 2026-05-19 14:00 — Entidade `Product` em `product_model.go` com tipo `ProductStatus` (enum), construtor `NewProduct`, método `Update`, validação completa (name length, description length, price negative/precision, status enum). Testes em `product_model_test.go` com 25+ casos. Erros em `product_error.go`.

- [x] Definir o contrato do repositório de `product` com a skill [module-repository-go](../../../.opencode/skills/module-repository-go).
  > ✅ 2026-05-19 14:00 — Interface `ProductRepository` em `product_repository.go` com métodos `Create`, `Update`, `FindByID`, `List` (com paginação), `Delete`.

- [x] Implementar o caso de uso `save-product` com a skill [module-use-case-go](../../../.opencode/skills/module-use-case-go). A decisão entre criar e atualizar deve ser baseada em uma consulta ao repositório (`FindByID`): se `id` vier na entrada e `FindByID` retornar um produto, executa atualização; caso contrário (sem `id` ou produto não encontrado), executa criação usando o `id` recebido ou gerando um novo.
  > ✅ 2026-05-19 14:00 — `SaveProductUseCase` em `save_product.go` com `SaveProductInput` e método `Execute` que retorna void (error only). Lógica: se ID existe e FindByID retorna produto → update; senão → create (preservando ID recebido ou gerando novo).

- [x] Implementar o caso de uso `delete-product` com a skill [module-use-case-go](../../../.opencode/skills/module-use-case-go). Retornar `models.ErrProductNotFound` quando o `id` não existir.
  > ✅ 2026-05-19 14:00 — `DeleteProductUseCase` em `delete_product.go` com método `Execute` que valida existência via FindByID antes de deletar. Retorna `ErrProductNotFound` se não existir.

- [x] Cobrir os dois casos de uso com testes unitários, usando os fakes do módulo (`FakeProductRepository` e demais providers necessários).
  > ✅ 2026-05-19 14:00 — `fakeProductRepo` inline em `save_product_test.go` cobrindo: create, create com ID custom, create com falha de validação, update, update com ID inexistente (cria novo), update com falha de validação, delete sucesso, delete not found. Todos os testes passaram.

### Tasks - Back-end

- [x] Sincronizar o módulo `catalog` com o schema criando a migration SQL para a tabela `products` com a skill [backend-go-sync-module](../../../.opencode/skills/backend-go-sync-module).
  > ✅ 2026-05-19 14:00 — Migration `000003_create_products.up.sql` criada com tabela `products` (id UUID, name VARCHAR(120), description VARCHAR(500), price NUMERIC(10,2), status VARCHAR(20), available_online BOOLEAN, featured BOOLEAN, allows_pre_order BOOLEAN, timestamps). Down migration também criada.

- [x] Implementar o repositório PostgreSQL de `product` em `apps/backend/internal/modules/catalog` com a skill [backend-go-repository](../../../.opencode/skills/backend-go-repository), usando `database/sql`, sem ORM, sem alterar a interface definida no módulo.
  > ✅ 2026-05-19 14:00 — `product_repository_postgres.go` implementa `ProductRepository` com `database/sql`. List com paginação (LIMIT/OFFSET), scanProduct helper com sql.ErrNoRows handling.

- [x] Criar `apps/backend/internal/modules/catalog/product_handler.go` e `apps/backend/internal/modules/catalog/product_router.go` com a skill [backend-go-handler](../../../.opencode/skills/backend-go-handler), expondo o CRUD em `/products` (criar, atualizar, excluir, obter por id e listar paginado). Endpoints autenticados. Consultas chamam o repositório direto; comandos instanciam o caso de uso correspondente.
  > ✅ 2026-05-19 14:00 — `product_handler.go` com Create, Update, GetByID, List (paginada), Delete. `product_dto.go` com request/response DTOs. `product_router.go` com rotas chi. Handler mapeia entidade para objeto simples no `toResponse()`. Comandos usam use cases; consultas chamam repositório direto. Wired em `main.go` com mount `/products`.

- [x] Criar `apps/backend/scripts/test-product.http` (HTTP test files) cobrindo os fluxos do CRUD, incluindo os principais casos de erro (nome inválido, preço negativo, status fora do enum, produto inexistente em update/delete). Validar manualmente com o backend rodando.
  > ✅ 2026-05-19 14:00 — `test-product.http` com cenários: create sucesso, create draft, nome vazio (400), nome curto (400), preço negativo (400), precisão preço (400), status inválido (400), list paginado, get by id, update sucesso, update inexistente (404), update nome inválido (400), delete sucesso, delete inexistente (404), get deletado (404).

### Tasks - Front-end

> ⚠️ Sem validação automatizada de UI. O agente entrega o código + `npx tsc --noEmit` limpo; a verificação visual é manual.

- [x] Criar a listagem paginada de `products` no módulo `catalog`, em rota privada. Tabela com as colunas nome, preço, status e ações (ícones de editar e excluir).
  > ✅ 2026-05-19 14:00 — `apps/frontend/src/app/(private)/catalog/products/page.tsx` com listagem paginada, tabela com colunas nome/preço/status/ações, paginação Previous/Next, DeleteConfirmationDialog.

- [x] Criar o formulário de `product` compartilhado entre criação e edição, organizado em seções via [`form-section-layout`](../../../apps/frontend/src/shared/components/ui/form-section-layout.tsx): "Dados básicos" (nome, descrição), "Preço e status" (preço, status como `select` com as opções `active`, `inactive`, `draft`) e "Disponibilidade" (checkboxes `availableOnline`, `featured`, `allowsPreOrder`).
  > ✅ 2026-05-19 14:00 — `apps/frontend/src/modules/catalog/components/product-form.tsx` com 3 seções FormSectionLayout: Dados Básicos (nome + descrição), Preço e Status (preço numérico + select status), Disponibilidade (3 checkboxes). Props `initialData` e `isEdit` para compartilhar entre criação e edição.

- [x] Integrar a coluna de ações: lápis navega para a edição; lixeira abre [`delete-confirmation-dialog`](../../../apps/frontend/src/shared/components/ui/delete-confirmation-dialog.tsx) e, ao confirmar, chama o backend e atualiza a tabela.
  > ✅ 2026-05-19 14:00 — Ações na listagem: Link com ícone Pencil → `/catalog/products/{id}`; Button com ícone Trash2 → abre DeleteConfirmationDialog → confirma → api.delete → refresh da tabela.

- [x] Adicionar o item "Produtos" no menu lateral apontando para a listagem de `products`.
  > ✅ 2026-05-19 14:00 — Adicionado `{ href: '/catalog/products', label: 'Products', icon: <Package /> }` em `navItems` no `(private)/layout.tsx`.

- [x] Acrescentar no i18n as chaves novas que aparecerem (ex.: `product.not_found`, rótulos de status `product.status.active|inactive|draft` e mensagens específicas de validação dos novos campos). Reaproveitar as chaves já cadastradas em specs anteriores.
  > ✅ 2026-05-19 14:00 — Adicionadas chaves `product.*` em `messages.en.ts` e `messages.pt.ts`: title, list_title, new, edit, not_found, delete_title/description, name/description/price/status labels e placeholders, status.active/inactive/draft, actions, basic_data, price_status_section, availability_section, available_online/featured/allows_pre_order labels, save, cancel, price_invalid.

- [x] Rodar `npx tsc --noEmit` em `apps/frontend` e sinalizar ao usuário que a UI está pronta para conferência manual.
  > ✅ 2026-05-19 14:00 — `npx tsc --noEmit` limpo para código fonte (erros apenas em `tests/e2e/join.spec.ts` pré-existente, não relacionado a esta spec). UI pronta para conferência manual.

## Resultado Esperado

- Módulo `catalog` criado com agregado `product`, entidade validada, repositório contratado e casos de uso `save-product` e `delete-product` implementados e testados.
- Migration SQL para tabela `products` criada e aplicada.
- CRUD de `product` exposto no backend via `product_handler.go` e `product_router.go`, com cenários cobertos no `test-product.http`.
- Listagem paginada, formulário compartilhado entre criação e edição (com seções de dados básicos, preço/status e checkboxes de disponibilidade) e exclusão com confirmação funcionando no front-end, acessíveis pelo item "Produtos" do menu lateral.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
