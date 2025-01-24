build:
	@go build -o bin/gobank

startdb:
	docker start gobank

execdb:
	docker exec -it gobank psql -U postgres

run:
	@./bin/gobank

test: 
	@go test ./...
