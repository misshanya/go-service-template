# Go Service Template

- Using **Echo** and **PostgreSQL**, with **pgx/v5**, **sqlc** and **goose**
- Dockerfile optimized with build cache

## Structure

```
.
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

4. Enjoy developing!

## Development notes

### Database-related

#### sqlc

- You should write your queries in `internal/db/sqlc/queries/`
- Generate Go code from SQL with `sqlc generate -f internal/db/sqlc/sqlc.yaml`

#### Migrations with goose

- Create migrations with `goose -dir internal/db/migrations create -s {migration_name} sql`

#### Swagger docs

- Generate docs with `swag init -g cmd/main.go`
- Format comments with `swag fmt`
