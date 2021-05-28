package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/namsral/flag"
	"github.com/streadway/amqp"
	"net/url"
	"poc-queue-service/model"
	"poc-queue-service/sql"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		amqpServerURL       = flag.String("amqp-server-url", "amqp://guest:guest@rabbitmq:5672", "AMQP server URL")
		dbConnectionString  = flag.String("db-connection-string", "postgres://poc:poc-123@database/poc?sslmode=disable", "Database connection string")
		dbMaxOpenConnection = flag.Int("db-max-open-connection", 1, "Max database open connection")
		dbMaxIdleConnection = flag.Int("db-max-idle-connection", 1, "Max database idle connection")
	)

	parseDBUrl, _ := url.Parse(*dbConnectionString)

	gormInstance, err := gorm.Open(parseDBUrl.Scheme, *dbConnectionString)
	if err != nil {
		panic("Failed to connect database " + err.Error())
	}

	defer gormInstance.Close()

	gormInstance.DB().SetConnMaxLifetime(60 * time.Minute)
	gormInstance.DB().SetMaxOpenConns(*dbMaxOpenConnection)
	gormInstance.DB().SetMaxIdleConns(*dbMaxIdleConnection)

	var dataRepository   = sql.NewSQLDataRepository(gormInstance)

	conn, err := amqp.Dial(*amqpServerURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	messages, err := ch.Consume(
		"order_product",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			fmt.Printf("Received message: %s\n", d.Body)
			splitMessage := strings.Split(string(d.Body), ":")
			productID := splitMessage[0]
			quantity := splitMessage[1]

			product, _ := dataRepository.ProductStore.FindOneBy(map[string]interface{}{
				"id": productID,
			})

			if !dataRepository.ProductStore.Exist(product) {
				fmt.Printf("Product id=%s not found", productID)
				continue
			}

			i, _ := strconv.Atoi(quantity)

			// Validate quantity
			if product.GetQuantity() < uint32(i) {
				fmt.Println("Invalid quantity")
				continue
			}

			// Create new order
			order := model.NewOrder(uint32(i), product.GetID())
			if err := dataRepository.OrderStore.Save(order); err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Printf("Create order id:%d success\n", order.GetID())
		}
	}()

	fmt.Println("Successfully connected to RabbitMQ instance")
	<-forever
}
