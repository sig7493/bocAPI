package api

import (
	//"os"
	//"log"
	"fmt"
	"encoding/json"
	//"sync"
	//"bytes"
	//"context"
	"database/sql"
	
	"strconv"

	//"github.com/joho/godotenv"

	_ "github.com/lib/pq"

)

// Create a custom Env struct which holds a connection pool.
type Env struct {
    db *sql.DB
}

func pg_connect() {
	pg_host := getEnvVariable("POSTGRES_HOST")
	pg_port,_ := strconv.Atoi(getEnvVariable("POSTGRES_PORT"))
	pg_user := getEnvVariable("POSTGRES_USER")
	pg_pass := getEnvVariable("POSTGRES_PASSWORD")
	pg_db := getEnvVariable("POSTGRES_DB")
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pg_host, pg_port, pg_user, pg_pass, pg_db)

	fmt.Printf(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
	panic(err)
	//return nil, err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
	panic(err)
	//return nil, err
	}

	// Create an instance of Env containing the connection pool.
	env := &Env{db: db}
	
	fmt.Println("Successfully connected!")

	resp, err := env.GetDetailsOfProcessRunDateId_request(20170907)

	if err != nil {
		fmt.Printf("%v", err)
	}

	out, err := json.Marshal(resp)

	if err != nil {
		fmt.Printf("Error converting to JSON!!!\n %v", err)
	}

	fmt.Printf("+%v\n", string(out))

	//return db, nil
	//return env
}

func (env *Env) GetDetailsOfProcessRunDateId_request(processrundateid int) (*GetDetailsOfProcessRunDateIdResponse, error) {

	resp := GetDetailsOfProcessRunDateIdResponse{}

	/* err := env.db.Ping()
	if err != nil {
	fmt.Printf("Postgresql is not alive!!!\n")
	return nil, err
	} */

	sqlStatement := `SELECT "PROCESS_RUN_DATE_ID",masked_ser_num, "BPS_SHIFT_ID", 
	"MACHINE_ID", "PRINT_BATCH_ID", "RDP_ID", "BN_STATUS_ID", 
	"OUTPUT_STACKER_ID", "CIRC_TRIAL_ID", "BPS_SHIFT_NB", "DEPOSIT_NB", 
	"ROW_COUNTER_NB", "LOAD_ID"
	FROM public.STG_notes_processing where STG_notes_processing."PROCESS_RUN_DATE_ID" = $1`

	rows, err:= env.db.Query(sqlStatement, 20170907)

	if err != nil {
		//panic(err)
		return nil, err
	  }

	defer rows.Close()

	for rows.Next() {

		var processrundateid_detail GetDetailsOfProcessRunDateId

		err := rows.Scan(&processrundateid_detail.Process_Run_Date_ID, &processrundateid_detail.Masked_Serial_Num, &processrundateid_detail.Bps_Shift_ID, 
			&processrundateid_detail.Machine_ID, &processrundateid_detail.Print_Batch_ID, &processrundateid_detail.Rdp_ID, &processrundateid_detail.Bn_Status_ID, 
			&processrundateid_detail.Output_Stacker_ID, &processrundateid_detail.Circ_Trial_ID, &processrundateid_detail.Bps_Shift_Nb, &processrundateid_detail.Deposit_Nb, 
			&processrundateid_detail.Row_Counter_NB, &processrundateid_detail.Load_ID)

		if err != nil {
			return nil, err
		}
		
		resp = append(resp, processrundateid_detail)
	}

	if err = rows.Err(); err != nil {
        return nil, err
	}
	//rows.Close()
	//db.Close()
	return &resp, nil

}