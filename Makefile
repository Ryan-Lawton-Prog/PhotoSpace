build:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go

build-api:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/api ./cmd/api/main.go

build-app:
	go mod download && go build -o ./.bin/app ./cmd/app/main.go

run-api: build-api
	cd deployments; docker-compose up --build server

test:	
	go test -v ./... -coverprofile=.test/coverage.out

coverage: test
	go tool cover -html=.test/coverage.out
	
dev-api:
	./utils/air -c ./deployments/.air.toml