package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/namsral/flag"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"poc-misreported-qty/model"
	"poc-misreported-qty/server"
	"poc-misreported-qty/sql"
	"poc-misreported-qty/util/queue"
	"poc-misreported-qty/util/validator"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	appURL    string
	db        *gorm.DB
	refreshDB = func() {
		res := db.Exec("drop schema public cascade; create schema public;")
		if res.Error != nil {
			log.Panicf("Failed to clear database. Err: %v", res.Error)
		}

		// Migrate database
		db.AutoMigrate(
			&model.Product{},
			&model.Order{},
		)
	}

	apiV1Index         string
	apiV1CreateProduct string
	apiV1GetAllProduct string
	apiV1CreateOrder   string
	apiV1GetAllOrder   string
)

func TestMain(m *testing.M) {
	flag.Parse()

	*debug = false

	var (
		gormInstance, _ = gorm.Open("postgres", "postgres://poc:poc-123@localhost:5437/poc-test?sslmode=disable")

		dataRepository   = sql.NewSQLDataRepository(gormInstance)
		validatorService = validator.NewDefaultValidatorService()
		queueService     = queue.NewDefaultQueueService("amqp://guest:guest@localhost:5673")
	)

	// set db global variable
	db = gormInstance

	// create server instance
	serverInstance := &server.DefaultAPIServer{
		Test: true,
		Config: &server.Config{
			Test: true,
		},
		APIV1Config: &server.APIV1Config{
			ValidatorService: validatorService,
			DataRepository:   dataRepository,
			QueueService:     queueService,
		},
	}

	serverInstance.InitEngine()
	serverInstance.RegisterRoutes()

	// run test server
	s := httptest.NewServer(serverInstance.Engine)
	appURL = s.URL

	// set endpoint
	apiV1Index = appURL
	apiV1CreateProduct = fmt.Sprintf("%s/api/v1/product/create", appURL)
	apiV1GetAllProduct = fmt.Sprintf("%s/api/v1/product", appURL)
	apiV1CreateOrder = fmt.Sprintf("%s/api/v1/order/create", appURL)
	apiV1GetAllOrder = fmt.Sprintf("%s/api/v1/order", appURL)

	exitCode := m.Run()
	s.Close()
	os.Exit(exitCode)
}

func TestIndex(t *testing.T) {
	Convey("Given symomath API instance", t, func() {
		refreshDB()
		Convey("Must has working index", func() {
			payload, status, err := httpGet(apiV1Index, nil, "")

			So(err, ShouldBeNil)
			So(status, ShouldEqual, 200)
			So(payload, ShouldResemble, map[string]interface{}{
				"status":  float64(200),
				"message": "api up and running",
			})
		})
	})
}

func TestAPIV1CreateProduct(t *testing.T) {
	Convey("Given poc API instance", t, func() {
		refreshDB()
		Convey("Should have error validation payload", func() {
			params := &url.Values{}
			params.Set("name", "Product test")
			params.Set("quantity", "10")

			payload, status, err := httpPost(apiV1CreateProduct, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 400)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(400),
				"error": map[string]interface{}{
					"errors": map[string]interface{}{
						"price": []interface{}{"required"},
					},
					"message": "Validation error",
				},
			})
		})

		Convey("Should successfully create new product", func() {
			params := &url.Values{}
			params.Set("name", "Product test")
			params.Set("quantity", "10")
			params.Set("price", "20000")

			payload, status, err := httpPost(apiV1CreateProduct, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 200)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(200),
				"results": map[string]interface{}{
					"message": "Create product success",
				},
			})
		})
	})
}

