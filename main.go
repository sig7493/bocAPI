package main

import (
	"os"
	"log"
	//"fmt"
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

	apihost := getEnvVariable("API_HOST")
	apiport := getEnvVariable("API_PORT")

	//fmt.Printf("eshost = %s and esport = %s \n", eshost, esport)

	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	/* e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3002", "http://10.175.166.9:3002"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	  })) */

	// Set up basic auth with username=foo and password=bar
	/* e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if username == "gops" && password == "password" {
				return true, nil
			}
			fmt.Printf("error!!!")
			return false, nil
		},
	})) */

	/* //Define Routes
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get_token/{userid}/{password}", api.GetTokenHandler).Methods("GET")
//	router.HandleFunc("/get_by_process_run_date_id/{token}/{PROCESS_RUN_DATE_ID}", api.GetByProcessRunDateIdHandler).Methods("GET")
//	router.HandleFunc("/get_between_process_run_date_id/{token}/{START_PROCESS_RUN_DATE_ID}/{END_PROCESS_RUN_DATE_ID}", api.GetBetweenProcessRunDateIdsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(apihost + ":" + apiport, router)) */

	// Route => handler
	e.GET("/generate_token/:hostip/:userid/:passwd", api.GetTokenHandler)

	/* e.GET("/count_by_process_run_date_id/:token/:processrundateid", api.GetByProcessRunDateIdHandler)
	e.GET("/counts_between_process_run_date_ids/:token/:fromprocessrundateid/:toprocessrundateid", api.GetBetweenProcessRunDateIdsHandler) */

	apigroup := e.Group("/api")
	api.UseSubroute(apigroup)
	//////apigroup.Use(middleware.JWT([]byte(getEnvVariable("POC_JWT_SECRET"))))
	//apigroup.GET("", api.Restricted)

	// Start server
	e.Logger.Fatal(e.Start(apihost + ":" + apiport))

  }