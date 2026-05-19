---
name: spec-backend-auth-basic-go
description: Orquestra a criacao da base de autenticacao do backend Go com chi, JWT e persistencia sem ORM.
---

# Spec Backend Auth Basic Go

## Objetivo

Orquestrar a base de autenticacao (cadastro e login) no backend Go, usando as skills especializadas de dominio, persistencia e handlers.

## Entradas obrigatorias

1. modulo de autenticacao
2. agregado principal (ex.: user)
3. entidade principal
4. atributos minimos do usuario
5. confirmacao de cadastro e login

## Sequencia deterministica

1. `config-new-module-go`
2. `config-package-shared-go` (se necessario)
3. `module-aggregate-go`
4. `module-entity-go`
5. `module-repository-go`
6. `module-use-case-go` (cadastro)
7. `module-use-case-go` (login)
8. `shared-validation-rule-go` quando faltar regra reutilizavel
9. `backend-go-provider-implementation` (hash e JWT)
10. `backend-go-sync-module`
11. `backend-go-repository`
12. `backend-go-config`
13. `backend-go-handler`

## Regras basicas

- JWT via `github.com/golang-jwt/jwt/v5` quando nao houver IdP externo.
- Senha sempre protegida (bcrypt/argon2).
- Sem ORM; usar `database/sql`.

## Saida esperada

- Modulo de autenticacao funcional.
- Persistencia integrada sem ORM.
- Endpoints de cadastro/login prontos.
