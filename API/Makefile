build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

run: build
	docker-compose up --build server

test:	
	go test -v ./... -coverprofile=.test/coverage.out

coverage: test
	go tool cover -html=.test/coverage.out
	