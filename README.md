# POC Queue Service - Online Store

### Context
Customers were able to put items in their cart, check out, and then pay. After several days, many of our customers received calls from
our Customer Service department stating that their orders have been canceled due to stock unavailability.

### Additional Fact
- Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
- Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus
requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.

### An explanation of why it happened and the solution
This usually occurs because the system processes orders asynchronously, i.e. if there are 2 or more customers making a product order at the same time, the system will perform all of these tasks simultaneously. This results in an incorrectly reported product quantity, as we must always check the product stock and ensure its availability. The solution is to use a queuing system with the FIFO concept, with this concept processes that previously occurred asynchronously become sequential and take turns.

!["Queue Flow"](https://i.ibb.co/rF6nBzH/Untitled-Diagram-1.png "Queue Flow")

### Requirements
- Go 1.14
- Docker and Docker Compose use [Docker CE](https://docs.docker.com/engine/installation) for Linux or [Docker Toolbox](https://www.docker.com/products/docker-toolbox) for Windows and Mac.

### Setting up Project

- Extract this repository.
- Install all project dependencies `go mod tidy`

### How to run test
- Run the db test `make db-test`
- Run the rabbitmq test `make rabbitmq-test`
- Run the test `make test`

### How to run app as docker container
- Run the docker compose
```bash
docker-compose up -d --build
```
- App running on `localhost:8080`
- RabbitMQ management console `localhost:15672`, username/pass: `guest`

### API Spec
- Create product
```bash
path: /api/v1/product/create
body:
{
  name: string,
  quantity: int,
  price: int
}
```
- Get all product
```bash
path: /api/v1/product
```  

- Create order
```bash
path: /api/v1/order/create
body:
{
  product_id: int,
  quantity: int
}
```

- Get all order
```bash
path: /api/v1/order
```

### List of Useful Commands
- `make db` -> Run PostgreSQL as docker container
- `make rabbitmq` -> Run RabbitMQ as docker container
- `make test` -> Run test to ensure all functionalities are work
- `make consumer` -> Run consumer for order product queue
- `make dev` -> Run server
- `make build` -> Build poc api binary file
- You can find another useful commands in `Makefile`.
