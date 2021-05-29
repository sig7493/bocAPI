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

type GetDetailsOfProcessRunDateId struct {
	Process_Run_Date_ID int `json:"process_run_date_id"`
	Masked_Serial_Num string `json:"masked_serial_num"`
	Bps_Shift_ID int `json:"bps_shift_id"`
	Machine_ID int `json:"machine_id"`
	Print_Batch_ID int `json:"print_batch_id"`
	Rdp_ID int `json:"rdp_id"`
	Bn_Status_ID int `json:"bn_status_id"`
	Output_Stacker_ID int `json:"output_stacker_id"`
	Circ_Trial_ID int `json:"circ_trial_id"`
	Bps_Shift_Nb int `json:"bps_shift_nb"`
	Deposit_Nb int `json:"deposit_nb"`
	Row_Counter_NB int `json:"row_counter_nb"`
	Load_ID int `json:"load_id"`
}

// swagger:model GetDetailsOfProcessRunDateIdResponse
// GetDetailsOfProcessRunDateIdResponse represents body of get_details_of_process_run_date_id response.
type GetDetailsOfProcessRunDateIdResponse []GetDetailsOfProcessRunDateId

type GetNotesValidityDetails struct {
	Denomination string `json:"denomination"`
	Image_path string `json:"image_path"`
	Rgb_color string `json:"rgb_color"`
	Rgb_val string `json:"rgb_val"`
	Serial_number string `json:"serial_number"`
	Process_Run_Date_ID int `json:"process_run_date_id"`
	Bps_Shift_ID int `json:"bps_shift_id"`
	Machine_ID int `json:"machine_id"`
	Print_Batch_ID int `json:"print_batch_id"`
	Rdp_ID int `json:"rdp_id"`
	Bn_Status_ID int `json:"bn_status_id"`
	Output_Stacker_ID int `json:"output_stacker_id"`
	Circ_Trial_ID int `json:"circ_trial_id"`
	Bps_Shift_Nb int `json:"bps_shift_nb"`
	Deposit_Nb int `json:"deposit_nb"`
	Row_Counter_NB int `json:"row_counter_nb"`
	Load_ID int `json:"load_id"`
}

type GetNotesValidityDetailsResponse []GetNotesValidityDetails