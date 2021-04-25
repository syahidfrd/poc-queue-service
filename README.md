# POC Misreported Quantity
### Requirements

- Golang = 1.14
- PostgreSQL = 12
- Docker Use [Docker CE](https://docs.docker.com/engine/installation) for Linux or [Docker Toolbox](https://www.docker.com/products/docker-toolbox) for Windows and Mac.
- Go modules is a collection of Go packages stored in a file tree with a `go.mod` file at its root.

### Setting up Project

- Extract this repository.
- Install all project dependencies.
- Run `make test` to run tests and make sure that all tests are passing.

### How to Run

- Service requires PostgreSQL.

```bash
# run postgres with docker
docker run --name poc -d -p 5432:5432 -e POSTGRES_USER=poc -e POSTGRES_PASSWORD=poc-123 -e POSTGRES_DB=poc postgres
# run postgres with make command
make db
```

### List of Useful Commands

- `make test` -> Run test to ensure all functionalities are work
- `make build` -> Build poc api binary file
- You can find another useful commands in `Makefile`.
