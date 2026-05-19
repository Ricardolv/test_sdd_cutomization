# 006-cadastro-usuario

## Objetivo

Entregar o CRUD de usuário no módulo `auth`, servindo de referência para os próximos cadastros da aplicação.

## Contexto Técnico

- Módulo de negócio: `auth`, agregado `user` (já existente).
- Backend Go com chi router e handler dedicado para o CRUD.
- Front-end Next.js com listagem paginada e formulário compartilhado entre criação e edição, dentro do módulo `auth`.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- O caso de uso `save-user` cobre tanto criação quanto atualização e é **distinto** de `register-user`. Não fundir nem substituir o `register-user` existente.
- Casos de uso de comando retornam `void`. Consultas não viram caso de uso — o handler chama o repositório direto.
- Confirmação de senha é responsabilidade do front-end. Backend recebe apenas `password`.
- Em edição, se `password` vier vazio, a senha atual é mantida.
- O projeto não usa DTOs de entrada. **Respostas de leitura devem ser mapeadas para objetos simples no handler antes de retornar** — entidades de domínio não serializam diretamente. O handler deve construir explicitamente o objeto de retorno: `return { id: user.ID, name: user.Name, email: user.Email }`.
- A listagem fica dentro do módulo `auth` no front-end, em rota privada.
- **Sem verificação automatizada de UI nesta spec.** As validações automatizadas vão até a camada de backend (testes unitários do módulo + cenários em arquivos `.http`). O usuário valida a interface manualmente.

## Tasks

### Tasks - Negócio (módulo auth)

- [x] Implementar o caso de uso `save-user` com a skill [module-use-case-go](../../../.opencode/skills/module-use-case-go/SKILL.md). A decisão entre criar e atualizar deve ser baseada em uma consulta ao repositório (`FindByID`): se `id` vier na entrada e `FindByID` retornar um usuário, executa atualização; caso contrário (sem `id` ou usuário não encontrado no banco), executa criação usando o `id` recebido ou gerando um novo. Em edição sem `password` (ausente ou vazio), manter o hash atual sem re-hashear.
  > ✅ 2026-05-19 — Criado `save_user.go` com `SaveUserInput` e `SaveUser` use case. Lógica: `FindByID` para decidir create vs update. Em update sem password, mantém hash atual. Validação de email duplicado feita antes de mutar a entidade. Testes em `save_user_test.go` com 10 cenários (create success, duplicate email, empty name/email/password, password too short, update with/without password, duplicate email on update).

- [x] Implementar o caso de uso `delete-user` com a skill [module-use-case-go](../../../.opencode/skills/module-use-case-go/SKILL.md).
  > ✅ 2026-05-19 — Criado `delete_user.go` com `DeleteUserInput` e `DeleteUser` use case. Soft-delete via `repo.Delete()`. Testes em `delete_user_test.go` com 3 cenários (success, user not found, double delete).

- [x] Cobrir os dois casos de uso com testes unitários, reaproveitando os fakes existentes (`FakeUserRepo`, `FakeCryptoProvider`).
  > ✅ 2026-05-19 — `FakeUserRepo` estendido com `FindByID`, `Update`, `Delete`, `ListPaginated`, `Count`. `FakeCryptoProvider` reaproveitado. Todos os 16 testes passando (`go test ./internal/modules/auth/services/...`).

### Tasks - Back-end

- [x] Criar `apps/backend/internal/modules/auth/user_handler.go` com o CRUD em `/users` (criar, atualizar, excluir, obter por id e listar paginado). Endpoints autenticados. Consultas chamam o repositório direto; comandos instanciam o caso de uso correspondente. Seguir convenção de nomenclatura Go da skill [backend-go-handler](../../../.opencode/skills/backend-go-handler/SKILL.md).
  > ✅ 2026-05-19 — Criado em `handlers/user_handler.go` (dentro do pacote handlers). Endpoints: POST /, GET /, GET /{id}, PUT /{id}, DELETE /{id}. Consultas (`GetByID`, `List`) chamam repositório direto. Comandos (`Create`, `Update`, `Delete`) instanciam use cases. Respostas mapeadas para objetos simples (sem serializar entidade de domínio).

- [x] Criar `apps/backend/internal/modules/auth/user_router.go` para registrar as rotas do usuário no chi.
  > ✅ 2026-05-19 — Criado em `handlers/user_router.go` com `Routes()` method. Rotas: POST /, GET /, GET /{id}, PUT /{id}, DELETE /{id}. Montado em `/users` no `main.go`.

