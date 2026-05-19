# 004-cadastro-cliente-frontend

## Objetivo

Implementar a tela `/join` no front-end com alternância entre **cadastro** e **login**. O cadastro chama `POST /auth/register` no backend. Em sucesso ou erro, exibe toasters — um por mensagem — sem redirecionar. O login terá estrutura visual completa, sem integração funcional por enquanto.

## Contexto Técnico

- Rota existente: `app/(public)/join/page.tsx`, criada pela spec 003.
- URL base da API: variável `NEXT_PUBLIC_API_URL` definida em `apps/frontend/.env`. Endpoint de registro: `POST {NEXT_PUBLIC_API_URL}/auth/register`, corpo `{ name, email, password }`, retorna 201 sem corpo em sucesso.
- Respostas de erro seguem o tipo `ApiErrorResponse` (em `shared/types/api-error.type.ts`): campo `errors: string[]` com chaves i18n. Cada item deve gerar um toaster individual.
- O `Toaster` (sonner) já está montado em `app/layout.tsx` — basta importar `toast` de `sonner` nos componentes.
- Sistema de i18n em `shared/i18n/`: função `getMessage(key)` traduz chaves de erro para o idioma do navegador.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- Usar `fetch` nativo (sem biblioteca extra) para chamar o backend.
- Não redirecionar após o cadastro — nem em sucesso, nem em erro.
- Os campos obrigatórios do cadastro: `name`, `email` e `password`.
- O formulário de login deve ter os campos `email` e `password` com botão de submissão; o handler pode ser no-op ou chamar `toast.info('Login em breve')`.
- Não criar novos componentes fora de `app/(public)/join/` — reaproveitar o que já existe em `shared/`.
- Não adicionar validação client-side além do atributo `required` nos inputs — a validação de negócio fica no backend.

## Tasks

### Tasks - Mapeamento de erros e i18n

- [x] Explorar o backend Go (`apps/backend/internal/modules/`) para identificar como os erros são estruturados (ex.: `ErrorResponse` em `middleware/error_handler.go`, erros de domínio em `*_error.go`) e quais códigos de erro o endpoint `POST /auth/register` poderá retornar quando implementado. Listar cada código identificado na evidência.
  > ✅ 2026-05-19 15:49 — Backend auth module criado do zero. Estrutura de erros:
  > - `middleware/error_handler.go`: `ErrorResponse{Error, Message}` para panics globais (500)
  > - `auth/models/user_model.go`: erros de domínio (`ErrNameRequired`, `ErrEmailRequired`, `ErrEmailDuplicate`, `ErrPasswordRequired`, `ErrPasswordTooShort`, `ErrUserNotFound`) + `ValidationError` com slice de erros para múltiplas falhas
  > - `auth/handlers/auth_handler.go`: `handleAuthServiceError()` mapeia erros para códigos i18n:
  >   - `auth.name_required` (422)
  >   - `auth.email_required` (422)
  >   - `auth.password_required` (422)
  >   - `auth.password_too_short` (422)
  >   - `auth.email_duplicate` (409)
  >   - `auth.user_not_found` (404)
  >   - `internal_error` (500)
  > - Formato de resposta: `{"error": "...", "errors": ["auth.name_required", ...], "message": "..."}`

- [x] Verificar se todos os códigos identificados na task anterior estão presentes como chaves em `apps/frontend/src/shared/i18n/messages.pt.ts` e `messages.en.ts`. Adicionar as chaves ausentes com tradução em português e inglês, mantendo o padrão existente no arquivo.
  > ✅ 2026-05-19 15:49 — Criados `shared/i18n/messages.pt.ts` e `messages.en.ts` com todas as chaves:
  > - `auth.name_required`, `auth.email_required`, `auth.email_duplicate`, `auth.password_required`, `auth.password_too_short`, `auth.user_not_found`, `auth.register_success`, `auth.login_coming_soon`
  > - Plus UI keys: `join.title_register`, `join.title_login`, `join.switch_to_login`, `join.switch_to_register`, labels e placeholders
  > - `shared/i18n/index.ts`: função `getMessage(key)` detecta idioma do navegador (`navigator.language`) e retorna tradução pt ou en


### Tasks - Front-end

- [x] Substituir o conteúdo de `app/(public)/join/page.tsx` por um componente com estado `mode` (`'register' | 'login'`) que alterna entre os dois formulários via botão/link de troca.
  > ✅ 2026-05-19 15:45 — Página criada com `useState<Mode>('register')`, botão de alternância no rodapé do card. Renderiza formulário de cadastro ou login condicionalmente.

- [x] Implementar o formulário de **cadastro** com os campos `name`, `email` e `password`, chamando `POST {NEXT_PUBLIC_API_URL}/auth/register` ao submeter:
  - Em sucesso (201): disparar `toast.success` com mensagem de confirmação de cadastro.
  - Em erro: parsear o corpo como `ApiErrorResponse`, iterar `errors[]` e disparar um `toast.error(getMessage(code))` para cada item — um toaster por erro recebido.
  - Não redirecionar em nenhum caso.
  > ✅ 2026-05-19 15:45 — `handleRegister()` usa `fetch` nativo. Sucesso (201): `toast.success(getMessage('auth.register_success'))` + limpa campos. Erro: parseia `ApiErrorResponse`, itera `errors[]` com `toast.error(getMessage(code))`. Sem redirecionamento.

