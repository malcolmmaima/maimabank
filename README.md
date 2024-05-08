# Maima Bank

MaimaBank is a service that provides APIs to enable the frontend to perform the following tasks:

1. Create and manage bank accounts, which consist of the account owner's name, balance, and currency.
2. Record all balance changes made to each account. Every time an amount is added to or subtracted from an account, an account entry record is generated.
3. Transfer funds between two accounts. This happens as a transaction, ensuring that both accounts' balances are either successfully updated or not updated at all.
4. Generate account statements for a given account, which consist of a list of account transfers.
5. Add support for multiple currency transfers. Each account can only hold one currency, and transfers between accounts with different currencies will be automatically converted using the latest exchange rate.

#### To-do
6. Add support for loan accounts. Loan accounts are accounts that can have a negative balance. When a loan account is created, a loan limit is set. The account balance can go down to the negative of the loan limit. For example, if the loan limit is 1000, the account balance can go down to -1000. When a loan account is created, the account balance is set to 0.
7. Add support for loan repayments. Loan accounts can be topped up by transferring funds from another account. The transfer amount is added to the loan account balance. For example, if the loan account balance is -500 and the transfer amount is 200, the new loan account balance will be -300.
8. Add support for loan repayments with interest. Loan accounts can be topped up by transferring funds from another account. The transfer amount is added to the loan account balance, and an interest amount is added to the transfer amount. For example, if the loan account balance is -500, the transfer amount is 200, and the interest rate is 10%, the new loan account balance will be -280.
9. Add support for loan repayments with interest and fees. Loan accounts can be topped up by transferring funds from another account. The transfer amount is added to the loan account balance, an interest amount is added to the transfer amount, and a fee is deducted from the transfer amount. For example, if the loan account balance is -500, the transfer amount is 200, the interest rate is 10%, and the fee is 5, the new loan account balance will be -285.
10. Add support for creating and managing savings accounts. Savings accounts are accounts that can only hold one currency and cannot have a negative balance. When a savings account is created, an interest rate is set. The interest rate is used to calculate the interest amount that is added to the account balance at the end of each month. For example, if the interest rate is 10% and the account balance is 1000, the interest amount added to the account balance at the end of the month will be 100.
11. Add support for creating and managing savings accounts with fees. Savings accounts are accounts that can only hold one currency and cannot have a negative balance. When a savings account is created, an interest rate is set. The interest rate is used to calculate the interest amount that is added to the account balance at the end of each month. For example, if the interest rate is 10% and the account balance is 1000, the interest amount added to the account balance at the end of the month will be 100. A fee is deducted from the account balance at the end of each month. For example, if the fee is 5, the account balance will be deducted by 5 at the end of each month.
12. Add Mpesa integration. Mpesa is a mobile money service that allows users to send and receive money. Mpesa integration allows users to transfer funds from their Mpesa account to their bank account and vice versa.
13. Security Upgrade. Add support for two-factor authentication, security questions, and other security features.
14. Do let me know of other features you would like to see in the service by raising an issue.

Last Step: Deploy to production and build an Android/Web app to consume the Bank API.

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

After installing the above tools you can proceed to setup your environment with below steps:

1. Run the following command to setup postgres:

   ```bash
   make postgres
   ```

2. Create a new database:

   ```bash
   make createdb
   ```

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
