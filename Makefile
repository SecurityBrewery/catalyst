.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

.PHONY: fmt-go
fmt-go:
	@echo "Formatting..."
	go mod tidy
	golangci-lint fmt ./...

.PHONY: fmt-ui
fmt-ui:
	@echo "Formatting..."
	cd ui && bun format

.PHONY: fmt
fmt: fmt-go fmt-ui

.PHONY: fix
fix:
	@echo "Fixing..."
	golangci-lint run --fix ./...

.PHONY: lint-go
lint-go:
	mkdir -p ui/dist
	touch ui/dist/index.html
	golangci-lint version
	golangci-lint run ./...

.PHONY: lint-ui
lint-ui:
	@echo "Linting..."
	cd ui && bun lint

.PHONY: lint
lint: lint-go lint-ui

.PHONY: test-go
test-go:
	@echo "Testing..."
	mkdir -p ui/dist
	touch ui/dist/index.html
	go test -v ./...

.PHONY: test-ui
test-ui:
	@echo "Testing..."
	cd ui && bun test

.PHONY: test
test: test-go test-ui

.PHONY: test-coverage
test-coverage:
	@echo "Testing with coverage..."
	mkdir -p ui/dist
	touch ui/dist/index.html
	go test -coverpkg=./... -coverprofile=coverage.out -count 1 ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: build-ui
build-ui:
	@echo "Building..."
	cd ui && bun install
	cd ui && bun build-only

.PHONY: build
build: build-ui
	@echo "Building..."
	go build -o catalyst .


.PHONY: build-linux
build-linux: build-ui
	@echo "Building..."
	GOOS=linux GOARCH=amd64 go build -o catalyst .

.PHONY: docker
docker: build-linux
	@echo "Building Docker image..."
	docker build -f docker/Dockerfile -t catalyst .

.PHONY: dev
dev:
	@echo "Running..."
	rm -rf catalyst_data
	go run . admin create admin@catalyst-soar.com 1234567890
	go run . fake-data
	go run . serve --app-url http://localhost:8090 --flags dev


.PHONY: dev2
dev2:
	@echo "Running..."
	rm -rf catalyst_data
	mkdir -p catalyst_data
	cp upgradetest/data/v0.14.1/data.db catalyst_data/data.db
	go run .

.PHONY: dev-10000
dev-10000:
	@echo "Running..."
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
	@echo "TailwindCSS..."
	# tailwindcss -i ./css/in.css -o ./static/output.css --cwd api/ui
	@echo "Done."

.PHONY: sqlc
sqlc:
	@echo "SQLC..."
	cd app/database && go tool sqlc generate
	@echo "Done."

.PHONY: openapi
openapi:
	@echo "OpenAPI..."
	cd app && go tool oapi-codegen --config=openapi/config.yml openapi/openapi.yml
	openapi-generator generate -i app/openapi/openapi.yml -g typescript-fetch -o ui/src/client
	@echo "Done."

.PHONY: generate
generate: sqlc tailwindcss openapi