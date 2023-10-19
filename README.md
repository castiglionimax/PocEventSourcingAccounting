# POC

This project involves building a REST API using the Event Sourcing architecture, implemented in Go programming language. The API utilizes Event Sourcing principles, where events are stored in a MongoDB event store. Additionally, projections are created using MySQL. The entire project is containerized using Docker Compose, making it easy to deploy and manage. This setup ensures a scalable and efficient system for handling events and projections, with MongoDB serving as the event store and MySQL for data projections. The Go programming language is used for the implementation, providing a robust foundation for the API
## Table of Contents
- [Installation](#installation)
- [Usage](#usage)

## installation

This project uses Docker Compose for easy setup and deployment. Make sure you have Docker and Docker Compose installed on your system. Follow these steps to get started:

1. Clone the repository:
```bash
   git clone https://github.com/castiglionimax/PocEventSourcingAccounting.git
   ```

2. Navigate to the project directory:

```bash
cd PocEventSourcingAccounting
   ```
3. Build and start the project containers:

```bash
docker-compose up --build
   ```
## Usage

To test the application, you can utilize Postman or a similar tool. Below are the curl commands:

To create a new account:
```sh
curl --location --request POST 'http://127.0.0.1:8080/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "juan",
    "account_number": 123
}'
`````
With the account ID obtained, make a deposit
```sh
curl --location --request POST 'http://127.0.0.1:8080/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
"account_id": "bf08ebb5-b470-490e-9b94-192b0e560dd3",
"transaction_type": "deposit",
"amount": 1000
}'
```

Or make a cash withdrawal.
```sh
curl --location --request POST 'http://127.0.0.1:8080/transactions' \
--header 'Content-Type: application/json' \
--data-raw '{
"account_id": "bf08ebb5-b470-490e-9b94-192b0e560dd3",
"transaction_type": "withdrawal",
"amount": 500
}'
```

Finally, to obtain an account balance.
```sh
curl --location --request GET 'http://127.0.0.1:8080/accounts/:account_id/balance'
```

To stop the project containers, you can run:
```bash
docker-compose down
   ```

