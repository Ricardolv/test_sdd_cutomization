# Regras de Nomenclatura

Convenções globais de nomes de arquivos e diretórios. Podem ser referenciadas por qualquer spec.

## Regra geral

- nomes de arquivos e diretórios em `kebab-case`, sempre minúsculas
- nomes devem indicar **responsabilidade**, não implementação
- quando fizer sentido, o sufixo deve explicitar o papel do arquivo
- não usar `PascalCase`, `camelCase` ou mistura de maiúsculas em diretórios

Exemplos de diretórios válidos: `shared`, `pages`, `examples`, `customer-settings`.

## Backend (Go)

| Sufixo              | Uso                                                    | Exemplo                        |
| ------------------- | ------------------------------------------------------ | ------------------------------ |
| `*_model.go`        | structs de domínio/entidades                           | `customer_model.go`            |
| `*_service.go`      | regras de negócio e orquestração                       | `customer_service.go`          |
| `*_repository.go`   | acesso a dados                                         | `customer_repository.go`       |
| `*_handler.go`      | handlers HTTP (chi)                                    | `customer_handler.go`          |
| `*_middleware.go`   | middlewares                                            | `auth_middleware.go`           |
| `*_router.go`       | roteamento                                             | `customer_router.go`           |
| `*_config.go`       | configuração                                           | `db_config.go`                 |
| `*_test.go`         | testes                                                 | `customer_service_test.go`     |
| `*_mock.go`         | mocks para testes                                      | `customer_repository_mock.go`  |
| `*_dto.go`          | data transfer objects (request/response)               | `customer_dto.go`              |
| `*_error.go`        | erros customizados                                     | `customer_error.go`            |

## Frontend (Next.js)

| Sufixo              | Uso                                                    | Exemplo                        |
| ------------------- | ------------------------------------------------------ | ------------------------------ |
| `*.page.tsx`        | páginas (App Router)                                   | `login/page.tsx`               |
| `*.component.tsx`   | componentes                                            | `customer-form.component.tsx`  |
| `*.context.tsx`     | contextos React e hooks associados                     | `toast.context.tsx`            |
| `*.provider.tsx`    | providers de composição, wrappers globais              | `app.providers.tsx`            |
| `*.hook.ts`         | custom hooks                                           | `use-auth.hook.ts`             |
| `*.store.ts`        | stores                                                 | `session.store.ts`             |
| `*.types.ts`        | tipos auxiliares                                       | `menu.types.ts`                |
| `*.api.ts`          | funções de chamada à API                               | `customer.api.ts`              |
| `*.spec.ts`         | testes automatizados                                   | `customer-form.spec.ts`        |

### Exemplos aplicados

**Backend (Go):**
- `customer_model.go`
- `customer_service.go`
- `customer_repository.go`
- `customer_handler.go`
- `auth_middleware.go`
- `customer_router.go`
- `customer_dto.go`
- `customer_error.go`
- `customer_service_test.go`

**Frontend (Next.js):**
- `login/page.tsx`
- `customer-form.component.tsx`
- `toast.context.tsx`
- `use-auth.hook.ts`
- `session.store.ts`
- `app.providers.tsx`
- `customer.api.ts`
- `menu.types.ts`

## Exceções controladas

Nomes exigidos por ferramentas ou convenções externas mantêm o formato original. Exemplos: `README.md`, `SKILL.md`, `go.mod`, `go.sum`, `package.json`, `tsconfig.json`, `spec.md`, `main.go`.

Fora dessas exceções, prefira sempre `kebab-case`.

## Regra de decisão

Se um nome estiver ambíguo, prefira a forma que deixe mais claro:

- o que o arquivo representa
- em que camada ele vive
- qual a responsabilidade principal
