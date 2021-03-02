package api

import (
/* 	"bytes"
	"context"
	"log"
	"net/url"
	"os"
	
	"strings"
	"time" */

	"fmt"
	//"encoding/json"
	"strconv"
	"net/http"
	//"github.com/gorilla/mux"
	"github.com/labstack/echo"

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
	Token string `json:"token"`
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
	Token string `json:"token" query:"token"`
}

type GetBetweenProcessRunDateIds struct {
			ProcessRunDateID int `json:"processrundateid"`
			Count int `json:"count"`
}

// swagger:model GetBetweenProcessRunDateIdsResponse
// GetBetweenProcessRunDateIdsResponse represents body of get_between_process_run_date_ids response.
type GetBetweenProcessRunDateIdsResponse []GetBetweenProcessRunDateIds


// GetTokenHandler handles incoming get_token requests
func GetTokenHandler(ctx echo.Context) error {
	/* vars := mux.Vars(r)
	userid := vars["userid"]
	password := vars["password"] */

	req := GetTokenRequest{}/*  */

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
    /* hostip := req.HostIP
	userid := req.UserId	
	password := req.Passwd */

	
	hostip := ctx.Param("hostip")
	userid := ctx.Param("userid")
	password := ctx.Param("passwd")

	// Perform check 

	resp := hostip + " sent by " + userid +" = " + password
	//fmt.Fprintf(w, "Token generated testtoken for user - %v with password - %v", userid, password)

	return ctx.JSON(http.StatusOK, resp)

}
func GetByProcessRunDateIdHandler(ctx echo.Context) error {
	/* vars := mux.Vars(r)
	processrundateID := vars["PROCESS_RUN_DATE_ID"]
	tokenID, ok := vars["token"] */

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

	token := ctx.Param("token")

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	fmt.Printf("token = %v\n", token)
	
	//if !ok {
	if token == "" {	
		//fmt.Fprintf(w, "Invalid Token. Please generate a token")
		resp := "Invalid Token. Please generate a token"
		return ctx.JSON(http.StatusUnauthorized, resp)
	}

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

	token := ctx.Param("token")

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	fmt.Printf("token = %v\n", token)

	if token == "" {	
		resp := "Invalid Token. Please generate a token"
		return ctx.JSON(http.StatusUnauthorized, resp)
	}

	//resp := "From : " + strconv.Itoa(fromprocessrundateID) + " To : " + strconv.Itoa(toprocessrundateID)  +" requested using token ->" + token

	resp := GetBetweenProcessRunDateIds_search_request(fromprocessrundateID, toprocessrundateID)

	return ctx.JSON(http.StatusOK, resp)
}