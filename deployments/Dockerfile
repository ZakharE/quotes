FROM golang:1.21-alpine

WORKDIR /app

COPY ../go.mod go.sum ./

# Download dependencies.
RUN apk add --no-cache make
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.15.0
RUN go mod download


COPY .. .

RUN make generate-all
RUN GOOS=linux go build -o quotes cmd/quotes/main.go
