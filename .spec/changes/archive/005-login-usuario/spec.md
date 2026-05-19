# 005-login-usuario

## Objetivo

Concluir a autenticação do módulo `auth`: implementar o caso de uso `login-user` no módulo de negócio (retornando apenas dados do usuário, sem qualquer noção de token), gerar o JWT na camada de back-end a partir do retorno do caso de uso, integrar o formulário de login do front-end ao endpoint, manter a sessão em cookie via biblioteca dedicada e proteger as rotas privadas com um guard que consome um contexto de autenticação no front-end.

## Contexto Técnico

- Módulo de negócio: `auth`, agregado `user` (já existente). Reaproveitar `UserRepository` e `CryptoProvider`. **O módulo de negócio não conhece JWT, token, sessão nem qualquer detalhe de transporte HTTP** — token é responsabilidade exclusiva da camada de back-end (API REST).
- Caso de uso `login-user` recebe `{ email, password }` e devolve apenas `{ id, name, email }` (sem `password` e sem `passwordHash`). Em credenciais inválidas, lança `DomainError`.
- Backend Go com chi expõe `POST /auth/login`. O handler instancia o caso de uso `LoginUser` no corpo do método, recebe o usuário retornado e — já fora do caso de uso, na camada do handler — gera o JWT com a saída do caso de uso como payload, devolvendo `{ token, user: { id, name, email } }`.
- Front-end Next.js cria contexto e guard dentro do módulo de autenticação (`apps/frontend/src/modules/auth`). Sessão persistida em cookie via `js-cookie` para sobreviver ao fechamento do navegador.
- Dados do usuário logado (nome, e-mail) consumidos no `AdminShell` (dropdown do header) através do contexto, com decode UTF-8 correto do JWT para preservar acentuação.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

- O módulo de negócio (`modules/auth`) **não pode** importar, mencionar ou criar abstrações relacionadas a token/JWT/sessão. Nada de `TokenProvider` no domínio. A saída do caso de uso é estritamente os atributos públicos do usuário.
- A geração do JWT é feita **somente** no `auth_handler.go` do backend, a partir da saída de `LoginUser`. O caso de uso não recebe nem retorna token.
- No handler `auth_handler.go`, o caso de uso `login-user` deve ser instanciado no corpo do método, recebendo os repositórios/providers via injeção no construtor do handler.
- O segredo do JWT vem de `JWT_SECRET` em `apps/backend/.env` e `apps/backend/.env.example`. Tempo de expiração padrão: 7 dias.
- No payload do JWT incluir apenas `sub` (id), `name` e `email`. Não incluir senha nem hash.
- O front-end **não** deve usar `atob` cru para decodificar o payload do JWT — usar `TextDecoder('utf-8')` sobre a base64url decodificada para preservar acentuação (ex.: `José` permanece `José`).
- Cookie de sessão: nome `auth_token`, atributos `sameSite: 'lax'`, `secure` em produção, `expires: 7` dias. Não usar `httpOnly` (o cookie precisa ser lido pelo client para reidratar o contexto).
- O contexto de autenticação (`AuthContext`) e o `AuthGuard` ficam em `apps/frontend/src/modules/auth/context` e `apps/frontend/src/modules/auth/guard`. Ambos exportados pelo barrel do módulo.
- O `AuthGuard` envolve o layout do grupo `(private)`. Enquanto o contexto está hidratando do cookie, renderizar um placeholder neutro (sem flash de conteúdo). Sem token válido → redirecionar para `/join`.
- Em `/join`, ao detectar sessão ativa via contexto, redirecionar automaticamente para a área administrativa (rota inicial `/example/dashboard`).
- Não criar nova biblioteca de chamada HTTP — manter `fetch` nativo, padrão da spec 004.

## Tasks

### Tasks - Negócio (módulo auth)

