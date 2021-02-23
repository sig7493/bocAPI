package api

import (
/* 	"bytes"
	"context"
	"log"
	"net/url"
	"os"
	"encoding/json"
	"strings"
	"time" */

	//i"fmt"
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

// GetTokenHandler handles incoming get_token requests
func GetTokenHandler(ctx echo.Context) error {
	/* vars := mux.Vars(r)
	userid := vars["userid"]
	password := vars["password"] */

	req := GetTokenRequest{}

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}
    hostip := req.HostIP
	userid := req.UserId	
	password := req.Passwd

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

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	processrundateID := req.ProcessRunDateID
	token := req.Token


	//if !ok {
	if token == "" {	
		//fmt.Fprintf(w, "Invalid Token. Please generate a token")
		resp := "Invalid Token. Please generate a token"
		return ctx.JSON(http.StatusUnauthorized, resp)
	}

	//fmt.Fprintf(w, "Get Process Run Date ID: %v with token %v", processrundateID, tokenID)

	resp := strconv.Itoa(processrundateID) + " requested using token ->" + token

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
