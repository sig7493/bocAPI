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

// GetNotesDestructionDetailsRequest represents body of get_notes_destruction_details request.
type GetNotesDestructionDetailsRequest struct {
	Printbatchid string `json:"printbatchid" query:"printbatchid"`
	Year int `json:"year" query:"year"`
	Qtr string `json:"qtr" query:"qtr"`
	Month string `json:"month" query:"month"`
	Denom string `json:"denom" query:"denom"`
}

type GetNotesDestructionDetails struct {
	Scrollid string `json:"scrollid"`
	//DestructionDetails string `json:"destructiondetails"`
	Print_Batch_ID string `json:"print_batch_id"`
	Bn_denom_en_nm string `json:"bn_denom_en_nm"`
	T_exceed_string_txt string `json:"t_exceed_string_txt"`
	Year_nb int `json:"year_nb"`
	Quarter_en_nm string `json:"quarter_en_nm"`
	Month_en_nm string `json:"month_en_nm"`
	Row_counter_nb int `json:"row_counter_nb"`
}

type GetNotesDestructionDetailsResponse []GetNotesDestructionDetails

type GetNotesDestructionAgg struct {
	Print_Batch_ID string `json:"print_batch_id"`
	Bn_denom_en_nm string `json:"bn_denom_en_nm"`
	T_exceed_string_txt string `json:"t_exceed_string_txt"`
	Year_nb int `json:"year_nb"`
	Quarter_en_nm string `json:"quarter_en_nm"`
	Month_en_nm string `json:"month_en_nm"`
	Sum_row_counter_nb int `json:"sum_row_counter_nb"`
	Rank int `json:"rank"`
}

type GetNotesDestructionAggResponse []GetNotesDestructionAgg

type GetVelocityDataRequest struct {
	Dn string `json:"dn" query:"dn"`
	Series string `json:"srs" query:"srs"`
	Day int64 `json:"day" query:"day"`
	Wk int64 `json:"wk" query:"wk"`
	Mn int64 `json:"mn" query:"mn"`
	Yr int64 `json:"yr" query:"yr"`
	Startdate string `json:"sd" query:"sd"`
	Enddate string `json:"ed" query:"ed"`
	Fittype string `json:"ft" query:"ft"`
	Sn string `json:"sn" query:"sn"`
	Stacker string `json:"sk" query:"sk"`
}

