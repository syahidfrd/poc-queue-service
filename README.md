# POC Misreported Quantity

### Objective
We are members of the engineering team of an online store. When we look at ratings for our online store application, we received the following
facts:
- Customers were able to put items in their cart, check out, and then pay. After several days, many of our customers received calls from
our Customer Service department stating that their orders have been canceled due to stock unavailability.
- These bad reviews generally come within a week after our 12.12 event, in which we held a large flash sale and set up other major
discounts to promote our store.
  
After checking in with our Customer Service and Order Processing departments, we received the following additional facts:
- Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
- The misreported items are those that performed very well on our 12.12 event.
- Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus
requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.
  
### An explanation of why it happened and the solution
This usually occurs because the system processes orders asynchronously, i.e. if there are 2 or more customers making a product order at the same time, the system will perform all of these tasks simultaneously. This results in an incorrectly reported product quantity, as we must always check the product stock and ensure its availability. The solution is to use a queuing system with the FIFO concept, with this concept processes that previously occurred asynchronously become sequential and take turns.

!["Queue Flow"](https://i.ibb.co/rF6nBzH/Untitled-Diagram-1.png "Queue Flow")

### Requirements

- Golang = 1.14
- PostgreSQL = 12
- Docker Use [Docker CE](https://docs.docker.com/engine/installation) for Linux or [Docker Toolbox](https://www.docker.com/products/docker-toolbox) for Windows and Mac.
- RabbitMQ = 3-management
- Go modules is a collection of Go packages stored in a file tree with a `go.mod` file at its root.

### Setting up Project

- Extract this repository.
- Install all project dependencies.
- Run `make test` to run tests and make sure that all tests are passing.

### List of Useful Commands
- `make db` -> Run PostgreSQL as docker container
- `make rabbitmq` -> Run RabbitMQ as docker container
- `make test` -> Run test to ensure all functionalities are work
- `make consumer` -> Run consumer for order product queue
- `make dev` -> Run server
- `make build` -> Build poc api binary file
- You can find another useful commands in `Makefile`.