func TestAPIV1CreateOrder(t *testing.T) {
	Convey("Given poc API instance", t, func() {
		refreshDB()
		Convey("Should have error validation payload", func() {
			params := &url.Values{}
			params.Set("quantity", "10")

			payload, status, err := httpPost(apiV1CreateOrder, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 400)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(400),
				"error": map[string]interface{}{
					"errors": map[string]interface{}{
						"product_id": []interface{}{"required"},
					},
					"message": "Validation error",
				},
			})
		})

		Convey("Should have error validation, product not found", func() {
			params := &url.Values{}
			params.Set("quantity", "10")
			params.Set("product_id", "1")

			payload, status, err := httpPost(apiV1CreateOrder, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 400)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(400),
				"error": map[string]interface{}{
					"errors": map[string]interface{}{
						"product_id": []interface{}{"not found"},
					},
					"message": "Validation error",
				},
			})
		})

		Convey("Should have error validation, invalid quantity", func() {
			// Create new product
			err := db.Exec("INSERT INTO products (name, quantity, price) VALUES ('Product Test', 10, 120000)").Error
			So(err, ShouldBeNil)

			params := &url.Values{}
			params.Set("quantity", "20")
			params.Set("product_id", "1")

			payload, status, err := httpPost(apiV1CreateOrder, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 400)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(400),
				"error": map[string]interface{}{
					"errors": map[string]interface{}{
						"quantity": []interface{}{"invalid"},
					},
					"message": "Validation error",
				},
			})
		})

		Convey("Should successfully create order", func() {
			// Create new product
			err := db.Exec("INSERT INTO products (name, quantity, price) VALUES ('Product Test', 10, 120000)").Error
			So(err, ShouldBeNil)

			params := &url.Values{}
			params.Set("quantity", "5")
			params.Set("product_id", "1")

			payload, status, err := httpPost(apiV1CreateOrder, params, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 200)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(200),
				"results": map[string]interface{}{
					"message": "Create order success",
				},
			})
		})
	})
}

func TestAPIV1GetAllProduct(t *testing.T) {
	Convey("Given poc API instance", t, func() {
		refreshDB()
		Convey("Should successfully return products data", func() {
			// Create sample product
			params := &url.Values{}
			params.Set("name", "Product test")
			params.Set("quantity", "10")
			params.Set("price", "20000")

			_, _, err := httpPost(apiV1CreateProduct, params, "")
			So(err, ShouldBeNil)

			var dbProductCreatedAt time.Time
			err = db.Table("products").Select("created_at").Row().Scan(&dbProductCreatedAt)
			So(err, ShouldBeNil)

			payload, status, err := httpGet(apiV1GetAllProduct, nil, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 200)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(200),
				"results": map[string]interface{}{
					"products": []interface{}{
						map[string]interface{}{
							"id":         float64(1),
							"name":       "Product test",
							"quantity":   float64(10),
							"price":      float64(20000),
							"created_at": dbProductCreatedAt.Format(time.RFC3339),
						},
					},
				},
			})
		})
	})
}

func TestAPIV1GetAllOrder(t *testing.T) {
	Convey("Given poc API instance", t, func() {
		refreshDB()
		Convey("Should successfully return orders data", func() {
			// Create sample product
			params := &url.Values{}
			params.Set("name", "Product test")
			params.Set("quantity", "10")
			params.Set("price", "20000")

			_, _, err := httpPost(apiV1CreateProduct, params, "")
			So(err, ShouldBeNil)

			var dbProductCreatedAt time.Time
			err = db.Table("products").Select("created_at").Where("id=?", 1).Row().Scan(&dbProductCreatedAt)
			So(err, ShouldBeNil)

			// Create sample order
			err = db.Exec("INSERT INTO orders (product_id, quantity, created_at) VALUES (1, 5, now())").Error
			So(err, ShouldBeNil)
			var dbOrderCreatedAt time.Time
			err = db.Table("orders").Select("created_at").Where("id=?", 1).Row().Scan(&dbOrderCreatedAt)
			So(err, ShouldBeNil)

			payload, status, err := httpGet(apiV1GetAllOrder, nil, "")
			So(err, ShouldBeNil)
			So(status, ShouldEqual, 200)
			So(payload, ShouldResemble, map[string]interface{}{
				"status": float64(200),
				"results": map[string]interface{}{
					"orders": []interface{}{
						map[string]interface{}{
							"id":         float64(1),
							"quantity":   float64(5),
							"created_at": dbOrderCreatedAt.Format(time.RFC3339),
							"product": map[string]interface{}{
								"id":         float64(1),
								"name":       "Product test",
								"quantity":   float64(10),
								"price":      float64(20000),
								"created_at": dbProductCreatedAt.Format(time.RFC3339),
							},
						},
					},
				},
			})
		})
	})
}
