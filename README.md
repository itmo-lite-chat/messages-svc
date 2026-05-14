# messages-svc

gRPC service for lite-chat messages.

## Run

```bash
env $(cat dev/.env | xargs) go run ./cmd
```

## New migration
```go
goose -dir ./internal/app/migrations create $YOUR_MIGRATION_NAME sql
```
