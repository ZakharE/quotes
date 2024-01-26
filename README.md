# Prerequisites

You should have `Docker` and `make` installed.

# Used technologies

- Go: 1.21
- Database: `postgreSql`
- Sql migrations: `goose`
- Api handlers and model code generation by api spec: `github.com/deepmap/oapi-codegen`
- Docker, docker-compose

# How to run

Simply run

```shell
make app-start
```

it will run `docker compose` command in detached mode

To stop application

```shell
make app-stop
```

# Available endpoints

All endpoint described in `api/quotes/quotes.yaml` file.
For GoLand users there is file with requests located at `requests/quotes.http`

# App design

Async refresh happens in `quote_refresher` daemon.

The logic is following:

Collect all tasks in not `success` status within the batch.

Then, aggregate tasks by currency pair.

Next, request quotes from the client and save them in the database.

If the request was unsuccessful, move the tasks to an `error` status.







