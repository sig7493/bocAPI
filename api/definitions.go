package api

import (
	"os"	
	"github.com/labstack/echo/middleware"
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