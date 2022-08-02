package docs

import (
	//"net/http"
    "github.com/sig7493/bocAPI/api"
)

// swagger:route GET /get_token get_token idofget_tokenEndpoint
// Get token authorizes user credentials and generates token to access the APIs.
// responses:
// 200: body:GetTokenResponse

// Description of response body Get token will generates token to access the APIs.
// swagger:response body:GetTokenResponse
type get_tokenResponseWrapper struct {
	// in:body
	Body api.GetTokenResponse
}

// swagger:parameters idofget_tokenEndpoint
type get_tokenParamsWrapper struct {
	// {"hostip": "127.0.0.1", "userid" : "gops", "passwd": "password"}.
	// in.body
	Body api.GetTokenRequest
}


// swagger:route GET /get_by_process_run_date_id get_by_process_run_date_id idofget_by_process_run_date_idEndpoint
// Get count of notes processed by a PROCESS_RUN_DATE_ID.
// responses:
// 200: body:GetByProcessRunDateIdResponse

// Get count of notes processed by a PROCESS_RUN_DATE_ID.
// swagger:response body:GetByProcessRunDateIdResponse
type get_by_process_run_date_idResponseWrapper struct {
	// in:body
	Body api.GetByProcessRunDateIdResponse
}

// swagger:parameters idofget_by_process_run_date_idEndpoint
type get_by_process_run_date_idParamsWrapper struct {
	// {"processrundateid": "20210721", "token" : "jkasdk57893#5^"}.
	// in.body
	Body api.GetByProcessRunDateIdRequest
}


// swagger:route GET /get_between_process_run_date_ids get_between_process_run_date_ids idofget_between_process_run_date_idsEndpoint
// Get count of notes processed between two PROCESS_RUN_DATE_IDs.
// responses:
// 200: body:GetBetweenProcessRunDateIdsResponse

// Get count of notes processed between two PROCESS_RUN_DATE_IDs.
// swagger:response body:GetBetweenProcessRunDateIdsResponse
type get_between_process_run_date_idsResponseWrapper struct {
	// in:body
	Body api.GetBetweenProcessRunDateIdsResponse
}

// swagger:parameters idofget_between_process_run_date_idsEndpoint
type get_between_process_run_date_idsParamsWrapper struct {
	// {"fromprocessrundateid": "20210721", "toprocessrundateid": "20210723", "token" : "jkasdk57893#5^"}.
	// in.body
	Body api.GetBetweenProcessRunDateIdsRequest
}