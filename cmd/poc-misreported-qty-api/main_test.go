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
	"poc-misreported-qty/util/validator"
	"testing"

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
		)
	}

	apiV1Index         string
	apiV1CreateProduct string
)

func TestMain(m *testing.M) {
	flag.Parse()

	*debug = false

	var (
		gormInstance, _ = gorm.Open("postgres", "postgres://poc:poc-123@localhost:5437/poc-test?sslmode=disable")

		dataRepository   = sql.NewSQLDataRepository(gormInstance)
		validatorService = validator.NewDefaultValidatorService()
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
