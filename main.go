package main

import (
	"os"
	"log"
	"fmt"
	//"net/http"
	"github.com/joho/godotenv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sig7493/bocAPI/api"
	// "github.com/gorilla/mux"

	/* "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag" */

	_ "github.com/sig7493/bocAPI/docs" // This line is necessary for go-swagger to find your docs!
)

func getEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

  func main() {

	eshost := getEnvVariable("ES_HOST")
	esport := getEnvVariable("ES_PORT")
	apihost := getEnvVariable("API_HOST")
	apiport := getEnvVariable("API_PORT")

	fmt.Printf("eshost = %s and esport = %s \n", eshost, esport)

	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set up basic auth with username=foo and password=bar
	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if username == "gops" && password == "password" {
				return true, nil
			}
			return false, nil
		},
	}))

	/* //Define Routes
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get_token/{userid}/{password}", api.GetTokenHandler).Methods("GET")
//	router.HandleFunc("/get_by_process_run_date_id/{token}/{PROCESS_RUN_DATE_ID}", api.GetByProcessRunDateIdHandler).Methods("GET")
//	router.HandleFunc("/get_between_process_run_date_id/{token}/{START_PROCESS_RUN_DATE_ID}/{END_PROCESS_RUN_DATE_ID}", api.GetBetweenProcessRunDateIdsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(apihost + ":" + apiport, router)) */

	// Route => handler
	e.GET("/get_token", api.GetTokenHandler)
	e.GET("/get_by_process_run_date_id", api.GetByProcessRunDateIdHandler)

	// Start server
	e.Logger.Fatal(e.Start(apihost + ":" + apiport))

  }