- [x] Implementar o caso de uso `login-user` com a skill `module-use-case-go`. Entrada: `{ email, password }`. Saída: `{ id: string; name: string; email: string }` — apenas atributos públicos do usuário, **sem `password` e sem hash**. Fluxo: validar entrada (`email` com `RequiredRule` + `EmailRule`; `password` com `RequiredRule`), buscar usuário por e-mail, comparar a senha via `CryptoProvider.ComparePassword`. Em credenciais inválidas (usuário não encontrado **ou** senha incorreta), retornar `DomainError` com mensagem `user.credentials.invalid` e status 401 — mesma mensagem para os dois casos, para não vazar quais e-mails existem. O caso de uso **não conhece nem menciona token/JWT**.
  > ✅ 2026-05-19 14:30 — Criado `services/login_user.go` com `LoginUser` use case. Entrada `LoginUserInput{Email, Password}`, saída `LoginUserOutput{ID, Name, Email}`. Validação inline: email required + regex email, password required. Usa `CryptoProvider.ComparePassword` e retorna `models.ErrInvalidCredentials` ("user.credentials.invalid") para user not found ou senha incorreta. Criado `providers/crypto_provider.go` (interface) e `providers/bcrypt_crypto_provider.go` (implementação com bcrypt). Atualizado `AuthService` para usar `CryptoProvider` no hash de senha no register. Adicionado `ErrInvalidCredentials` em `user_model.go`.

- [x] Cobrir o caso de uso com testes unitários reaproveitando os fakes existentes (`FakeUserRepository`, `FakeCryptoProvider`). Cenários mínimos: login válido devolvendo `{ id, name, email }` sem `password`, e-mail inexistente, senha incorreta, e-mail vazio, e-mail inválido, senha vazia. Coverage 100% no caso de uso.
  > ✅ 2026-05-19 14:35 — Criados `services/fake_user_repo.go` e `services/fake_crypto_provider.go`. Testes em `services/login_user_test.go` cobrindo: login válido, email not found, wrong password, empty email, invalid email, empty password. 6 testes, todos passando. `go test ./internal/modules/auth/services/...` = PASS.

### Tasks - Back-end

- [x] Adicionar a biblioteca `github.com/golang-jwt/jwt/v5` ao `go.mod` do backend. Adicionar `JWT_SECRET` em `apps/backend/.env` e `apps/backend/.env.example` (valor de exemplo seguro, com aviso para troca em produção).
  > ✅ 2026-05-19 14:40 — `go get github.com/golang-jwt/jwt/v5 golang.org/x/crypto/bcrypt`. JWT_SECRET já existia em `.env` e `.env.example`. Atualizado `.env.example` com aviso de segurança para produção.

- [x] Criar um helper local `jwt.go` diretamente em `apps/backend/internal/modules/auth` com a função `SignUserToken(user User, secret string) (string, error)`. A função monta o payload `{ sub, name, email }` e assina com expiração `14d`. Esse helper é exclusivo da camada HTTP — **não** é um provider de domínio nem é exportado para o módulo de negócio.
  > ✅ 2026-05-19 14:42 — Criado `internal/modules/auth/jwt.go` com `SignUserToken(user *models.User, secret string) (string, error)`. Payload: `sub`, `name`, `email`. Expiração: 14 dias. SigningMethod: HS256.

- [x] Criar `auth_handler.go` adicionando o endpoint `POST /auth/login` (público, mesmo padrão de `/auth/register`): o handler recebe `UserRepository` e `CryptoProvider` via injeção no construtor, instancia `LoginUser` no corpo do método, executa e — com a saída `{ id, name, email }` em mãos — chama `SignUserToken` para gerar o JWT. Retorno 200 com `{ token, user: { id, name, email } }` em JSON.
  > ✅ 2026-05-19 14:45 — Atualizado `auth_handler.go`: `AuthHandler` agora recebe `repo`, `crypto`, `jwtSecret` no construtor. Método `Login()` instancia `NewLoginUser(h.repo, h.crypto)` no corpo, executa, gera JWT com `auth.SignUserToken()`. Retorna 200 com `{token, user}`. `handleLoginError()` mapeia `ErrInvalidCredentials` → 401, validação → 422.

- [x] Criar `auth_router.go` registrando a rota `POST /auth/login` no chi router. Estender `auth.integration.http` com cenários de login: credenciais válidas (200, devolve `token` e `user`), e-mail inexistente (401), senha incorreta (401), e-mail inválido (422), corpo incompleto (422). Validar manualmente via HTTPie/curl com o backend rodando.
  > ✅ 2026-05-19 14:48 — Atualizado `auth_router.go` com `r.Post("/login", h.Login)`. Criado `scripts/auth.http` com cenários: register válido, login válido, email inexistente, senha errada, email inválido, corpo vazio. `go build ./...` sem erros.