type GetVelocityData struct {
	E_Emission_Ce string `json:"e_emission_ce"`
	Dnm_Bank_Note_Denom_Am string `json:"dnm_bank_note_denom_am"`
	Sn_Serial_Number string `json:"sn_serial_number"`
	Area_Note_Area_Rt string `json:"area_note_area_rt"`
	Boc_Loc_Id string `json:"boc_loc_id"`
	Creases_Crumple_Score string `json:"creases_crumple_score"`
	Crn_Miss_Corner_Areas_Sum string `json:"crn_miss_corner_areas_sum"`
	Etr_Closed_Edge_Tears_Sum string `json:"etr_closed_edge_tears_sum"`
	Fed_Folded_Edge_Areas_Sum string `json:"fed_folded_edge_areas_sum"`
	Fhlg_Foil_Hlgraphic_Effect_Ic string `json:"fhlg_foil_hlgraphic_effect_ic"`
	Fld_Fold_Corner_Areas_Sum string `json:"fld_fold_corner_areas_sum"`
	Foil_Fitness string `json:"foil_fitness"`
	Fol_Foil_Area_Miss_Sum string `json:"fol_foil_area_miss_sum"`
	Foreign_Marks_Total_Back string `json:"foreign_marks_total_back"`
	Foreign_Marks_Total_Front string `json:"foreign_marks_total_front"`
	Gcwf_Graffiti_Over_Win_Foil string `json:"gcwf_graffiti_over_win_foil"`
	Gtb_Graffiti_On_Back_Sum string `json:"gtb_graffiti_on_back_sum"`
	Gtf_Graffiti_On_Front_Sum string `json:"gtf_graffiti_on_front_sum"`
	Gwb_Max_Graffiti_On_Back string `json:"gwb_max_graffiti_on_back"`
	Gwf_Max_Graffiti_On_Front string `json:"gwf_max_graffiti_on_front"`
	Hilo_Hi_Low_Note_Ride_Nb string `json:"hilo_hi_low_note_ride_nb"`
	Hl_Hole_Areas_Sum string `json:"hl_hole_areas_sum"`
	Iwb_Ink_Wear_On_Back_Nb string `json:"iwb_ink_wear_on_back_nb"`
	Iwf_Ink_Wear_On_Front_Nb string `json:"iwf_ink_wear_on_front_nb"`
	Len_Note_Length_Nb string `json:"len_note_length_nb"`
	Machine_Ce string `json:"machine_ce"`
	Mcrn_Max_Miss_Corner_Area string `json:"mcrn_max_miss_corner_area"`
	Med_Miss_Edge_Areas_Sum string `json:"med_miss_edge_areas_sum"`
	Metr_Max_Closed_Edge_Tear string `json:"metr_max_closed_edge_tear"`
	Mfed_Max_Folded_Edge_Area string `json:"mfed_max_folded_edge_area"`
	Mfld_Max_Folded_Corner_Area string `json:"mfld_max_folded_corner_area"`
	Mfol_Max_Foil_Scratch_Length string `json:"mfol_max_foil_scratch_length"`
	Mhl_Max_Hole_Area string `json:"mhl_max_hole_area"`
	Miwb_Max_Ink_Wear_On_Back string `json:"miwb_max_ink_wear_on_back"`
	Miwf_Max_Ink_Wear_On_Front string `json:"miwf_max_ink_wear_on_front"`
	Mmed_Max_Miss_Edge_Area string `json:"mmed_max_miss_edge_area"`
	Mmtr_Max_Closed_Tear string `json:"mmtr_max_closed_tear"`
	Motr_Max_Open_Tear_Length string `json:"motr_max_open_tear_length"`
	Mtr_Closed_Tears_Sum string `json:"mtr_closed_tears_sum"`
	Optically_Variable_Ink_Presence string `json:"optically_variable_ink_presence"`
	Optically_Variable_Ink_Score string `json:"optically_variable_ink_score"`
	Otr_Open_Tear_Lengths_Sum string `json:"otr_open_tear_lengths_sum"`
	Owl_Opacification_Wear_Level string `json:"owl_opacification_wear_level"`
	O_Orientation_Ce string `json:"o_orientation_ce"`
	Rcrn_Max_Miss_Corner_Region string `json:"rcrn_max_miss_corner_region"`
	Retr_Max_Clsd_Edge_Tear_Region string `json:"retr_max_clsd_edge_tear_region"`
	Rfed_Max_Folded_Edge_Region string `json:"rfed_max_folded_edge_region"`
	Rfld_Max_Fold_Corner_Region string `json:"rfld_max_fold_corner_region"`
	Rgtb_Max_Graffiti_Back_Region string `json:"rgtb_max_graffiti_back_region"`
	Rgtf_Max_Graffiti_Front_Region string `json:"rgtf_max_graffiti_front_region"`
	Rhl_Max_Hole_Region string `json:"rhl_max_hole_region"`
	Riwb_Max_Ink_Wear_Back_Region string `json:"riwb_max_ink_wear_back_region"`
	Riwf_Max_Ink_Wear_Front_Region string `json:"riwf_max_ink_wear_front_region"`
	Rmed_Max_Miss_Edge_Region string `json:"rmed_max_miss_edge_region"`
	Rmtr_Max_Closed_Tear_Region string `json:"rmtr_max_closed_tear_region"`
	Rotr_Max_Open_Tear_Region string `json:"rotr_max_open_tear_region"`
	Rtap_Tape_Region string `json:"rtap_tape_region"`
	Sfol_Foil_Scratch_Lengths_Sum string `json:"sfol_foil_scratch_lengths_sum"`
	Skew_Note_Skew_Nb string `json:"skew_note_skew_nb"`
	Slb_Soil_On_Back_Sum string `json:"slb_soil_on_back_sum"`
	Slf_Soil_On_Front_Sum string `json:"slf_soil_on_front_sum"` 
	Stnb_Staining_Discolor_Back string `json:"stnb_staining_discolor_back"`
	Stnf_Staining_Discolor_Front string `json:"stnf_staining_discolor_front"` 
	Swf_Small_Wind_Feature_Ic string `json:"swf_small_wind_feature_ic"`
	Tacf_Tactile_Feature_Ic string `json:"tacf_tactile_feature_ic"`
	Tape_Tape_Areas_Sum string `json:"tape_tape_areas_sum"`
	Wid_Note_Width_Nb string `json:"wid_note_width_nb"`
	Adi string `json:"adi"`
	Adi_Day int `json:"adi_day"`
	Adi_Mnth int `json:"adi_mnth"`
	Adi_Process_Run_Date_Id string `json:"adi_process_run_date_id"`
	Adi_Qtr int `json:"adi_qtr"`
	Adi_Wk int `json:"adi_wk"`
	Adi_Yr int `json:"adi_yr"`
	Bps_Shift_Nb string `json:"bps_shift_nb"`
	Deposit_Nb string `json:"deposit_nb"`
	Fi_Ce string `json:"fi_ce"`
	Output_Stacker_En_Nm string `json:"output_stacker_en_nm"`
	Rdp_Ce string `json:"rdp_ce"`
}

