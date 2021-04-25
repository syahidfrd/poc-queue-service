package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/namsral/flag"
	"log"
	"net/url"
	"poc-misreported-qty/server"
	"poc-misreported-qty/sql"
	"poc-misreported-qty/util/queue"
	"poc-misreported-qty/util/validator"
	"time"
)

var (
	port                = flag.String("port", "8080", "Port where the app listen to")
	debug               = flag.Bool("debug", true, "Run app in debug mode")
	dbConnectionString  = flag.String("db-connection-string", "postgres://poc:poc-123@localhost/poc?sslmode=disable", "Database connection string")
	dbMaxOpenConnection = flag.Int("db-max-open-connection", 1, "Max database open connection")
	dbMaxIdleConnection = flag.Int("db-max-idle-connection", 1, "Max database idle connection")
	amqpServerURL       = flag.String("amqp-server-url", "amqp://guest:guest@localhost:5672", "AMQP server URL")
)

func main() {
	flag.Parse()

	parseDBUrl, _ := url.Parse(*dbConnectionString)

	gormInstance, err := gorm.Open(parseDBUrl.Scheme, *dbConnectionString)
	if err != nil {
		panic("Failed to connect database " + err.Error())
	}

	defer gormInstance.Close()

	gormInstance.DB().SetConnMaxLifetime(60 * time.Minute)
	gormInstance.DB().SetMaxOpenConns(*dbMaxOpenConnection)
	gormInstance.DB().SetMaxIdleConns(*dbMaxIdleConnection)

	var (
		dataRepository   = sql.NewSQLDataRepository(gormInstance)
		validatorService = validator.NewDefaultValidatorService()
		queueService     = queue.NewDefaultQueueService(*amqpServerURL)
	)

	var defaultServer = &server.DefaultAPIServer{
		Config: &server.Config{
			ListenPort: *port,
			Debug:      *debug,
		},
		APIV1Config: &server.APIV1Config{
			DataRepository:   dataRepository,
			ValidatorService: validatorService,
			QueueService:     queueService,
		},
	}

	defaultServer.InitEngine()
	defaultServer.RegisterRoutes()

	var serverInstance server.Service = defaultServer
	err = serverInstance.Run()

	if err != nil {
		log.Panic(err)
	}
}
