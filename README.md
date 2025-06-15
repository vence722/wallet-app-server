# Wallet App Server
A simple and clean wallet app backend server implemented in Go

## User Requirements
According to the user requirement, the following functions are covered:
- User can deposit money into his/her wallet
- User can withdraw money from his/her wallet
- User can send money to another user
- User can check his/her wallet balance
- User can view his/her transaction history

In addition, I think it makes sense that only the authenticated user can be able to access our app, so this implementation will also cover an API endpoint for user to login and get authenticated.

## Project Design

### General Design
The system are mainly consists of three small modules: user, wallet and transaction. 
- user module is responsible for user login, user information maintenance and user activity tracking
- wallet module is responsible for wallet information maintenance, wallet balance checking, deposit/withdraw
- transaction module is responsible for transfer, transaction history query

According to this design, I'll split the system API controllers, services and repositories into separate .go files.

### Key Decision Made for the Design
- How to handle concurrent requests for wallet balance change (deposit, withdraw, transfer)?

    Use Postgres DB transaction. Use SELECT...FOR UPDATE to lock the wallet balance at the begining of the transaction so that another concurrent DB session won't get dirty value. Commit the transaction only when all the update queries are run successfully, otherwise roll back the transaction to recover the state to the beginning of the request, and return error to the client.

- What is the mechanism to authenticate the user to call the APIs?

    Use an access token granted by the /user/login endpoint. After the user successfully logins, an access token will be generated and stored in Redis. The later requests sent to the API server are expected to have a bearer token (in HTTP `Authorization` header) sent together. At the backend, the authentication middleware will verify the access token by parsing the `Authorization` header to obtain the access token and then verify it from Redis. Error will be return if the provided access token cannot be verified.

- How to keep track of all the users and wallets in the system? 

    Two tables are related to the tracking/auditing requirements: `txn_history` and `user_activity`, but they have slightly different purpose. The `txn_history` is mainly used for tracking money related events and targeting a wallet. The `user_activiy`, on the other hand, is used for tracking user events, including money related and non-related events (e.g. login). This separation provides more flexibility for implementing auditing or compliance requirements.

### UML
![](docs/wallet_app_uml.png)

### List of API endpoints

|HTTP Method|Endpoint|Description|
|-|-|-|
|POST|/api/v1/auth/login|User login, and get access token|
|GET|/api/v1/wallet/list|List wallets by user ID|
|POST|/api/v1/wallet/deposit|Deposit to a spcified wallet|
|POST|/api/v1/wallet/withdraw|Withdraw from a specified wallet|
|POST|/api/v1/wallet/checkBalance|Checks wallet balance|
|POST|/api/v1/transaction/transfer|Transfer money from user's wallet to another|
|POST|/api/v1/transaction/history|List transaction history by wallet ID|

The detail API specification can be found in [the OpenAPI spec](api/wallet_app_api_specification.yml)

### Project Structure
For reviewer to begin to review the code, here's a brief introduction about the folder structure of this project:
```
api/ ---------------------> the API specification documents (e.g. OpenAPI/Swagger yaml)
app/ ---------------------> the root of the wallet app source code
    - config/ ------------> app configuration related go files
    - constant/ ----------> global constant shared by all the project
    - controller/ --------> MVC controllers, the entry point of each API endpoints
    - db/ ----------------> DB module, responsible for the database connection
    - entity/ ------------> DB entities to map each DB table, defined in GORM framework standarded
    - logger/ ------------> a logger wrapper to provide an abstract layer for the underlying log library
    - middleware/ --------> custom GIN middlewares
    - model/ -------------> model structs to store data, to be passed through service and controller layers
    - redis/ -------------> Redis module, responsible for the Redis connection
    - repository/ --------> all DB operations defined here, to be called by service layer
    - service/ -----------> all business logic defined here, to be called by controller layer
    - util/ --------------> provides some util functions shared by the project
    - app.go -------------> the entry point of the server, including the logic for initialization and starting the GIN server
    - routes.go ----------> config all the API routes for the server
cmd/ ---------------------> the root of all executable files
    - main.go ------------> the main entry point of the program
database/ ----------------> defines some DB schema sql files
dist/ --------------------> the root of target project directory 
docs/ --------------------> document related items
tests/ -------------------> test related files
    - end2end/ -----------> end-to-end test related files
tools/ -------------------> provide useful executables
    - password_hasher/ ---> a small util to generate password hash used by this project
build_xxx_xxx.sh ---------> build scripts to provide the executable file
```

Basically, the main business logic flow is `controller --> service --> repository`, a traditional [3-tier architecture](https://www.ibm.com/think/topics/three-tier-architecture).

## Installation

(1) Clone this repository

(2) Runs the build script on the top level of the project directory.
- `build_linux_amd64.sh` for running on 64-bit Intel CPU Linux machine (typical for your server)
- `build_mac_arm64.sh` for running on 64-bit Appple CPU MacOS machine (typical for your laptop)
- If none of these are available for your case, simply check and modify the content inside any of the scripts. They're actually doing very simple things - running go build with specified GOOS and GOARCH environment variables and specify the build target to the `dist/` directory

(3) Verify if the `wallet-app-server` executable file is generated in the `dist/` directory

## Configuration
Under the `dist/` directory, you could find `config.toml` file. This is where all the configuration for this server are stored.

- `Server` section contains some basic configuration of the app (e.g. hostname, port, session expire time)
- `Logging` section is responsible for the configuration of the log files
- `DB` section is where you config the database connection
- `Redis` section is where you config the Redis connection

If you want to use our Docker based local testing environment directly, then no need to change the configurations.

## Setup Testing Envirionment
If you have `Docker` installed on your current machine, you could run the following script from your `project root` to spawn a docker-compose with Postgres and Redis:
```
cd tests/end2end
./init_all.sh
```
Make sure you don't have local Postgres or Redis running on the default ports (5432 and 6379) before you spawn the docker-compose.

The script will also insert some test data into the Postgres DB for end-to-end testing.

## Start API Server
If you have finished the `Installation` and `Configuration` steps, go to the `project root` and then run the following commands to start the server

```
cd dist
./wallet-app-server
```
This will run the server in foreground mode. If you want to run it in background mode, run the following commands:

```
cd dist
./start.sh
```

And you can stop the server running in background mode:
```
./stop.sh
```

## Testing

### End-to-end Testing (recommended)
Simply run the following commands:
```
cd tests/end2end
./init_all.sh
./start_test.sh
```

The test scripts are inside `tests/end2end` directory

`init_all.sh` will start a clean docker compose of Postgres and Redis. If existing one is running, the script will tear it down first. Then it'll also insert testing data into DB tables (the data can be found in `test_data.sql`)

`start_test.sh` will trigger the end-to-end test cases to start. It'll log some useful information about the API request/response to help you understand what's happening behind.

### Unit Testing
```
go test ./app/...
```
FACT: This project doesn't include many unit test cases (only for some important stateless functions), since I'm still finding a effective way to create unit test cases for the business logic part, which has a lot of dependencies that need to create mock up objects.

## Area of improvements
- More unit testing to cover all important functions
- Run load testing and do performance optimization
- List endpoint for transaction history should have pagination
- Support K8S deployment

## Features wishlist
- User profile maintenance
- User wallets information maintenance
- Endpoint to list user activities (login, deposit, withdraw, transfer, etc.)

## Time spent on the test
<=72 hours