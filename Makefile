build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go
	cd internal/web && npm run build

build-api:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go


build-ui:
	cd internal/web && npm run build

run-api: build-api
	docker-compose -f "docker-compose.yml" up --build server

run-ui:
	serve -p 3000 -d ./internal/web/dist

test:	
	go test -v ./... -coverprofile=./test/.report/coverage.out

coverage: test
	go tool cover -html=./test/.report/coverage.out
	
dev-api:
	./utils/air -c ./deployments/.air.toml

dev-ui:
	cd internal/web && npm run dev