type GetVelocityDataResponse []GetVelocityData

type GetVelocityDataAllResponse struct {
	ScrollID string `json:"scrollid"`
	Total int `json:"total"`
	Payload GetVelocityDataResponse `json:"payload"`
}

type GetVelocityAggWearCatData struct{
	Series int `json:"series"`
	Year int `json:"year"`
	Qtr int `json:"qtr"`
	Denom int `json:"denom"`
	FLD_Notes_with_Zero_Count float64 `json:"fld_notes_with_zero_count"`
	FLD_Average_Excluding_Zero float64 `json:"fld_average_excluding_zero"`
	FLD_Non_Zero_Pct float64 `json:"fld_non_zero_pct"`
	FLD_gt_eq_Threshold_Pct float64 `json:"fld_gt_eq_threshold_pct"`
	FLD_Average_Reading float64 `json:"fld_average_reading"`
	FED_Notes_with_Zero_Count float64 `json:"fed_notes_with_zero_count"`
	FED_Average_Excluding_Zero float64 `json:"fed_average_excluding_zero"`
	FED_Non_Zero_Pct float64 `json:"fed_non_zero_pct"`
	FED_gt_eq_Threshold_Pct float64 `json:"fed_gt_eq_threshold_pct"`
	FED_Average_Reading float64 `json:"fed_average_reading"`
	CRS_Notes_with_Zero_Count float64 `json:"crs_notes_with_zero_count"`
	CRS_Excluding_Zero_Mean float64 `json:"crs_excluding_zero_mean"`
	CRS_Non_Zero_Pct float64 `json:"crs_non_zero_pct"`
	CRS_gt_eq_Threshold_Pct float64 `json:"crs_gt_eq_threshold_pct"`
	CRS_Average_Reading float64 `json:"crs_average_reading"`
	HL_Notes_with_Zero_Count float64 `json:"hl_notes_with_zero_count"`
	HL_Excluding_Zero_Mean float64 `json:"hl_excluding_zero_mean"`
	HL_Non_Zero_Pct float64 `json:"hl_non_zero_pct"`
	HL_gt_eq_Threshold_Pct float64 `json:"hl_gt_eq_threshold_pct"`
	HL_Average_Reading float64 `json:"hl_average_reading"`
	OTR_Notes_with_Zero_Count float64 `json:"otr_notes_with_zero_count"`
	OTR_Excluding_Zero_Mean float64 `json:"otr_excluding_zero_mean"`
	OTR_Non_Zero_Pct float64 `json:"otr_non_zero_pct"`
	OTR_gt_eq_Threshold_Pct float64 `json:"otr_gt_eq_threshold_pct"`
	OTR_Average_Reading float64 `json:"otr_average_reading"`
	MOTR_Notes_with_Zero_Count float64 `json:"motr_notes_with_zero_count"`
	MOTR_Excluding_Zero_Mean float64 `json:"motr_excluding_zero_mean"`
	MOTR_Non_Zero_Pct float64 `json:"motr_non_zero_pct"`
	MOTR_gt_eq_Threshold_Pct float64 `json:"motr_gt_eq_threshold_pct"`
	MOTR_Average_Reading float64 `json:"motr_average_reading"`
	ETR_MTR_Notes_with_Zero_Count float64 `json:"etr_mtr_notes_with_zero_count"`
	ETR_MTR_Excluding_Zero_Mean float64 `json:"etr_mtr_excluding_zero_mean"`
	ETR_MTR_Non_Zero_Pct float64 `json:"etr_mtr_non_zero_pct"`
	ETR_MTR_gt_eq_Threshold_Pct float64 `json:"etr_mtr_gt_eq_threshold_pct"`
	ETR_MTR_Average_Reading float64 `json:"etr_mtr_average_reading"`
	METR_MMTR_Notes_with_Zero_Count float64 `json:"metr_mmtr_notes_with_zero_count"`
	METR_MMTR_Excluding_Zero_Mean float64 `json:"metr_mmtr_excluding_zero_mean"`
	METR_MMTR_Non_Zero_Pct float64 `json:"metr_mmtr_non_zero_pct"`
	METR_MMTR_gt_eq_Threshold_Pct float64 `json:"metr_mmtr_gt_eq_threshold_pct"`
	METR_MMTR_Average_Reading float64 `json:"metr_mmtr_average_reading"`
	CRN_Notes_with_Zero_Count float64 `json:"crn_notes_with_zero_count"`
	CRN_Excluding_Zero_Mean float64 `json:"crn_excluding_zero_mean"`
	CRN_Non_Zero_Pct float64 `json:"crn_non_zero_pct"`
	CRN_gt_eq_Threshold_Pct float64 `json:"crn_gt_eq_threshold_pct"`
	CRN_Average_Reading float64 `json:"crn_average_reading"`
	MED_Notes_with_Zero_Count float64 `json:"med_notes_with_zero_count"`
	MED_Excluding_Zero_Mean float64 `json:"med_excluding_zero_mean"`
	MED_Non_Zero_Pct float64 `json:"med_non_zero_pct"`
	MED_gt_eq_Threshold_Pct float64 `json:"med_gt_eq_threshold_pct"`
	MED_Average_Reading float64 `json:"med_average_reading"`
	TAPE_Notes_with_Zero_Count float64 `json:"tape_notes_with_zero_count"`
	TAPE_Excluding_Zero_Mean float64 `json:"tape_excluding_zero_mean"`
	TAPE_Non_Zero_Pct float64 `json:"tape_non_zero_pct"`
	TAPE_gt_eq_Threshold_Pct float64 `json:"tape_gt_eq_threshold_pct"`
	TAPE_Average_Reading float64 `json:"tape_average_reading"`
	FOL_Notes_with_Zero_Count float64 `json:"fol_notes_with_zero_count"`
	FOL_Excluding_Zero_Mean float64 `json:"fol_excluding_zero_mean"`
	FOL_Non_Zero_Pct float64 `json:"fol_non_zero_pct"`
	FOL_gt_eq_Threshold_Pct float64 `json:"fol_gt_eq_threshold_pct"`
	FOL_Average_Reading float64 `json:"fol_average_reading"`
	SFOL_Notes_with_Zero_Count float64 `json:"sfol_notes_with_zero_count"`
	SFOL_Excluding_Zero_Mean float64 `json:"sfol_excluding_zero_mean"`
	SFOL_Non_Zero_Pct float64 `json:"sfol_non_zero_pct"`
	SFOL_gt_eq_Threshold_Pct float64 `json:"sfol_gt_eq_threshold_pct"`
	SFOL_Average_Reading float64 `json:"sfol_average_reading"`
	MFOL_Notes_with_Zero_Count float64 `json:"mfol_notes_with_zero_count"`
	MFOL_Excluding_Zero_Mean float64 `json:"mfol_excluding_zero_mean"`
	MFOL_Non_Zero_Pct float64 `json:"mfol_non_zero_pct"`
	MFOL_gt_eq_Threshold_Pct float64 `json:"mfol_gt_eq_threshold_pct"`
	MFOL_Average_Reading float64 `json:"mfol_average_reading"`
	IWB_Notes_with_Zero_Count float64 `json:"iwb_notes_with_zero_count"`
	IWB_Excluding_Zero_Mean float64 `json:"iwb_excluding_zero_mean"`
	IWB_Non_Zero_Pct float64 `json:"iwb_non_zero_pct"`
	IWB_gt_eq_Threshold_Pct float64 `json:"iwb_gt_eq_threshold_pct"`
	IWB_Average_Reading float64 `json:"iwb_average_reading"`
	IWF_Notes_with_Zero_Count float64 `json:"iwf_notes_with_zero_count"`
	IWF_Excluding_Zero_Mean float64 `json:"iwf_excluding_zero_mean"`
	IWF_Non_Zero_Pct float64 `json:"iwf_non_zero_pct"`
	IWF_gt_eq_Threshold_Pct float64 `json:"iwf_gt_eq_threshold_pct"`
	IWF_Average_Reading float64 `json:"iwf_average_reading"`
	GTF_Notes_with_Zero_Count float64 `json:"gtf_notes_with_zero_count"`
	GTF_Excluding_Zero_Mean float64 `json:"gtf_excluding_zero_mean"`
	GTF_Non_Zero_Pct float64 `json:"gtf_non_zero_pct"`
	GTF_gt_eq_Threshold_Pct float64 `json:"gtf_gt_eq_threshold_pct"`
	GTF_Average_Reading float64 `json:"gtf_average_reading"`
	GTB_Notes_with_Zero_Count float64 `json:"gtb_notes_with_zero_count"`
	GTB_Excluding_Zero_Mean float64 `json:"gtb_excluding_zero_mean"`
	GTB_Non_Zero_Pct float64 `json:"gtb_non_zero_pct"`
	GTB_gt_eq_Threshold_Pct float64 `json:"gtb_gt_eq_threshold_pct"`
	GTB_Average_Reading float64 `json:"gtb_average_reading"`
	GCWF_Notes_with_Zero_Count float64 `json:"gcwf_notes_with_zero_count"`
	GCWF_Excluding_Zero_Mean float64 `json:"gcwf_excluding_zero_mean"`
	GCWF_Non_Zero_Pct float64 `json:"gcwf_non_zero_pct"`
	GCWF_gt_eq_Threshold_Pct float64 `json:"gcwf_gt_eq_threshold_pct"`
	GCWF_Average_Reading float64 `json:"gcwf_average_reading"`
	OWL_Notes_with_Zero_Count float64 `json:"owl_notes_with_zero_count"`
	OWL_Excluding_Zero_Mean float64 `json:"owl_excluding_zero_mean"`
	OWL_Non_Zero_Pct float64 `json:"owl_non_zero_pct"`
	OWL_gt_eq_Threshold_Pct float64 `json:"owl_gt_eq_threshold_pct"`
	OWL_Average_Reading float64 `json:"owl_average_reading"`
	STNF_Notes_with_Zero_Count float64 `json:"stnf_notes_with_zero_count"`
	STNF_Excluding_Zero_Mean float64 `json:"stnf_excluding_zero_mean"`
	STNF_Non_Zero_Pct float64 `json:"stnf_non_zero_pct"`
	STNF_gt_eq_Threshold_Pct float64 `json:"stnf_gt_eq_threshold_pct"`
	STNF_Average_Reading float64 `json:"stnf_average_reading"`
	STNB_Notes_with_Zero_Count float64 `json:"stnb_notes_with_zero_count"`
	STNB_Excluding_Zero_Mean float64 `json:"stnb_excluding_zero_mean"`
	STNB_Non_Zero_Pct float64 `json:"stnb_non_zero_pct"`
	STNB_gt_eq_Threshold_Pct float64 `json:"stnb_gt_eq_threshold_pct"`
	STNB_Average_Reading float64 `json:"stnb_average_reading"`
	FMTF_Notes_with_Zero_Count float64 `json:"fmtf_notes_with_zero_count"`
	FMTF_Excluding_Zero_Mean float64 `json:"fmtf_excluding_zero_mean"`
	FMTF_Non_Zero_Pct float64 `json:"fmtf_non_zero_pct"`
	FMTF_gt_eq_Threshold_Pct float64 `json:"fmtf_gt_eq_threshold_pct"`
	FMTF_Average_Reading float64 `json:"fmtf_average_reading"`
	FMTB_Notes_with_Zero_Count float64 `json:"fmtb_notes_with_zero_count"`
	FMTB_Excluding_Zero_Mean float64 `json:"fmtb_excluding_zero_mean"`
	FMTB_Non_Zero_Pct float64 `json:"fmtb_non_zero_pct"`
	FMTB_gt_eq_Threshold_Pct float64 `json:"fmtb_gt_eq_threshold_pct"`
	FMTB_Average_Reading float64 `json:"FMTB_Average_Reading"`
	FFIT_Notes_with_Zero_Count float64 `json:"ffit_notes_with_zero_count"`
	FFIT_Excluding_Zero_Mean float64 `json:"ffit_excluding_zero_mean"`
	FFIT_Non_Zero_Pct float64 `json:"ffit_non_zero_pct"`
	FFIT_gt_eq_Threshold_Pct float64 `json:"ffit_gt_eq_threshold_pct"`
	FFIT_Average_Reading float64 `json:"ffit_average_reading"`
	OVIP_Notes_with_Zero_Count float64 `json:"ovip_notes_with_zero_count"`
	OVIP_Excluding_Zero_Mean float64 `json:"ovip_excluding_zero_mean"`
	OVIP_Non_Zero_Pct float64 `json:"ovip_non_zero_pct"`
	OVIP_gt_eq_Threshold_Pct float64 `json:"ovip_gt_eq_threshold_pct"`
	OVIP_Average_Reading float64 `json:"ovip_average_reading"`
	OVIS_Notes_with_Zero_Count float64 `json:"ovis_notes_with_zero_count"`
	OVIS_Excluding_Zero_Mean float64 `json:"ovis_excluding_zero_mean"`
	OVIS_Non_Zero_Pct float64 `json:"ovis_non_zero_pct"`
	OVIS_gt_eq_Threshold_Pct float64 `json:"ovis_gt_eq_threshold_pct"`
	OVIS_Average_Reading float64 `json:"ovis_average_reading"`
	}

type GetVelocityAggWearCatDataResponse []GetVelocityAggWearCatData

type GetVelocityAggWearCategoryDataResponse struct {
	Total int `json:"total"`
	Payload GetVelocityAggWearCatDataResponse `json:"payload"`
}