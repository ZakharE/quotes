# Prerequisites

To run application you must have `Docker` installed.

# Used technologies

- Go: 1.21
- Database: `postgreSql`
- Sql migrations: `goose`
- Api handlers and model code generation by api spec: `github.com/deepmap/oapi-codegen`
- Docker, docker-compose

# How to run

## If you have Make installed
Simply run

```shell
make app-start
```

it will run `docker compose` command in detached mode

To stop application

```shell
make app-stop
```

## If you have only docker installed

Run following command from project directory

```shell
docker-compose -f ./deployments/docker-compose.yml  up  --build -d
```

To stop the app run following from project directory

```shell
docker-compose -f ./deployments/docker-compose.yml  up  --build -d
```

# Available endpoints

All endpoint described in `api/quotes/quotes.yaml` file.
For GoLand users there is file with requests located at `requests/quotes.http`

# Refresh logic

Async refresh happens in `quote_refresher` daemon.

The logic is following:

1. Collect all the tasks that are not in `success` and `last_update_attempt` was more than 1 minute ago.
2. Aggregate tasks by currency pair.
3. Request quotes from the client
4. Update task status and data in the database. If the request was unsuccessful, set the tasks an `error` status.









