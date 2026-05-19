# NNN-cadastro-{{entidade}}

> Template de cadastro (CRUD) de uma entidade qualquer dentro de um módulo de negócio **existente**.
> Antes de usar, substitua os placeholders abaixo e remova esta seção de instruções.
>
> **Placeholders:**
> - `{{entidade}}` — nome da entidade no singular, kebab-case (ex.: `customer`, `produto`, `pedido`).
> - `{{Entidade}}` — mesma entidade em PascalCase (ex.: `Customer`, `Produto`, `Pedido`).
> - `{{entidades}}` — plural em kebab-case usado em rotas/URLs (ex.: `customers`, `produtos`, `pedidos`).
> - `{{modulo}}` — módulo de negócio onde o cadastro vive (ex.: `auth`, `catalog`, `vendas`). **Assume-se que o módulo já existe**; este template não cria módulos novos.
> - `{{campos}}` — lista dos campos da entidade com tipo e validação (ex.: `name string` (required), `email string` (email format)).
> - `{{colunas-listagem}}` — colunas exibidas na tabela da listagem (ex.: nome, e-mail).
> - `{{secoes-formulario}}` — agrupamento de campos em seções do formulário (ex.: "Dados básicos" (nome, e-mail)).
> - `{{rotulo-menu}}` — rótulo do item adicionado ao menu lateral (ex.: "Clientes", "Produtos").

## Objetivo

Entregar o CRUD de `{{entidade}}` no módulo `{{modulo}}`, com model, persistência, endpoints e interface de listagem e formulário compartilhado entre criação e edição.

## Contexto Técnico

- Módulo de negócio: `{{modulo}}` (já existente), agregado `{{entidade}}`.
- Backend Go com **chi** para roteamento, handlers em `internal/modules/{{modulo}}/handlers/`, persistência via `database/sql` com PostgreSQL.
- Front-end Next.js com listagem paginada e formulário compartilhado entre criação e edição, dentro do módulo `{{modulo}}` em rota privada.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- O handler `save-{{entidade}}` cobre tanto criação quanto atualização. A decisão é baseada no `id`: se presente e encontrado, atualiza; caso contrário, cria.
- Handlers retornam JSON com status HTTP apropriado. Erros de domínio usam `http.Error` com mensagem JSON padronizada.
- Queries usam `database/sql` diretamente, sem ORM. Prepared statements para operações repetidas.
- A listagem fica dentro do módulo `{{modulo}}` no front-end, em rota privada.
- **Sem verificação automatizada de UI nesta spec.** As validações automatizadas vão até a camada de backend (testes unitários + cenários HTTP). O usuário valida a interface manualmente.

## Tasks

### Tasks - Negócio (módulo {{modulo}})

- [ ] Criar o agregado `{{entidade}}` dentro do módulo `{{modulo}}` com struct em `models/{{entidade}}_model.go`, campos: {{campos}}.

- [ ] Definir a interface do repositório em `repositories/{{entidade}}_repository.go` com métodos: `Create`, `Update`, `Delete`, `FindByID`, `List`.

- [ ] Implementar o serviço em `services/{{entidade}}_service.go` com `Save` (criar/atualizar) e `Delete`. Validações de domínio no serviço.

- [ ] Implementar o repositório em `repositories/{{entidade}}_repository.go` com queries SQL para PostgreSQL.

- [ ] Cobrir serviço e repositório com testes unitários em `services/{{entidade}}_service_test.go` e `repositories/{{entidade}}_repository_test.go`.

### Tasks - Back-end

- [ ] Criar handlers em `handlers/{{entidade}}_handler.go` com `Create`, `Update`, `Delete`, `GetByID`, `List`. Endpoints em `/{{entidades}}`. Handlers autenticados.

- [ ] Criar router em `handlers/{{entidade}}_router.go` registrando rotas no `chi` com middleware de autenticação.

- [ ] Criar `scripts/{{entidade}}.http` (Rest Client) cobrindo os fluxos do CRUD, incluindo os principais casos de erro. Validar manualmente com o backend rodando.

- [ ] Criar migration para a tabela `{{entidades}}` (usar `migrate` ou `goose`).

### Tasks - Front-end

> ⚠️ Sem validação automatizada de UI. O agente entrega o código + `npx tsc --noEmit` limpo; a verificação visual é manual.

- [ ] Criar a listagem paginada de `{{entidades}}` no módulo `{{modulo}}`, em rota privada. Tabela com as colunas {{colunas-listagem}} e ações (ícones de editar e excluir).

- [ ] Criar o formulário de `{{entidade}}` compartilhado entre criação e edição, organizado em seções: {{secoes-formulario}}.

- [ ] Criar funções de chamada à API em `src/shared/api/{{entidade}}.api.ts` (create, update, delete, getById, list).

- [ ] Integrar a coluna de ações: lápis navega para a edição; lixeira abre diálogo de confirmação e, ao confirmar, chama o backend e atualiza a tabela.

- [ ] Adicionar o item "{{rotulo-menu}}" no menu lateral apontando para a listagem de `{{entidades}}`.

- [ ] Acrescentar no i18n as chaves novas que aparecerem (ex.: `{{entidade}}.not_found` e mensagens específicas de validação). Reaproveitar as chaves já cadastradas em specs anteriores.

- [ ] Rodar `npx tsc --noEmit` em `apps/frontend` e sinalizar ao usuário que a UI está pronta para conferência manual.

## Resultado Esperado

- Model `{{Entidade}}` com struct validada, interface de repositório e serviço implementados e testados.
- Tabela `{{entidades}}` criada no PostgreSQL com migration aplicada.
- CRUD de `{{entidade}}` exposto no backend via handlers Go com chi, com cenários cobertos no `{{entidade}}.http`.
- Listagem paginada, formulário compartilhado entre criação e edição e exclusão com confirmação funcionando no front-end, acessíveis pelo item "{{rotulo-menu}}" do menu lateral.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
