run-container:
	docker-compose --env-file ./.env up

run-tests:
	go test -v  ./internal/service ./internal/storage/postgres ./internal/delivery