dev:
	go run cmd/poc-queue-service/main.go

consumer:
	go run cmd/order-product-consumer/main.go

db:
	docker run --name poc -d -p 5432:5432 -e POSTGRES_USER=poc -e POSTGRES_PASSWORD=poc-123 -e POSTGRES_DB=poc postgres

db-test:
	docker run --name poc-test -d -p 5437:5432 -e POSTGRES_USER=poc -e POSTGRES_PASSWORD=poc-123 -e POSTGRES_DB=poc-test postgres

rabbitmq-test:
	docker run -d --hostname my-rabbit-test --name some-rabbit-test -p 15673:15672 -p 5673:5672 rabbitmq:3-management

rabbitmq:
	docker run -d --hostname my-rabbit --name some-rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

build:
	rm -f engine
	go build -o engine cmd/poc-queue-service/main.go

test:
	go test ./cmd/poc-queue-service -v