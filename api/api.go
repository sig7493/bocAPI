package api

import (
	 "bytes"
	 //"fmt"
/*	"context"
	
	"net/url"
	
	
	"strings" */
	//"os"
	"time"

	"fmt"
	"encoding/json"
	"strconv"
	"net/http"
	//"github.com/gorilla/mux"
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"

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

/* 	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"golang.org/x/text/message" */

	"github.com/allegro/bigcache/v3"

)

var cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

func UseSubroute(group *echo.Group) {
	// group.GET("/count_by_process_run_date_id/:processrundateid", GetByProcessRunDateIdHandler, IsLoggedIn)
	// group.GET("/counts_between_process_run_date_ids/:fromprocessrundateid/:toprocessrundateid", GetBetweenProcessRunDateIdsHandler, IsLoggedIn)
	// group.GET("/counts_of_all_process_run_date_ids", GetAllProcessRunDateIdsHandler, IsLoggedIn)
	group.GET("/count_by_process_run_date_id/:processrundateid", GetByProcessRunDateIdHandler)
	group.GET("/counts_between_process_run_date_ids/:fromprocessrundateid/:toprocessrundateid", GetBetweenProcessRunDateIdsHandler)
	group.GET("/counts_of_all_process_run_date_ids", GetAllProcessRunDateIdsHandler)
	group.GET("/notes_validity_details", GetNotesValidityDetailsHandler)
	group.GET("/notes_invalidity_details", GetNotesInValidityDetailsHandler)
	group.GET("/notes_destruction_aggregate", GetNotesDestructionAggHandler)
	group.GET("/notes_destruction_details/:printbatchid/:year/:qtr/:month/:denom/:from/:scrollid/", GetNotesDestructionDetailsHandler)
	group.GET("/notes_destruction_details/*", GetNotesDestructionDetailsHandler)
	group.GET("/velocity", GetVelocityDataHandler)
	group.GET("/velocity_agg_wear_category", GetVelocityAggWearCategoryDataHandler)
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

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	//claims["hostip"] = hostip
	claims["name"] = userid
	// claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	claims["exp"] = time.Now().Add(time.Minute * 90000).Unix()

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

func GetAllProcessRunDateIdsHandler(ctx echo.Context) error {

	restricted := Restricted(ctx)

	fmt.Printf("restricted = %v\n", restricted)

	//Uncomment pg_coonect() to connect to postgresql
	//pg_connect()
	/* tmpresp, err := env.GetDetailsOfProcessRunDateId_request(20170907)

	if err != nil {
		return err
	}

	fmt.Printf("+%v\n", tmpresp) */

	resp := GetBetweenProcessRunDateIdsResponse{}
	var buf bytes.Buffer

	cache_resp, cache_err := cache.Get("all-process-run-date-ids-key")

	if cache_err == nil {
		fmt.Println("Fetching Cached Response...")
		fmt.Println(string(cache_resp))

		if err := json.NewDecoder(bytes.NewBuffer(cache_resp)).Decode(&resp); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		  }
	} else {
		fmt.Println("Fetching from ElasticSearch...")
		resp = *GetAllProcessRunDateIds_search_request()

		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		  }

		cache.Set("all-process-run-date-ids-key", []byte(buf.Bytes()))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func GetNotesValidityDetailsHandler(ctx echo.Context) error {

	resp := GetNotesValidityDetails_search_request() 

	return ctx.JSON(http.StatusOK, resp)
}

func GetNotesInValidityDetailsHandler(ctx echo.Context) error {

	resp := GetNotesInValidityDetails_search_request() 

	return ctx.JSON(http.StatusOK, resp)
}
func GetNotesDestructionAggHandler(ctx echo.Context) error {

	resp := GetNotesDestructionAggResponse{}
	var buf bytes.Buffer

	cache_resp, cache_err := cache.Get("notes-destruction-agg-key")

	if cache_err == nil {
		fmt.Println("Fetching Cached Response...")
		fmt.Println(string(cache_resp))

		if err := json.NewDecoder(bytes.NewBuffer(cache_resp)).Decode(&resp); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		  }
	} else {
		fmt.Println("Fetching from ElasticSearch...")
		resp = *GetNotesDestructionAgg_search_request()

		if err := json.NewEncoder(&buf).Encode(resp); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		  }

		cache.Set("notes-destruction-agg-key", []byte(buf.Bytes()))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func GetNotesDestructionDetailsHandler(ctx echo.Context) error {

	req := GetNotesDestructionDetailsRequest{}

	printbatchid := ctx.Param("printbatchid")
	if (len(printbatchid) == 0) {
		err := "Missing printbatchid"
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		printbatchid = ""
	}

	year, err := strconv.Atoi(ctx.Param("year"))
	if err != nil {
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		year = time.Now().Year()
	}

	qtr := ctx.Param("qtr")
	if len(qtr) == 0 {
		err := "Missing qtr"
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		qtr = "*"
	}

	month := ctx.Param("month")
	if len(month) == 0 {
		err := "Missing month"
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		month = "*"
	}

	denom := ctx.Param("denom")
	if len(denom) == 0 {
		err := "Missing denom"
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		denom = "*"
	}

	from_val, err := strconv.Atoi(ctx.Param("from"))
	if err != nil {
		fmt.Printf("%v", err)
		//return ctx.JSON(http.StatusBadRequest, err)
		from_val = 0
	}

	scroll_id := ctx.Param("scrollid")
	if len(scroll_id) == 0{
		scroll_id = "None"
	}

	/* from, err := strconv.Atoi(ctx.Param("from"))
	if err != nil {
		fmt.Printf("%v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	size, err := strconv.Atoi(ctx.Param("size"))
	if err != nil {
		fmt.Printf("%v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	} */

	if err := ctx.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	resp := GetNotesDestructionDetails_search_request(printbatchid, year, qtr, month, denom, from_val, scroll_id)
	//resp := GetNotesDestructionDetails_search_request(printbatchid, year, qtr, month, denom)

	//resp := "All Good...maybe"

	return ctx.JSON(http.StatusOK, resp)

}

func GetVelocityDataHandler(ctx echo.Context) error {

	//opts := VelocityOpts{}

	//velocitydata := GetVelocityDataRequest{}
	
	dn := ctx.QueryParam("dn")
	if dn == "" {
		dn = "*"
	}

	sr := ctx.QueryParam("srs")
	if sr == "" {
		sr = "*"
	}
	
	day := ctx.QueryParam("day")
	if day == "" {
		day = "*"
	}
	
	wk := ctx.QueryParam("wk")
	if wk == "" {
		wk = "*"
	}

	mn := ctx.QueryParam("mn")
	if mn == "" {
		mn = "*"
	}

	yr := ctx.QueryParam("yr")
	if yr == "" {
		yr = "*"
	}

	startdate := ctx.QueryParam("sd")
	if startdate == "" {
		startdate = "1970-01-01"
	}
	
	enddate := ctx.QueryParam("ed")
	if enddate == "" {
		enddate = "2122-12-31"
	}

	// fittype := ctx.QueryParam("ft")
	// if fittype == "*" {
	// 	fittype = "*"
	// }

	sn := ctx.QueryParam("sn")
	if sn == "" {
		sn = "*"
	}

	stacker := ctx.QueryParam("sk")
	if stacker == "" {
		stacker = "*"
	}

	from := ctx.QueryParam("from")
	if from == "" {
		from = "0"
	}

	scroll_id := ctx.QueryParam("sid")
	if scroll_id == "" {
		scroll_id = "None"
	}
	fmt.Println(dn, sr, day, wk, mn, yr, startdate, enddate, sn, stacker)

	resp := GetVelocityData_Search_Request(sr, dn, stacker, sn, day, wk, mn, yr, startdate, enddate, from, scroll_id)

	// return nil
	return ctx.JSON(http.StatusOK, resp)
	
}


func GetVelocityAggWearCategoryDataHandler(ctx echo.Context) error {

	resp := GetVelocityAggWearCategoryData_Search_Request()

	return ctx.JSON(http.StatusOK, resp)
}

/* func TranslateHandler(ctx echo.Context) error {

	// define a direct translation from English to French
	message.SetString(language.French, "In %v people speak %v.", "In %v spreekt men %v.")



} */

