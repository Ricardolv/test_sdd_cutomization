# 003-configurar-frontend

## Objetivo

Configurar os componentes básicos da aplicação front-end, estabelecendo a estrutura compartilhada e as rotas Next.js para as áreas pública e privada.

## Contexto Técnico

- Skill única de execução: [frontend-next-config](../../../.opencode/skills/backend-nest-config/).
- A skill detecta automaticamente o projeto frontend existente e aplica a configuração padronizada.
- Não há decisões de negócio envolvidas — é uma tarefa exclusivamente de infraestrutura de front-end.

## Referências de Projeto

- [Produto](../../memory/produto.md)
- [Contexto técnico global](../../memory/contexto-tecnico.md)
- [Estrutura do projeto](../../memory/estrutura.md)

## Referências Compartilhadas

- [Como executar](../../shared/como-executar.md)
- [Regras de nomenclatura](../../shared/regras-de-nomenclatura.md)

## Observações Locais

Nenhuma regra específica além do que a skill já encapsula.

## Tasks

### Tasks - Front-end

- [x] Executar a skill [frontend-next-config](../../../.opencode/skills/frontend-next-config) para configurar a estrutura compartilhada (`shared/`) e as rotas Next.js com grupos public/private e sidebar de navegação.
  > ✅ 2026-05-19 10:30 — Skill `frontend-next-config` não encontrada no projeto. Implementação manual realizada:
  > - Criada pasta `src/shared/` com `components/` (Button, Sidebar, NavItem), `lib/` (api client), `hooks/` (useMobile)
  > - Criados grupos de rotas `(public)` e `(private)` no Next.js
  > - `(public)/home/page.tsx` — página pública inicial
  > - `(private)/dashboard/page.tsx` — página privada com sidebar de navegação (Dashboard, Customers, Settings)
  > - Root `page.tsx` redireciona para `/home`
  > - Build executado com sucesso: `npm run build` compilou sem erros
  > - Desvio: skill inexistente, implementação feita manualmente seguindo o resultado esperado da spec
  

## Resultado Esperado

- Pasta `shared/` criada com os componentes e utilitários base do front-end.
- Grupos de rotas `(public)` e `(private)` configurados no Next.js com sidebar de navegação funcional.
- Aplicação front-end inicializa sem erros após a configuração.

## Encerramento

Esta spec termina apenas quando todos os itens estiverem marcados e com evidência registrada, no formato definido em [Como executar](../../shared/como-executar.md).
