package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	v1 "poc-misreported-qty/api/v1"
	"poc-misreported-qty/model"
	"poc-misreported-qty/util/queue"
	"poc-misreported-qty/util/validator"
	"time"
)

type (
	Config struct {
		ListenPort string
		Debug      bool
		Test       bool
	}

	APIV1Config struct {
		ValidatorService validator.Service
		QueueService     queue.Service
		DataRepository   model.DataRepository
	}

	DefaultAPIServer struct {
		Test        bool
		Config      *Config
		APIV1Config *APIV1Config
		Engine      *gin.Engine
	}
)

func (d *DefaultAPIServer) InitEngine() {
	if d.Config.Test {
		gin.SetMode("test")
	} else if d.Config.Debug {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	d.Engine = gin.New()

	d.Engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
		MaxAge:          12 * time.Hour,
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	d.Engine.Use(analyticMiddleware(d.Config.Test, []string{}))
}

func initProductHandler(d *DefaultAPIServer) (productHandler *v1.ProductHandler) {
	productHandler = v1.NewProductHandler(
		d.APIV1Config.ValidatorService,
		d.APIV1Config.DataRepository.ProductStore,
	)
	return
}

func initOrderHandler(d *DefaultAPIServer) (orderHandler *v1.OrderHandler) {
	orderHandler = v1.NewOrderHandler(
		d.APIV1Config.ValidatorService,
		d.APIV1Config.QueueService,
		d.APIV1Config.DataRepository.ProductStore,
		d.APIV1Config.DataRepository.OrderStore,
	)
	return
}

func (d *DefaultAPIServer) RegisterRoutes() {
	d.Engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "api up and running",
		})
	})

	var (
		productHandler = initProductHandler(d)
	)

	apiV1Public := d.Engine.Group("/api/v1")
	apiV1Public.POST("/product/create", productHandler.CreateProduct)
}

func (d *DefaultAPIServer) Run() (err error) {
	listenAddress := fmt.Sprintf(":%s", d.Config.ListenPort)
	err = d.Engine.Run(listenAddress)
	return
}
