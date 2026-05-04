.PHONY: run build test migrate swagger clean

run:
	go run cmd/server/main.go

build:
	go build -o server cmd/server/main.go

test:
	go test ./...

migrate:
	mysql -u root -p qiaoyi_community < migrations/001_init.up.sql

swagger:
	swag init -g cmd/server/main.go

clean:
	rm -f server
