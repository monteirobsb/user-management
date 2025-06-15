# Projeto: Aplicação de Gerenciamento de Usuários com Go, Vue.js e PostgreSQL

Este é um projeto full-stack para criar um sistema de gerenciamento de usuários.

## Tecnologias

* **Backend:** Go, Gin
* **Frontend:** Vue.js
* **Banco de Dados:** PostgreSQL
* **DevOps:** Docker

---

## Esquema do Banco de Dados

### Tabela: `users`

A tabela principal para armazenar os dados dos usuários.

| Coluna        | Tipo         | Restrições                          | Descrição                                                  |
| :------------ | :----------- | :---------------------------------- | :--------------------------------------------------------- |
| `id`          | `UUID`       | `PRIMARY KEY`                       | Identificador único do usuário (gerado via `gen_random_uuid()`) |
| `name`        | `VARCHAR(255)`| `NOT NULL`                          | Nome completo do usuário                                   |
| `email`       | `VARCHAR(255)`| `UNIQUE NOT NULL`                   | Endereço de e-mail (usado para login)                      |
| `password_hash`| `TEXT`       | `NOT NULL`                          | Hash da senha do usuário                                   |
| `created_at`  | `TIMESTAMPTZ`| `NOT NULL DEFAULT NOW()`            | Data e hora de criação do registro                         |
| `updated_at`  | `TIMESTAMPTZ`| `NOT NULL DEFAULT NOW()`            | Data e hora da última atualização                          |

---

## Como Executar o Projeto

1.  Clone o repositório.
2.  Renomeie o arquivo `.env.example` para `.env` e, se desejar, altere as credenciais.
3.  Execute o comando na raiz do projeto:
    ```bash
    docker-compose up --build
    ```
4.  Acesse o frontend em `http://localhost:8080`.
5.  A API estará disponível em `http://localhost:8000`.