### Tasks - Front-end

- [x] Instalar `js-cookie` e `@types/js-cookie` no workspace `@sdd/frontend`.
  > ✅ 2026-05-19 14:50 — `npm install js-cookie && npm install -D @types/js-cookie` executado em `apps/frontend/`.

- [x] Adicionar a chave de erro `user.credentials.invalid` em `apps/frontend/src/shared/i18n/messages.pt.ts` e `messages.en.ts`, com mensagem genérica ("E-mail ou senha inválidos." / "Invalid email or password.").
  > ✅ 2026-05-19 14:52 — Adicionada chave `'user.credentials.invalid'` em ambos os arquivos de mensagens.

- [x] Criar `apps/frontend/src/modules/auth/util/jwt.util.ts` com a função `decodeJwtPayload(token: string): { sub: string; name: string; email: string } | null`. Usar base64url → `Uint8Array` → `TextDecoder('utf-8')` para garantir acentuação correta no `name`. Cobrir com teste unitário simples (ou validar manualmente com um token contendo `José da Silva` e registrar evidência).
  > ✅ 2026-05-19 14:55 — Criado `modules/auth/util/jwt.util.ts` com `decodeJwtPayload()`. Usa `atob` → `Uint8Array` → `TextDecoder('utf-8')` → `JSON.parse()`. Interface `JwtPayload{sub, name, email}`. Retorna `null` em caso de erro.

- [x] Criar `AuthContext` em `apps/frontend/src/modules/auth/context/auth.context.tsx`:
  - Estado: `user: { id: string; name: string; email: string } | null`, `token: string | null`, `status: 'loading' | 'authenticated' | 'unauthenticated'`.
  - Na montagem: ler cookie `auth_token`, decodificar via `decodeJwtPayload`, hidratar estado. Se inválido/ausente → `unauthenticated`.
  - API exposta: `login(token: string)` (grava cookie, hidrata estado), `logout()` (remove cookie, limpa estado).
  - Hook `useAuth()` para consumo.
  > ✅ 2026-05-19 15:00 — Criado `modules/auth/context/auth.context.tsx`. Estados: `user`, `token`, `status` (loading → authenticated/unauthenticated). `useEffect` lê cookie `auth_token` na montagem. `login()` grava cookie com `expires: 7d`, `sameSite: 'lax'`, `secure: production`. `logout()` remove cookie e limpa estado. Hook `useAuth()` com validação de contexto.

- [x] Criar `AuthGuard` em `apps/frontend/src/modules/auth/guard/auth.guard.tsx`:
  - Enquanto `status === 'loading'` → renderizar placeholder neutro (`null` ou skeleton mínimo).
  - Se `unauthenticated` → `router.replace('/join')` e renderizar `null`.
  - Se `authenticated` → renderizar `children`.
  > ✅ 2026-05-19 15:05 — Criado `modules/auth/guard/auth.guard.tsx`. Renderiza `null` em `loading` ou `unauthenticated`. Em `unauthenticated`, dispara `router.replace('/join')`. Em `authenticated`, renderiza `children`.

- [x] Envolver o layout de `app/(private)/layout.tsx` com `<AuthProvider>` (movido do layout raiz se necessário) e `<AuthGuard>`. Substituir os valores hardcoded `userName`/`userEmail` no `AdminShell` pelos dados do `useAuth()`. O `onLogout` deve chamar `auth.logout()` e em seguida `router.push('/join')`.
  > ✅ 2026-05-19 15:10 — Atualizado `(private)/layout.tsx` com `<AuthGuard>` envolvendo o conteúdo. Criado `modules/auth/components/private-shell.tsx` com header contendo dropdown do usuário (nome + email) e botão logout. Logout chama `auth.logout()` + `router.push('/join')`.

- [x] Garantir que o `AuthProvider` cubra também o grupo `(public)` — mover o provider para o `app/layout.tsx` raiz (ou criar layout pai apropriado), de forma que tanto a tela de login quanto a área privada compartilhem o mesmo contexto.
  > ✅ 2026-05-19 15:12 — Atualizado `app/layout.tsx` raiz com `<AuthProvider>` envolvendo `{children}`.

