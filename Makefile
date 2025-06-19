#########
## install
#########

.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

.PHONY: install-ui
install-ui:
	cd ui && bun install

.PHONY: install-playwright
install-playwright:
	cd ui && bun install && bun install:e2e

#########
## fmt
#########

.PHONY: fmt-go
fmt-go:
	go mod tidy
	golangci-lint fmt ./...

.PHONY: fmt-ui
fmt-ui:
	cd ui && bun format

.PHONY: fmt
fmt: fmt-go fmt-ui

#########
## fix
#########

.PHONY: fix-go
fix-go:
	golangci-lint run --fix ./...

.PHONY: fix-ui
fix-ui:
	cd ui && bun lint --fix

.PHONY: fix
fix: fix-go fix-ui

#########
## lint
#########

.PHONY: lint-go
lint-go:
	golangci-lint version
	golangci-lint run ./...

.PHONY: lint-ui
lint-ui:
	cd ui && bun lint --max-warnings 0

.PHONY: lint
lint: lint-go lint-ui

#########
## test
#########

.PHONY: test-go
test-go:
	go test -v ./...

.PHONY: test-ui
test-ui:
	cd ui && bun test src

.PHONY: test-short
test-short: test-go test-ui

.PHONY: test-playwright
test-playwright:
	cd ui && bun test:e2e

.PHONY: test-playwright-ui
test-playwright-ui:
	cd ui && bun test:e2e:ui

.PHONY: test-upgrade-playwright
test-upgrade-playwright:
	./upgradetest/test_all.sh

.PHONY: test
test: test-short test-playwright test-upgrade-playwright

.PHONY: test-coverage
test-coverage:
	go test -coverpkg=./... -coverprofile=coverage.out -count 1 ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

##########
## build
##########

.PHONY: build-ui
build-ui:
	cd ui && bun build-only
	touch ui/dist/.keep

.PHONY: build
build: build-ui
	go build -o catalyst .

.PHONY: build-linux
build-linux: build-ui
	GOOS=linux GOARCH=amd64 go build -o catalyst .

.PHONY: docker
docker: build-linux
	docker build -f docker/Dockerfile -t catalyst .

############
## run
############

.PHONY: reset_data
reset_data:
	rm -rf catalyst_data

.PHONY: copy_existing_data
copy_existing_data: reset_data
	mkdir -p catalyst_data
	cp upgradetest/data/v0.14.1/data.db catalyst_data/data.db

.PHONY: dev
dev: reset_data
	go run . admin create admin@catalyst-soar.com 1234567890
	go run . fake-data
	UI_DEVSERVER=http://localhost:3000 go run . serve --app-url http://localhost:8090 --flags dev

.PHONY: dev_upgrade
dev_upgrade: copy_existing_data
	UI_DEVSERVER=http://localhost:3000 go run . serve --app-url http://localhost:8090 --flags dev

.PHONY: dev-playwright
dev-playwright: reset_data
	go run . admin create admin@catalyst-soar.com 1234567890
	go run . serve --app-url http://localhost:8090 --flags dev

.PHONY: dev-10000
dev-10000: reset_data
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

#########
## generate
#########

.PHONY: sqlc
sqlc:
	rm -rf app/database/sqlc
	cd app/database && go tool sqlc generate
	sed -i.bak 's/Queries/ReadQueries/g' app/database/sqlc/read.sql.go
	rm -f app/database/sqlc/read.sql.go.bak
	sed -i.bak 's/Queries/WriteQueries/g' app/database/sqlc/write.sql.go
	rm -f app/database/sqlc/write.sql.go.bak
	cp app/database/sqlc.db.go.tmpl app/database/sqlc/db.go

.PHONY: openapi-go
openapi-go:
	go tool oapi-codegen --config=app/openapi/config.yml openapi.yml

.PHONY: openapi-ui
openapi-ui:
	rm -rf ui/src/client
	cd ui && bun generate

.PHONY: openapi
openapi: openapi-go openapi-ui

.PHONY: generate-go
generate-go: openapi-go sqlc fmt-go

.PHONY: generate-ui
generate-ui: openapi-ui fmt-ui

.PHONY: generate
generate: generate-go generate-ui