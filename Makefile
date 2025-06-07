.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

.PHONY: fmt-go
fmt-go:
	go mod tidy
	golangci-lint fmt ./...

.PHONY: fmt-ui
fmt-ui:
	cd ui && bun format

.PHONY: fmt
fmt: fmt-go fmt-ui

.PHONY: fix
fix:
	golangci-lint run --fix ./...

.PHONY: lint-go
lint-go:
	golangci-lint version
	golangci-lint run ./...

.PHONY: lint-ui
lint-ui:
	cd ui && bun lint

.PHONY: lint
lint: lint-go lint-ui

.PHONY: test-go
test-go:
	go test -v ./...

.PHONY: test-ui
test-ui:
	cd ui && bun test src

.PHONY: test-playwright
test-playwright:
	cd playwright && bun test:e2e

.PHONY: test
test: test-go test-ui test-playwright

.PHONY: test-coverage
test-coverage:
	go test -coverpkg=./... -coverprofile=coverage.out -count 1 ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: install-ui
install-ui:
	cd ui && bun install

.PHONY: install-playwright
install-playwright:
	cd playwright && bun install && bun install:e2e
	
.PHONY: build-ui
build-ui:
	cd ui && bun build-only

.PHONY: build
build: build-ui
	go build -o catalyst .

.PHONY: build-linux
build-linux: build-ui
	GOOS=linux GOARCH=amd64 go build -o catalyst .

.PHONY: docker
docker: build-linux
	docker build -f docker/Dockerfile -t catalyst .

.PHONY: dev
dev:
	rm -rf catalyst_data
	go run . admin create admin@catalyst-soar.com 1234567890
	go run . fake-data
	go run . serve --app-url http://localhost:8090 --flags dev

.PHONY: dev2
dev2:
	rm -rf catalyst_data
	mkdir -p catalyst_data
	cp upgradetest/data/v0.14.1/data.db catalyst_data/data.db
	UI_DEVSERVER=http://localhost:3000 go run .

.PHONY: dev-playwright
dev-playwright:
	rm -rf catalyst_data
	mkdir -p catalyst_data
	cp upgradetest/data/v0.14.1/data.db catalyst_data/data.db
	go run .

.PHONY: dev-10000
dev-10000:
	rm -rf catalyst_data
	go run . admin create admin@catalyst-soar.com 1234567890
	go run . fake-data --users 100 --tickets 10000
	go run . serve --app-url http://localhost:8090 --flags dev

.PHONY: default-data
default-data:
	rm -rf catalyst_data
	go run . default-data

.PHONY: serve-ui
serve-ui:
	cd ui && bun dev --port 3000

.PHONY: tailwindcss
tailwindcss:
	# tailwindcss -i ./css/in.css -o ./static/output.css --cwd api/ui

.PHONY: sqlc
sqlc:
	cd app/database && go tool sqlc generate

.PHONY: openapi-go
openapi-go:
	cd app && go tool oapi-codegen --config=openapi/config.yml openapi/openapi.yml
	rm -rf ui/src/client

.PHONY: openapi-ui
openapi-ui:
	cd ui && bun generate

.PHONY: openapi
openapi: openapi-go openapi-ui

.PHONY: generate-go
generate-go: openapi-go sqlc fmt-go

.PHONY: generate-ui
generate-ui: openapi-ui tailwindcss fmt-ui

.PHONY: generate
generate: generate-go generate-ui