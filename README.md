# maimabank

MaimaBank is a service that provides APIs to enable the frontend to perform the following tasks:

1. Create and manage bank accounts, which consist of the account owner's name, balance, and currency.
2. Record all balance changes made to each account. Every time an amount is added to or subtracted from an account, an account entry record is generated.
3. Transfer funds between two accounts. This happens as a transaction, ensuring that both accounts' balances are either successfully updated or not updated at all.

## Setup local development

### Install tools

Before setting up the service locally, you need to install the following tools:

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
  brew install golang-migrate
  ```

  Windows using Chocolatey

  ```bash
  choco install golang-migrate
  ```

- [DB Docs](https://dbdocs.io/docs)

  ```bash
  npm install -g dbdocs
  dbdocs login
  ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

  ```bash
  npm install -g @dbml/cli
  dbml2sql --version
  ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

  ```bash
  brew install sqlc
  ```

  Windows

  ```bash
  go get github.com/kyleconroy/sqlc/cmd/sqlc
  ```

- [Gomock](https://github.com/golang/mock)

  ```bash
  go install github.com/golang/mock/mockgen@v1.6.0
  ```

### Setup infrastructure

Coming soon...

### How to generate code

Use the following commands to generate code:

- Generate schema SQL file with DBML:

  ```bash
  make db_schema
  ```

- Generate SQL CRUD with sqlc:

  ```bash
  make sqlc
  ```

- Generate DB mock with gomock:

  ```bash
  make mock
  ```

- Create a new db migration:

  ```bash
  migrate create -ext sql -dir db/migration -seq <migration_name>
  ```

### How to run

Use the following commands to run the service:

- Run server:

  ```bash
  make server
  ```

- Run test:

  ```bash
  make test
  ```