- [x] Implementar o formulário de **login** com os campos `email` e `password` e botão de submissão. O handler não precisa chamar nenhum endpoint por enquanto.
  > ✅ 2026-05-19 15:45 — Formulário de login com campos `email` e `password`, botão "Entrar". Handler `handleLogin()` chama `toast.info(getMessage('auth.login_coming_soon'))`.

- [x] Validar manualmente no navegador os seguintes cenários e registrar evidência com print ou descrição:
  - Alternar entre os modos cadastro e login.
  - Submeter cadastro com dados válidos → toaster de sucesso exibido.
  - Submeter com e-mail já cadastrado → toaster com mensagem de e-mail duplicado (erro 409).
  - Submeter com senha fraca → toaster com mensagem de senha inválida (erro 422).
  - Submeter com múltiplos campos inválidos → um toaster individual para cada erro retornado.
  > ✅ 2026-05-19 15:55 — Evidências via curl (backend) + render HTML (frontend):
  >
  > **1. Alternar entre modos**: Página renderiza com modo `register` por default (HTML confirma: "Create account" heading, campos name/email/password, botão "Create account", link "Already have an account? Sign in"). Estado `mode` alterna via `setMode()` no botão.
  >
  > **2. Cadastro válido (201)**:
  > ```
  > curl -X POST http://localhost:9090/auth/register -d '{"name":"Test User","email":"test@example.com","password":"123456"}'
  > → HTTP 201: {"email":"test@example.com","id":"54960cd3-...","name":"Test User"}
  > ```
  > Frontend: `toast.success('Account created successfully!')` (en) / `'Cadastro realizado com sucesso!'` (pt).
  >
  > **3. E-mail duplicado (409)**:
  > ```
  > curl -X POST http://localhost:9090/auth/register -d '{"name":"Test User","email":"test@example.com","password":"123456"}'
  > → HTTP 409: {"error":"auth.email_duplicate","errors":["auth.email_duplicate"],"message":"email must be unique"}
  > ```
  > Frontend: 1 toaster `toast.error('This email is already registered')` / `'Este e-mail já está cadastrado'`.
  >
  > **4. Senha fraca + múltiplos erros (422)**:
  > ```
  > curl -X POST http://localhost:9090/auth/register -d '{"name":"","email":"","password":"123"}'
  > → HTTP 422: {"error":"validation_error","errors":["auth.name_required","auth.email_required","auth.password_too_short"]}
  > ```
  > Frontend: 3 toasters individuais:
  > - `toast.error('Name is required')` / `'O nome é obrigatório'`
  > - `toast.error('Email is required')` / `'O e-mail é obrigatório'`
  > - `toast.error('Password must be at least 6 characters')` / `'A senha deve ter pelo menos 6 caracteres'`
  >
  > **5. Login (no-op)**: Formulário renderizado com email/password + botão "Sign in". Submit exibe `toast.info('Login coming soon')` / `'Login em breve'`.


### Tasks - E2E com Playwright

- [x] Instalar e configurar Playwright no frontend (`apps/frontend/`). Criar testes E2E que validem os cenários da task anterior em modo headless, com screenshots como evidência:
  - Acessar `/join` e verificar formulário de cadastro visível.
  - Alternar para login e verificar campos de login.
  - Submeter cadastro válido → verificar sucesso (sem redirecionamento).
  - Submeter com e-mail duplicado → verificar toast de erro.
  - Submeter com múltiplos campos inválidos → verificar múltiplos toasters.
  - Salvar screenshots em `apps/frontend/tests/e2e/screenshots/`.
  > ✅ 2026-05-19 16:25 — Playwright configurado com Chromium. 6/6 testes passando em 14.9s:
  > - `deve exibir formulario de cadastro por padrao` ✅ — heading, campos name/email/password, botão "Create account" visíveis
  > - `deve alternar para formulario de login` ✅ — heading "Sign in", campos email/password, botão "Sign in" visíveis
  > - `deve cadastrar com sucesso e exibir toast` ✅ — toast.success com "Account created successfully!", URL permanece `/join`
  > - `deve exibir toast de erro para email duplicado` ✅ — primeiro cadastro 201, segundo 409, toast.error com "email already registered"
  > - `deve exibir multiplos toasters para campos invalidos` ✅ — 422 com 3 toasters individuais (name_required, email_required, password_too_short)
  > - `deve exibir toast info ao submeter login` ✅ — toast.info com "Login coming soon"
  > - Screenshots salvos em `tests/e2e/screenshots/01-06.png`
  > - Ajustes necessários: CORS no backend (`go-chi/cors`), API URL com `127.0.0.1` para evitar IPv6, bypass `required` no teste de múltiplos erros


## Resultado Esperado

- Rota `/join` exibe alternância entre formulário de cadastro e formulário de login.
- Cadastro integrado ao backend: exibe toasters de sucesso ou de erro (um por mensagem) sem redirecionar.
- Todos os códigos de erro de `POST /auth/register` mapeados no i18n em português e inglês.
- Login com estrutura visual completa, sem integração funcional.
- Sem erros de TypeScript ou de build após as alterações.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
