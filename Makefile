dev:
	go run cmd/poc-misreported-qty-api/main.go

db:
	docker run --name poc -d -p 5432:5432 -e POSTGRES_USER=poc -e POSTGRES_PASSWORD=poc-123 -e POSTGRES_DB=poc postgres

db-test:
	docker run --name poc-test -d -p 5437:5432 -e POSTGRES_USER=poc -e POSTGRES_PASSWORD=poc-123 -e POSTGRES_DB=poc-test postgres

build:
	rm -f engine
	go build -o engine cmd/poc-misreported-qty-api/main.go

test:
	go test ./cmd/poc-misreported-qty-api -v