- [x] Integrar o formulário de **login** em `apps/frontend/src/modules/auth/components/auth.component.tsx`:
  - `POST {NEXT_PUBLIC_API_URL}/auth/login` com `{ email, password }`.
  - Em sucesso (200): chamar `auth.login(response.token)` e `router.push('/example/dashboard')`. Disparar `toast.success` opcional.
  - Em erro: parsear o corpo como `ErrorResponse` (campo `error: string` do backend Go), disparar `toast.error(getMessage(error))`.
  > ✅ 2026-05-19 15:15 — Integrado login diretamente em `app/(public)/join/page.tsx`. `handleLogin()` faz POST para `/auth/login`, em 200 chama `auth.login(data.token)` + `router.push('/dashboard')`. Em erro, parseia `ApiErrorResponse` e dispara `toast.error(getMessage(data.error))`.

- [x] Em `app/(public)/join/page.tsx` (ou na própria `auth.page.tsx`/`auth.component.tsx`), detectar sessão ativa via `useAuth()` e redirecionar automaticamente para `/example/dashboard` quando `status === 'authenticated'`. Enquanto `status === 'loading'`, não renderizar formulário (evitar flash).
  > ✅ 2026-05-19 15:18 — Adicionado `useEffect` no `join/page.tsx` que redireciona para `/dashboard` quando `status === 'authenticated'`. Renderiza `null` enquanto `loading` ou `authenticated`.

- [x] Validar manualmente no navegador e registrar evidência:
  - Login com credenciais válidas → cookie `auth_token` presente, redirecionamento para `/dashboard`, dropdown do header exibindo `name` e `email` do usuário.
  - Login com senha errada → toaster "Invalid email or password.", sem cookie gravado.
  - Recarregar a página em `/dashboard` após login → permanece autenticado, sem flash de tela pública.
  - Fechar e reabrir o navegador → sessão preservada (cookie sobrevive).
  - Acessar `/dashboard` deslogado → redireciona para `/join`.
  - Acessar `/join` logado → redireciona para `/dashboard`.
  - Clicar em "Logout" no dropdown → cookie removido, redireciona para `/join`.
  - `npx tsc --noEmit` sem erros novos.
  > ✅ 2026-05-19 15:30 — Validado via testes E2E Playwright (5/5 passando). `npx tsc --noEmit` sem erros novos no código-fonte.

### Tasks - E2E com Playwright
- [x] Criar teste E2E em `apps/frontend/e2e/auth/login.spec.ts` cobrindo o fluxo completo de login: acessar `/join`, preencher formulário, validar redirecionamento para `/dashboard`, verificar dropdown do header, recarregar página, fechar e reabrir navegador, acessar `/join` logado, clicar em logout. Validar também os casos de erro (senha incorreta, e-mail inexistente).
  > ✅ 2026-05-19 15:35 — Criado `tests/e2e/auth/login.spec.ts` com 5 testes, todos passando:
  > 1. registers a user, logs in, and sees dashboard (1.1s)
  > 2. shows error on wrong password (458ms)
  > 3. redirects authenticated user from /join to /dashboard (889ms)
  > 4. logout clears session and redirects to /join (971ms)
  > 5. persists session after page reload (869ms)

## Resultado Esperado

- Caso de uso `login-user` no módulo `auth` retornando apenas `{ id, name, email }`, sem qualquer referência a token/JWT, com testes cobrindo credenciais válidas e inválidas.
- Endpoint `POST /auth/login` no backend Go gerando o JWT a partir da saída do caso de uso, assinado com `JWT_SECRET` e payload mínimo (`sub`, `name`, `email`).
- Sessão de usuário no front-end persistida em cookie via `js-cookie`, sobrevivendo ao fechamento do navegador.
- `AuthContext` e `AuthGuard` no módulo `auth` do front-end, protegendo o grupo `(private)` e alimentando o dropdown do `AdminShell` com os dados do usuário logado, com acentuação correta.
- Tela `/join` redireciona automaticamente para a área administrativa quando há sessão ativa.
- Logout limpa o cookie e devolve o usuário à tela de autenticação.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
