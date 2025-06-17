# Projeto: Aplicação de Gerenciamento de Usuários com Go, Vue.js e PostgreSQL

Este é um projeto full-stack para criar um sistema de gerenciamento de usuários.

## Tecnologias

* **Backend:** Go, Gin
* **Frontend:** Vue.js
* **Banco de Dados:** PostgreSQL
* **DevOps:** Docker

---

## Configuração do Backend

O backend requer a configuração de variáveis de ambiente para seu correto funcionamento. Estas variáveis são tipicamente definidas em um arquivo `.env` na raiz do projeto backend (ou configuradas diretamente no ambiente de execução).

### Variáveis de Ambiente

Copie o arquivo `.env.example` para `.env` e preencha os valores necessários.

| Variável          | Obrigatório | Descrição                                                                                                | Exemplo/Padrão |
| :---------------- | :---------- | :------------------------------------------------------------------------------------------------------- | :------------- |
| `JWT_SECRET_KEY`  | **Sim**     | Chave secreta para assinar os tokens JWT. Crítica para a segurança da autenticação. **A aplicação não iniciará sem esta chave.** | `sua_chave_secreta_super_segura` |
| `API_PORT`        | Não         | Porta em que a API do backend será executada.                                                            | `8000`         |
| `DATABASE_HOST`   | **Sim**     | Endereço do servidor do banco de dados PostgreSQL.                                                       | `db` (nome do serviço Docker) |
| `POSTGRES_USER`   | **Sim**     | Nome de usuário para conexão com o PostgreSQL.                                                           | `user`         |
| `POSTGRES_PASSWORD`| **Sim**     | Senha para o usuário do PostgreSQL.                                                                      | `password`     |
| `POSTGRES_DB`     | **Sim**     | Nome do banco de dados no PostgreSQL.                                                                      | `userdb`       |
| `DATABASE_PORT`   | **Sim**     | Porta do servidor PostgreSQL.                                                                            | `5432`         |
| `DATABASE_SSLMODE`| Não         | Modo de SSL para a conexão com o PostgreSQL (`disable`, `require`, `verify-full`, etc.).                 | `disable`      |

**Nota:** A aplicação backend irá falhar ao iniciar se as variáveis obrigatórias (`JWT_SECRET_KEY` e as de conexão com o banco de dados) não estiverem definidas.

---

## API Endpoints do Backend

A API do backend é servida sob o prefixo `/api`.

### Autenticação

*   **`POST /api/login`**
    *   **Corpo da Requisição (JSON):**
        ```json
        {
          "email": "user@example.com",
          "password": "yourpassword"
        }
        ```
    *   **Resposta de Sucesso (200 OK):**
        ```json
        {
          "token": "jwt_token_aqui"
        }
        ```
    *   **Respostas de Erro:**
        *   `400 Bad Request`: Payload inválido ou dados ausentes.
        *   `401 Unauthorized`: Credenciais inválidas ou usuário não encontrado.

### Gerenciamento de Usuários

As rotas de gerenciamento de usuários (exceto a criação) são protegidas e requerem um token JWT válido no cabeçalho `Authorization: Bearer <token>`.

*   **`POST /api/users`** (Criação de Usuário - Rota Pública)
    *   **Corpo da Requisição (`models.UserCreateRequest`):**
        ```json
        {
          "name": "John Doe",
          "email": "john.doe@example.com",
          "password": "yoursecurepassword"
        }
        ```
    *   **Regras de Validação:**
        *   `name`: Obrigatório, não pode ser vazio.
        *   `email`: Obrigatório, deve ser um formato de e-mail válido.
        *   `password`: Obrigatório, mínimo de 8 caracteres.
    *   **Resposta de Sucesso (201 Created):** Retorna o objeto do usuário criado (sem o hash da senha).
        ```json
        {
          "id": "uuid-string-aqui",
          "name": "John Doe",
          "email": "john.doe@example.com",
          "created_at": "timestamp",
          "updated_at": "timestamp"
        }
        ```
    *   **Respostas de Erro:**
        *   `400 Bad Request`: Falha na validação dos dados de entrada. O corpo da resposta geralmente contém detalhes sobre os campos inválidos.
        *   `500 Internal Server Error`: Erro ao processar a criação do usuário (e.g., e-mail já existente, falha no banco de dados).

