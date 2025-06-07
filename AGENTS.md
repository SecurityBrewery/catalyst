# Contributor Guide

## Dev Environment Tips
- Go Backend:
    - The sql migrations are located in the `app/database/migrations` directory, queries in `app/database/query.sql`.
    - Never change existing migrations, add new ones when needed.
    - `sqlc` is used to generate Go code from SQL queries. Make sure to run `make sqlc` after modifying any SQL files.
    - Generated Go code is located in `app/database/sqlc` and should not be modified directly.
    - The OpenAPI spec is located at `app/openapi/openapi.yaml`. This file defines the API endpoints and their schemas.
    - `oapi-codegen` and `openapi-generator` are used to generate Go code from OpenAPI specifications. Run `make openapi` after modifying any OpenAPI files.
    - The generated Go code from OpenAPI is located in `app/openapi/gen.go` and should not be modified directly.
- Vue Frontend:
    - The frontend is located in the `/ui` folder.
    - Use `make install-ui` to install the necessary dependencies.
    - Use `make build-ui` to build the frontend.
    - Use `make dev-ui` to start the frontend in development mode.
    - The OpenAPI spec is also used to generate TypeScript types for the frontend. Run `make openapi` after modifying any OpenAPI files.
    - The generated TypeScript types are located in `ui/src/client` and should not be modified directly.

## Testing Instructions
- Go Backend:
    - Use `make fmt-go` to format the Go codebase.
    - Use `make test-go` to run all Go tests in the workspace.
    - Use `make lint-go` to run all Go linters in the workspace.
    - Format the codebase before running linters and tests.
    - Fix any test or linter errors until the whole suite is green.
    - Add or update tests for the code you change, even if nobody asked.
- Vue Frontend:
    - Use `make fmt-ui` to format the ui codebase.
    - Use `make test-ui` to run all ui tests in the workspace.
    - Use `make lint-ui` to run all ui linters in the workspace.
- Playwright End-to-End Tests:
    - Playwright tests are located in the `/playwright` directory.
    - Use `make install-ui build-ui install-playwright` to install the necessary dependencies for Playwright tests.
    - Use `make test-playwright` to run the Playwright end-to-end tests.

## PR instructions
- Use semantic commit messages like `feat: add new feature` or `fix: correct a bug`.