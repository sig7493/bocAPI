package api

import (
/* 	"bytes"
	"context"
	"log"
	"net/url"
	
	
	"strings" */
	"os"
	"time"

	"fmt"
	//"encoding/json"
	"strconv"
	"net/http"
	//"github.com/gorilla/mux"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/dgrijalva/jwt-go"

	/* "google.golang.org/grpc"

	"golang.org/x/oauth2" */

	/* "github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag" */

)

var mysecret = []byte(os.Getenv("POC_JWT_SECRET"))

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
    SigningKey: mysecret,
})

// GetTokenRequest represents body of get_token request.
type GetTokenRequest struct {
	HostIP string `json:"hostip"`
	UserId string `json:"userid"`
	Passwd string `json:"passwd"`
}

// swagger:model GetTokenResponse
// GetTokenResponse represents body of get_token response.
type GetTokenResponse struct {
	Token string `json:"token"`
}

// GetByProcessRunDateIdRequest represents body of get_by_process_run_date_id request.
type GetByProcessRunDateIdRequest struct {
	ProcessRunDateID int `json:"processrundateid"`
}

// swagger:model GetByProcessRunDateIdResponse
// GetByProcessRunDateIdResponse represents body of get_by_process_run_date_id response.
type GetByProcessRunDateIdResponse struct {
	ProcessRunDateID int `json:"processrundateid"`
	Count int `json:"count"`
}

// GetBetweenProcessRunDateIdsRequest represents body of get_between_process_run_date_ids request.
type GetBetweenProcessRunDateIdsRequest struct {
	FromProcessRunDateID int `json:"fromprocessrundateid" query:"fromprocessrundateid"`
	ToProcessRunDateID int `json:"toprocessrundateid" query:"toprocessrundateid"`
}

type GetBetweenProcessRunDateIds struct {
			ProcessRunDateID int `json:"processrundateid"`
			Count int `json:"count"`
}

// swagger:model GetBetweenProcessRunDateIdsResponse
// GetBetweenProcessRunDateIdsResponse represents body of get_between_process_run_date_ids response.
type GetBetweenProcessRunDateIdsResponse []GetBetweenProcessRunDateIds

func UseSubroute(group *echo.Group) {
	group.GET("/count_by_process_run_date_id/:processrundateid", GetByProcessRunDateIdHandler, IsLoggedIn)
	group.GET("/counts_between_process_run_date_ids/:fromprocessrundateid/:toprocessrundateid", GetBetweenProcessRunDateIdsHandler, IsLoggedIn)
}

// restricted handles jwt token validation
func Restricted(ctx echo.Context) error {
	if temp := ctx.Get("user"); temp != nil {
		user := temp.(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		fmt.Printf("name = %s", name)
		//return ctx.String(http.StatusOK, "Welcome "+name+"! \n")
		return nil
		
	}
	return echo.ErrUnauthorized
}

// GetTokenHandler handles incoming get_token requests
func GetTokenHandler(ctx echo.Context) error {
	/* vars := mux.Vars(r)
	userid := vars["userid"]
	password := vars["password"] */

    /* hostip := req.HostIP
	userid := req.UserId	
	password := req.Passwd */

	
	//hostip := ctx.Param("hostip")
	userid := ctx.Param("userid")
	//password := ctx.Param("passwd")

	// Perform check in AD if userID exists

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	//claims["hostip"] = hostip
	claims["name"] = userid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mysecret)

	if err != nil {
		fmt.Errorf("Something went wrong generating the token : %s", err.Error())
		return err
	}

	resp := GetTokenResponse{Token: tokenString}

	//resp := tokenString
	//resp := "host= " + hostip + " for user = " + userid
	//fmt.Fprintf(w, "Token generated testtoken for user - %v with password - %v", userid, password)

	return ctx.JSON(http.StatusOK, resp)

}
func GetByProcessRunDateIdHandler(ctx echo.Context) error {
	/* vars := mux.Vars(r)
	processrundateID := vars["PROCESS_RUN_DATE_ID"]
	tokenID, ok := vars["token"] */

	restricted := Restricted(ctx)

	fmt.Printf("restricted = %v\n", restricted)

	req := GetByProcessRunDateIdRequest{}
	
	/* params := ctx.QueryParam("params")
	fmt.Printf("%v\n", params)

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	json.Unmarshal([]byte(params), &req)
	fmt.Printf("req = %v\n", req)

	processrundateID := req.ProcessRunDateID
	token := req.Token */

	processrundateID, err := strconv.Atoi(ctx.Param("processrundateid"))
	if err != nil {
		fmt.Printf("%v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	//token := ctx.Param("token")

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	//fmt.Printf("token = %v\n", token)
	
	/* //if !ok {
	if token == "" {	
		//fmt.Fprintf(w, "Invalid Token. Please generate a token")
		resp := "Invalid Token. Please generate a token"
		return ctx.JSON(http.StatusUnauthorized, resp)
	} */

	//fmt.Fprintf(w, "Get Process Run Date ID: %v with token %v", processrundateID, tokenID)

	//resp := strconv.Itoa(processrundateID) + " requested using token ->" + token

	get_es_cluster_info()
	resp := GetByProcessRunDateId_search_request(processrundateID)
	

	return ctx.JSON(http.StatusOK, resp)
}
/*
func GetBetweenProcessRunDateIdsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	startprocessrundateID := vars["START_PROCESS_RUN_DATE_ID"]
	endprocessrundateID := vars["END_PROCESS_RUN_DATE_ID"]
	tokenID, ok := vars["token"]

	if !ok {
		fmt.Fprintf(w, "Invalid Token. Please generate a token")
		return
	}

	fmt.Fprintf(w, "Get between %v and %v Process Run Date IDs with token %v", startprocessrundateID, endprocessrundateID, tokenID)
} */

func GetBetweenProcessRunDateIdsHandler(ctx echo.Context) error {
	
	restricted := Restricted(ctx)

	fmt.Printf("restricted = %v\n", restricted)

	req := GetBetweenProcessRunDateIdsRequest{}
	
	/* params := ctx.QueryParam("params")
	fmt.Printf("%v\n", params)
	
	if err := ctx.Bind(&req); err != nil {
		//return echo.ErrBadRequest
		fmt.Printf("%v\n", err)
	}
	
	json.Unmarshal([]byte(params), &req)
	fmt.Printf("req = %v\n", req)

	fromprocessrundateID := req.FromProcessRunDateID
	toprocessrundateID := req.ToProcessRunDateID
	token := req.Token */

	fromprocessrundateID, err := strconv.Atoi(ctx.Param("fromprocessrundateid"))
	if err != nil {
		fmt.Printf("%v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	toprocessrundateID, err := strconv.Atoi(ctx.Param("toprocessrundateid"))
	if err != nil {
		fmt.Printf("%v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	//token := ctx.Param("token")

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	/* fmt.Printf("token = %v\n", token)

	if token == "" {	
		resp := "Invalid Token. Please generate a token"
		return ctx.JSON(http.StatusUnauthorized, resp)
	} */

	//resp := "From : " + strconv.Itoa(fromprocessrundateID) + " To : " + strconv.Itoa(toprocessrundateID)  +" requested using token ->" + token

	resp := GetBetweenProcessRunDateIds_search_request(fromprocessrundateID, toprocessrundateID)

	return ctx.JSON(http.StatusOK, resp)
}