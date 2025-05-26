FROM arm64v8/golang

WORKDIR /app/

COPY ./.bin/api .
COPY ./config/docker-config.yml ./config/config.yml

COPY ./cmd/api/main.go ./cmd/api/main.go
COPY ./config/init.go ./config/init.go
COPY ./internal/api ./internal/api
COPY ./internal/pkg ./internal/pkg

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go
EXPOSE 8000
CMD [ "api" ]
