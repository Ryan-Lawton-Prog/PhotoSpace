build-api:
	cd API; go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go

run-api: build-api
	cd API; docker-compose up --build server

test-api:	
	cd API; go test -v ./... -coverprofile=.test/coverage.out

coverage-api: test-api
	cd API; go tool cover -html=.test/coverage.out
	
dev-api:
	cd API; ../utils/air -c .air.toml