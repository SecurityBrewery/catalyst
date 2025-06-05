# Contributor Guide

## Dev Environment Tips
- The sql migrations are located in the `app/database/migrations` directory, queries in `app/database/query.sql`.
- `sqlc` is used to generate Go code from SQL queries. Make sure to run `make sqlc` after modifying any SQL files.
- Generated Go code is located in `app/database/sqlc` and should not be modified directly.
- The OpenAPI spec is located at `app/openapi/openapi.yaml`. This file defines the API endpoints and their schemas.
- `oapi-codegen` and `openapi-generator` are used to generate Go code from OpenAPI specifications. Run `make openapi` after modifying any OpenAPI files.
- The generated Go code from OpenAPI is located in `app/openapi/gen.go` and should not be modified directly.

## Testing Instructions
- Use `make test` to run all tests in the workspace.
- Use `make lint` to run all linters in the workspace.
- Fix any test or linter errors until the whole suite is green.
- Add or update tests for the code you change, even if nobody asked.

## PR instructions
- Use semantic commit messages like `feat: add new feature` or `fix: correct a bug`.