*   **`PUT /api/users/:id`** (Atualização de Usuário - Rota Protegida)
    *   **Parâmetro de URL:** `id` - UUID do usuário a ser atualizado.
    *   **Corpo da Requisição (`models.UserUpdateRequest` - campos opcionais):**
        *   Para atualizar o nome:
            ```json
            {
              "name": "Jane Doe"
            }
            ```
        *   Para atualizar o e-mail:
            ```json
            {
              "email": "jane.doe@newexample.com"
            }
            ```
    *   **Regras de Validação:**
        *   Se `name` for fornecido, não pode ser uma string vazia.
        *   Se `email` for fornecido, deve ser um formato de e-mail válido.
        *   A senha **não pode** ser atualizada através deste endpoint.
    *   **Resposta de Sucesso (200 OK):** Retorna o objeto do usuário atualizado.
    *   **Respostas de Erro:**
        *   `400 Bad Request`: Falha na validação dos dados de entrada ou ID de usuário inválido.
        *   `404 Not Found`: Usuário com o ID fornecido não encontrado.
        *   `500 Internal Server Error`: Erro ao processar a atualização.

*   **`GET /api/users`** (Listar Usuários - Rota Protegida)
    *   Retorna uma lista de todos os usuários.

*   **`GET /api/users/:id`** (Buscar Usuário por ID - Rota Protegida)
    *   Retorna os detalhes do usuário especificado.

*   **`DELETE /api/users/:id`** (Deletar Usuário - Rota Protegida)
    *   Remove o usuário especificado.

---

## Esquema do Banco de Dados

### Tabela: `users`

A tabela principal para armazenar os dados dos usuários.

| Coluna        | Tipo         | Restrições                          | Descrição                                                  |
| :------------ | :----------- | :---------------------------------- | :--------------------------------------------------------- |
| `id`          | `UUID`       | `PRIMARY KEY`                       | Identificador único do usuário (gerado automaticamente pelo backend via GORM hook) |
| `name`        | `VARCHAR(255)`| `NOT NULL`                          | Nome completo do usuário                                   |
| `email`       | `VARCHAR(255)`| `UNIQUE NOT NULL`                   | Endereço de e-mail (usado para login)                      |
| `password_hash`| `TEXT`       | `NOT NULL`                          | Hash da senha do usuário                                   |
| `created_at`  | `TIMESTAMPTZ`| `NOT NULL`                          | Data e hora de criação do registro (gerenciado pelo GORM)  |
| `updated_at`  | `TIMESTAMPTZ`| `NOT NULL`                          | Data e hora da última atualização (gerenciado pelo GORM)   |

**Nota:** O `id` do usuário é um UUID gerado pelo backend na criação do usuário (via hook do GORM). Os campos `created_at` e `updated_at` são gerenciados automaticamente pelo GORM.

---

## Como Executar o Projeto

1.  **Clonar o Repositório:**
    ```bash
    git clone <URL_DO_REPOSITORIO>
    cd <NOME_DO_DIRETORIO_DO_PROJETO>
    ```

2.  **Configurar Variáveis de Ambiente do Backend:**
    *   Navegue até o diretório do backend (ex: `cd backend`).
    *   Copie o arquivo `.env.example` para `.env`.
        ```bash
        cp .env.example .env
        ```
    *   Edite o arquivo `.env` e preencha **todas** as variáveis de ambiente obrigatórias listadas na seção "Configuração do Backend". Isso é crucial para o funcionamento da API, especialmente `JWT_SECRET_KEY` e as credenciais do banco de dados.

3.  **Executar com Docker Compose:**
    *   Volte para o diretório raiz do projeto (onde o arquivo `docker-compose.yml` está localizado).
    *   Execute o comando:
        ```bash
        docker-compose up --build
        ```
    *   Este comando irá construir as imagens Docker para o backend, frontend e o banco de dados PostgreSQL, e então iniciar os contêineres. A flag `--build` garante que as imagens sejam reconstruídas se houverem alterações no código ou Dockerfiles.

4.  **Acessar a Aplicação:**
    *   **Frontend:** Abra seu navegador e acesse `http://localhost:8080` (a porta pode variar dependendo da configuração no `docker-compose.yml` para o serviço `frontend`).
    *   **API Backend:** A API estará disponível em `http://localhost:8000` (ou a porta definida na variável de ambiente `API_PORT` e mapeada no `docker-compose.yml` para o serviço `backend`).

---
**Próximos Passos (Sugestões)**
* Implementar refresh tokens para melhorar a segurança e a experiência do usuário.
* Adicionar mais testes unitários e de integração para garantir a robustez do código.
* Implementar paginação para a listagem de usuários, especialmente se o número de usuários crescer.
* Adicionar roles (funções) e permissões para usuários, permitindo um controle de acesso mais granular.
* Melhorar o tratamento de erros e logging em toda a aplicação.
* Considerar o uso de migrations mais avançadas para o banco de dados (ex: Goose, GORM's migration tool).