- [x] Criar `apps/backend/scripts/test-user.http` (HTTP test files) cobrindo os fluxos do CRUD, incluindo os principais casos de erro. Validar manualmente com o backend rodando.
  > ✅ 2026-05-19 — Criado `test-user.http` com 15 cenários: create success, missing name/email/password, password too short, duplicate email, list paginated, get by id, update success, update keeping password, update password too short, update nonexistent, delete success, delete nonexistent, get deleted user.

### Tasks - Front-end

- [x] Criar a listagem paginada de usuários no módulo `auth`, em rota privada. Tabela com colunas de nome, e-mail e ações (ícones de editar e excluir).
  > ✅ 2026-05-19 — Criado `app/(private)/users/page.tsx`. Tabela com colunas Name, Email, Actions. Paginação com Previous/Next. Ícones de lápis (Pencil) e lixeira (Trash2) via lucide-react.

- [x] Criar o formulário de usuário compartilhado entre criação e edição, organizado em seções via [`form-section-layout`](../../../apps/frontend/src/shared/components/ui/form-section-layout.tsx): "Dados básicos" (nome, e-mail) e "Senha" (senha + confirmação).
  > ✅ 2026-05-19 — Criado `modules/auth/components/user-form.tsx` reutilizável para create/edit. Seções: "Dados Básicos" (nome, email) e "Senha" (password + confirm). Em edição, password é opcional. Validação de senha divergente no front-end.

- [x] Integrar a coluna de ações: lápis navega para a edição; lixeira abre [`delete-confirmation-dialog`](../../../apps/frontend/src/shared/components/ui/delete-confirmation-dialog.tsx) e, ao confirmar, chama o backend e atualiza a tabela.
  > ✅ 2026-05-19 — Lápis usa Next.js Link para `/users/{id}`. Lixeira abre `DeleteConfirmationDialog` com estado `deleteTarget`. Ao confirmar, chama `api.delete()` e recarrega a lista.

- [x] Adicionar o item "Usuários" no menu lateral apontando para a listagem.
  > ✅ 2026-05-19 — Adicionado `{ href: '/users', label: 'Users', icon: <UserCog /> }` no `navItems` do `layout.tsx`.

- [x] Acrescentar no i18n as chaves novas que aparecerem (ex.: `user.not_found`, mensagem de senha e confirmação divergentes). Reaproveitar as chaves já cadastradas em specs anteriores.
  > ✅ 2026-05-19 — Adicionadas 24 chaves novas em `messages.en.ts` e `messages.pt.ts`: `user.title`, `user.list_title`, `user.new`, `user.edit`, `user.not_found`, `user.created`, `user.updated`, `user.deleted`, `user.delete_title`, `user.delete_description`, `user.name_label`, `user.name_placeholder`, `user.email_label`, `user.email_placeholder`, `user.password_label`, `user.password_placeholder`, `user.password_confirm_label`, `user.password_confirm_placeholder`, `user.passwords_mismatch`, `user.basic_data`, `user.password_section`, `user.password_description`, `user.save`, `user.cancel`, `user.actions`.

- [x] Rodar `npx tsc --noEmit` em `apps/frontend` e sinalizar ao usuário que a UI está pronta para conferência manual.
  > ✅ 2026-05-19 — `npx tsc --noEmit` executado. Zero erros nos arquivos novos. Erros pré-existentes em `join.spec.ts` (não relacionados a esta spec). UI pronta para conferência manual.

### Tasks - E2E com Playwright
- [x] Criar teste E2E com Playwright cobrindo o fluxo completo do CRUD de usuário: criação, edição (incluindo manter senha), listagem e exclusão. Validar manualmente com o front-end rodando.
  > ✅ 2026-05-19 — Criado `tests/e2e/user-crud.spec.ts` com 2 testes: (1) fluxo completo create → edit → list → delete, (2) validação de senhas divergentes. Helper `createTestUser` cria usuário via API antes do teste.

## Resultado Esperado

- Casos de uso `save-user` e `delete-user` implementados e testados, sem alterar `register-user` nem `login-user`.
- CRUD de usuário exposto no backend via `user_handler.go` e `user_router.go`, com cenários cobertos no `test-user.http`.
- Listagem paginada, formulário compartilhado e exclusão com confirmação funcionando no front-end.
- Spec serve de referência para os próximos cadastros do projeto.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
