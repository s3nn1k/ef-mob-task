run-container:
	docker-compose --env-file ./.env up

run-tests:
	go test -v  ./internal/storage/postgres ./internal/service ./internal/delivery