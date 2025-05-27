build:
	cd backend && go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go
	cd frontend && npm i && npm run build

build-api:
	cd backend && go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go

build-ui:
	cd frontend && npm run build

run-api: build-api
	docker-compose -f "docker-compose.yml" up --build api

run-ui:
	serve -p 3000 -d ./frontend/dist

container:
	docker-compose -f "docker-compose.yml" up --build

db:
	docker-compose -f "docker-compose.yml" up -d db

test:	
	cd backend && go test -v ./... -coverprofile=./test/.report/coverage.out

coverage: test
	go tool cover -html=./test/.report/coverage.out
	
dev-api:
	./utils/air -c ./backend/deployments/.air.toml

dev-ui:
	cd internal/web && npm run dev

dev-container:
	docker-compose watch