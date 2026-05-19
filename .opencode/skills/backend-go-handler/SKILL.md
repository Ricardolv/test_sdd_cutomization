---
name: backend-go-handler
description: Cria handlers HTTP no backend Go para expor casos de uso via chi, com testes usando httptest.
---

# Backend Go Handler

## Objetivo

Criar ou atualizar handlers HTTP em `apps/backend` para expor casos de uso de `modules/*` usando chi, mantendo o handler fino e focado em HTTP.

## Entradas obrigatorias

1. modulo
2. agregado
3. caso de uso
4. metodo HTTP
5. path do endpoint

## Workflow

1. Ler o caso de uso real e sua entrada/saida.
2. Criar handler em `internal/http/handlers/<modulo>/`.
3. Montar router com chi:
   - `r.Route("/<modulo>", func(r chi.Router) { ... })`
4. Mapear payloads:
   - `json.Decoder` para body.
   - `chi.URLParam` para params.
5. Usar middlewares de erro e autenticacao existentes.
6. Criar teste com `net/http/httptest` e `chi.NewRouter()`.

## Guardrails

- Nao colocar regra de negocio no handler.
- Reutilizar contratos do caso de uso quando possivel.
- Nao duplicar tratamento de erro se existir helper compartilhado.

## Saida esperada

- Handler registrado no router principal.
- Teste cobrindo sucesso e erro basico.
