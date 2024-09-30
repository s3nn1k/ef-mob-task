run-container:
	docker-compose --env-file ./.env up

run-tests:
	go test -v  ./internal/storage/postgres ./internal/service ./internal/delivery

create-swagger:
	swag init -g internal/delivery/handler.go
	swag init -g cmd/main.go