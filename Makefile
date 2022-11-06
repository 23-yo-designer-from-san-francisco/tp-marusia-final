.PHONY: server
server:
	go build -tags=jsoniter -o bin/api/marusia -v ./cmd/server

.PHONY: swagger
swagger:
	swag init -g ./cmd/server/main.go
