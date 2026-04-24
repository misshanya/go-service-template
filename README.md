# Go Service Template

## Tech stack

- [Echo](https://echo.labstack.com) with oapi-codegen
- PostgreSQL with [pgx/v5](https://github.com/jackc/pgx), [sqlc](https://sqlc.dev) and [goose](https://github.com/pressly/goose)
- S3 with [AWS SDK](https://github.com/aws/aws-sdk-go-v2) ([RustFS](https://rustfs.com) included in compose)
- Valkey
- Kafka with [segmentio/kafka-go](https://github.com/segmentio/kafka-go)

## Structure

```
.
├── api/
│   └── openapi.yaml - OpenAPI spec for the service, used for code generation and docs
├── cmd/
│   └── main.go - entry point for the whole application, loads a config, creates an app and manages its lifecycle
├── internal
│   ├── app/
│   │   └── app.go - initialization of all components
│   ├── config/
│   │   └── config.go - config structs and loading from .env
│   ├── db/
│   │   ├── migrate.go - simple function to migrate db
│   │   ├── migrations/ - there would be your migrations
│   │   └── sqlc/ - sqlc config and queries
│   ├── errorz/ - domain-level errors
│   ├── models/ - domain-level models
│   ├── repository/ - storage, repo level
│   ├── service/ - use-case, business-logic and whatever else it's called
│   └── transport/ - top-level transport communications
│       └── http/ - specifically http transport
├── pkg/
│   └── dberrors/ - some sugar for db errors processing
│       └── is_unique_violation.go
```

## How to start

1. Clone repo and cd into it

    ```shell
   git clone https://github.com/misshanya/go-service-template my-ultimate-project
   cd my-ultimate-project
    ```

2. Remove `.git/` to use it as a base for your own project

    ```shell
   rm -rf .git
    ```
   
3. Update the module name in go.mod (IDE can help refactor the code accordingly)

4. Enjoy developing! (don't forget to read dev notes)

## Development notes

### Infrastructure

- Use required connectors in `internal/app/app.go` and remove unused ones  
  Uncomment or remove corresponding services in compose

### Database-related

#### sqlc

- You should write your queries in `internal/db/sqlc/queries/`
- Generate Go code from SQL with `sqlc generate -f internal/db/sqlc/sqlc.yaml`

#### Migrations with goose

- Create migrations with `goose -dir internal/db/migrations create -s {migration_name} sql`

### API

- Swagger runs on `/api/v1/swagger`

#### OpenAPI codegen

- Generate code from OpenAPI spec with `oapi-codegen -config oapi-codegen.yml api/openapi.yaml`
