package api

import (
	"os"
	"io"
	"log"
	"fmt"
	"encoding/json"
	"sync"
	"bytes"
	"context"
	//"reflect"
	"strings"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	//"github.com/sig7493/bocAPI/api/es"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/estransport"
	"github.com/tidwall/gjson"
)

var (
    r  map[string]interface{}
    wg sync.WaitGroup
  )

func getEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func es_connect() (*elasticsearch.Client, error) {
	eshosts := getEnvVariable("ES_HOST")
	fmt.Printf("%v\n", eshosts)
	cfg := elasticsearch.Config{
		Addresses: []string{eshosts},
		Logger: &estransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err!= nil {
		log.Fatal("Unable to connect to elasticsearch \n %s", err)
	}
	return es, err
}

func get_es_cluster_info() {

	es,_ := es_connect()
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s\n", elasticsearch.Version)
	log.Printf("Server: %s\n", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
}

func GetByProcessRunDateId_search_request(processrundateid int) *GetByProcessRunDateIdResponse {

	resp := GetByProcessRunDateIdResponse{}
	es,_ := es_connect()

	 // Build the request body.
	 var buf bytes.Buffer
	 query := map[string]interface{}{
	   "query": map[string]interface{}{
		 "match": map[string]interface{}{
		   "PROCESS_RUN_DATE_ID": processrundateid,
		 },
	   },
	 }

	 if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	  }
	
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("banknotes_aggregate-*"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		log.Fatalf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  log.Fatalf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  log.Fatalf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
	  log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	  // Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64)))
		resp.ProcessRunDateID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		resp.Count = resp.Count + int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))
	}

	log.Println(strings.Repeat("=", 37))

	return &resp

}


func GetBetweenProcessRunDateIds_search_request(fromprocessrundateid int, toprocessrundateid int) *GetBetweenProcessRunDateIdsResponse {

	resp := GetBetweenProcessRunDateIdsResponse{}
	es,_ := es_connect()

	 // Build the request body.
	 var buf bytes.Buffer
	 query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"PROCESS_RUN_DATE_ID": map[string]interface{}{
								"gte": fromprocessrundateid,
								"lte": toprocessrundateid,
							},
						},
					},
				},
			},
		},
	}

	 if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	  }
	
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("banknotes_aggregate-*"),
		es.Search.WithBody(&buf),
		es.Search.WithSize(200),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		log.Fatalf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  log.Fatalf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  log.Fatalf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
	  log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	  // Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64)))
		//resp.ProcessRunDateID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		//resp.Count = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))
		resp = append(resp, GetBetweenProcessRunDateIds{int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))})
	}

	log.Println(strings.Repeat("=", 37))

	return &resp
}


func GetNotesValidityDetails_search_request() *GetNotesValidityDetailsResponse {
	resp := GetNotesValidityDetailsResponse{}
	es,_ := es_connect()

	// var compositeSize = 2000
	var buf bytes.Buffer

	var Denomination = ""
	var Image_path = ""
	var Rgb_color = ""
	var Rgb_val = ""
	var Serial_number = ""
	var Process_Run_Date_ID = 0
	var Bps_Shift_ID = 0
	var Machine_ID = 0
	var Print_Batch_ID = 0
	var Rdp_ID = 0
	var Bn_Status_ID = 0
	var Output_Stacker_ID = 0
	var Circ_Trial_ID = 0
	var Bps_Shift_Nb = 0
	var Deposit_Nb = 0
	var Row_Counter_NB = 0
	var Load_ID = 0
	

	query := map[string]interface{}{
		"query": map[string]interface{}{
		  "exists": map[string]interface{}{
			  "field": "SER_NUM",
		  },
		},"collapse": map[string]interface{}{
			"field": "SER_NUM.keyword",
		  },"size": 4000,
	  }
 
	  if err := json.NewEncoder(&buf).Encode(query); err != nil {
		 log.Fatalf("Error encoding query: %s", err)
	   }
	 
	 // Perform the search request.
	 res, err := es.Search(
		 es.Search.WithContext(context.Background()),
		 es.Search.WithIndex("notes_integration_index-*"),
		 es.Search.WithBody(&buf),
		 es.Search.WithTrackTotalHits(true),
		 es.Search.WithPretty(),
	   )
	   if err != nil {
		 log.Fatalf("Error getting response: %s", err)
	   }
	   defer res.Body.Close()
 
	   if res.IsError() {
		 var e map[string]interface{}
		 if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		   log.Fatalf("Error parsing the response body: %s", err)
		 } else {
		   // Print the response status and error information.
		   log.Fatalf("[%s] %s: %s",
			 res.Status(),
			 e["error"].(map[string]interface{})["type"],
			 e["error"].(map[string]interface{})["reason"],
		   )
		 }
	   }
 
	   if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		 log.Fatalf("Error parsing the response body: %s", err)
	   }
	   // Print the response status, number of results, and request duration.
	   log.Printf(
		 "[%s] %d hits; took: %dms",
		 res.Status(),
		 int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		 int(r["took"].(float64)),
	   )
 
	   // Print the ID and document source for each hit.
	 for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		 log.Printf(" serial_number=%s, image_path=%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"], hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"])
		 
		 Denomination = hit.(map[string]interface{})["_source"].(map[string]interface{})["denomination"].(string)
		 Image_path = hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"].(string)
		 Rgb_color = hit.(map[string]interface{})["_source"].(map[string]interface{})["rgb_color"].(string)
		 Rgb_val = hit.(map[string]interface{})["_source"].(map[string]interface{})["rgb_val"].(string)
		 Serial_number = hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"].(string)

		 Process_Run_Date_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		 Bps_Shift_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"].(float64))
		 Machine_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"].(float64))
		 Print_Batch_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(float64))
		 Rdp_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"].(float64))
		 Bn_Status_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"].(float64))
		 Output_Stacker_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"].(float64))
		 Circ_Trial_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"].(float64))
		 Bps_Shift_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"].(float64))
		 Deposit_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"].(float64))
		 Row_Counter_NB = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"].(float64))
		 Load_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"].(float64))
		 
/* 		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"] == nil {
			Process_Run_Date_ID = 0
		 } else {
			Process_Run_Date_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		 }

		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"] == nil {
			Bps_Shift_ID = 0
		 } else {
			 Bps_Shift_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"].(float64))
		 }

		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"] == nil {
			Machine_ID = 0
		} else {
			Machine_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"].(float64))
		}
		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"] == nil {
			Print_Batch_ID = 0
			} else {
			Print_Batch_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"] == nil {
				Rdp_ID = 0
			} else {
				Rdp_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"] == nil {
				Bn_Status_ID = 0
			} else {
				Bn_Status_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"] == nil {
				Output_Stacker_ID = 0
			} else {
				Output_Stacker_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"] == nil {
				Circ_Trial_ID = 0
			} else {
				Circ_Trial_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"] == nil {
				Bps_Shift_Nb = 0
			} else {
				Bps_Shift_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"] == nil {
				Deposit_Nb = 0
			} else {
				Deposit_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"] == nil {
				Row_Counter_NB = 0
			} else {
				Row_Counter_NB = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"] == nil {
				Load_ID = 0
			} else {
				Load_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"].(float64))
			} */	
		 
		
		
		 resp = append(resp, GetNotesValidityDetails{Denomination, Image_path, Rgb_color, Rgb_val, Serial_number, Process_Run_Date_ID, Bps_Shift_ID, Machine_ID, Print_Batch_ID, Rdp_ID, Bn_Status_ID, Output_Stacker_ID, Circ_Trial_ID, Bps_Shift_Nb, Deposit_Nb, Row_Counter_NB, Load_ID})
	 }
 
	 log.Println(strings.Repeat("=", 37))
 
	 return &resp

}

func GetNotesInValidityDetails_search_request() *GetNotesValidityDetailsResponse {
	resp := GetNotesValidityDetailsResponse{}
	es,_ := es_connect()

	// var compositeSize = 2000
	var buf bytes.Buffer

	var Denomination = ""
	var Image_path = ""
	var Rgb_color = ""
	var Rgb_val = ""
	var Serial_number = ""
	var Process_Run_Date_ID = 0
	var Bps_Shift_ID = 0
	var Machine_ID = 0
	var Print_Batch_ID = 0
	var Rdp_ID = 0
	var Bn_Status_ID = 0
	var Output_Stacker_ID = 0
	var Circ_Trial_ID = 0
	var Bps_Shift_Nb = 0
	var Deposit_Nb = 0
	var Row_Counter_NB = 0
	var Load_ID = 0
	

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": map[string]interface{}{
					"exists": map[string]interface{}{
						"field": "SER_NUM",
						},
				},
			},
		},"collapse": map[string]interface{}{
			"field": "serial_number.keyword",
		  },"size": 4000,
	  }
 
	  if err := json.NewEncoder(&buf).Encode(query); err != nil {
		 log.Fatalf("Error encoding query: %s", err)
	   }
	 
	 // Perform the search request.
	 res, err := es.Search(
		 es.Search.WithContext(context.Background()),
		 es.Search.WithIndex("notes_integration_index-*"),
		 es.Search.WithBody(&buf),
		 es.Search.WithTrackTotalHits(true),
		 es.Search.WithPretty(),
	   )
	   if err != nil {
		 log.Fatalf("Error getting response: %s", err)
	   }
	   defer res.Body.Close()
 
	   if res.IsError() {
		 var e map[string]interface{}
		 if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		   log.Fatalf("Error parsing the response body: %s", err)
		 } else {
		   // Print the response status and error information.
		   log.Fatalf("[%s] %s: %s",
			 res.Status(),
			 e["error"].(map[string]interface{})["type"],
			 e["error"].(map[string]interface{})["reason"],
		   )
		 }
	   }
 
	   if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		 log.Fatalf("Error parsing the response body: %s", err)
	   }
	   // Print the response status, number of results, and request duration.
	   log.Printf(
		 "[%s] %d hits; took: %dms",
		 res.Status(),
		 int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		 int(r["took"].(float64)),
	   )
 
	   // Print the ID and document source for each hit.
	 for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		 log.Printf(" serial_number=%s, image_path=%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"], hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"])
		 
		 Denomination = hit.(map[string]interface{})["_source"].(map[string]interface{})["denomination"].(string)
		 Image_path = hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"].(string)
		 Rgb_color = hit.(map[string]interface{})["_source"].(map[string]interface{})["rgb_color"].(string)
		 Rgb_val = hit.(map[string]interface{})["_source"].(map[string]interface{})["rgb_val"].(string)
		 Serial_number = hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"].(string)

		 
 		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"] == nil {
			Process_Run_Date_ID = 0
		 } else {
			Process_Run_Date_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		 }

		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"] == nil {
			Bps_Shift_ID = 0
		 } else {
			 Bps_Shift_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"].(float64))
		 }

		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"] == nil {
			Machine_ID = 0
		} else {
			Machine_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"].(float64))
		}
		 if hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"] == nil {
			Print_Batch_ID = 0
			} else {
			Print_Batch_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"] == nil {
				Rdp_ID = 0
			} else {
				Rdp_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"] == nil {
				Bn_Status_ID = 0
			} else {
				Bn_Status_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"] == nil {
				Output_Stacker_ID = 0
			} else {
				Output_Stacker_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"] == nil {
				Circ_Trial_ID = 0
			} else {
				Circ_Trial_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"] == nil {
				Bps_Shift_Nb = 0
			} else {
				Bps_Shift_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"] == nil {
				Deposit_Nb = 0
			} else {
				Deposit_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"] == nil {
				Row_Counter_NB = 0
			} else {
				Row_Counter_NB = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"].(float64))
			}
		if hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"] == nil {
				Load_ID = 0
			} else {
				Load_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"].(float64))
			} 	
		 resp = append(resp, GetNotesValidityDetails{Denomination, Image_path, Rgb_color, Rgb_val, Serial_number, Process_Run_Date_ID, Bps_Shift_ID, Machine_ID, Print_Batch_ID, Rdp_ID, Bn_Status_ID, Output_Stacker_ID, Circ_Trial_ID, Bps_Shift_Nb, Deposit_Nb, Row_Counter_NB, Load_ID})
	 }
 
	 log.Println(strings.Repeat("=", 37))
 
	 return &resp

}



func GetAllProcessRunDateIds_search_request() *GetBetweenProcessRunDateIdsResponse {

	resp := GetBetweenProcessRunDateIdsResponse{}
	es,_ := es_connect()

	var compositeSize = 1000

	field := "PROCESS_RUN_DATE_ID"

	 // Build the request body.
	 var buf bytes.Buffer

	 // This query didn't work when there are multiple sources in it since it is an array of sources in ES. Unable to pass the array of sources
/*  		query := map[string]interface{}{
			"size": 0,
			"aggs": map[string]interface{}{
				"ALL_PROCESS_RUN_DATE_ID": map[string]interface{}{
					"composite": map[string]interface{}{
						"size": compositeSize,
						"sources": map[string]interface{}{
							field: map[string]interface{}{
								"terms": map[string]interface{}{
									"field": field,
									},
								},
							"count": map[string]interface{}{
									"terms": map[string]interface{}{
										"field": "count",
										},
									},
					},
				},
			},
		},
	} */
	// This query works as the second source is moved to aggs
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"ALL_PROCESS_RUN_DATE_ID": map[string]interface{}{
				"composite": map[string]interface{}{
					"size": compositeSize,
					"sources": map[string]interface{}{
						field: map[string]interface{}{
							"terms": map[string]interface{}{
								"field": field,
								},
							},
				},
			},
			"aggs": map[string]interface{}{
				"count": map[string]interface{}{ 
					"max": map[string]interface{}{ 
						"field": "count", 
						}, 
					},
			  },	
		},
	},
}

	 if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	  }
	
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("banknotes_aggregate-*"),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		log.Fatalf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  log.Fatalf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  log.Fatalf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
		  log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	if int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)) > 0 {  
	for _, agg := range r["aggregations"].(map[string]interface{})["ALL_PROCESS_RUN_DATE_ID"].(map[string]interface{})["buckets"].([]interface{}){
		////log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64)))
		//log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)))
		resp = append(resp, GetBetweenProcessRunDateIds{int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64))})
		}
	}
	

	log.Println(strings.Repeat("=", 37))

	return &resp
}

func GetNotesDestructionAgg_search_request() *GetNotesDestructionAggResponse {
	resp := GetNotesDestructionAggResponse{}

	es,_ := es_connect()

	// Build the request body.
	var buf bytes.Buffer
	//query := map[string]interface{}{"query": map[string]interface{}{"match_all": map[string]interface{}{}}}

	query := map[string]interface{}{
		"sort": []map[string]interface{}{
				{
				  "PRINT_BATCH_ID.keyword": map[string]interface{}{
					"order": "asc",
				  },
				 "SUM_ROW_COUNTER_NB": map[string]interface{}{
					"order": "desc",
				  }, 
				},
			  },"query": map[string]interface{}{"match_all": map[string]interface{}{}}}

   if err := json.NewEncoder(&buf).Encode(query); err != nil {
	log.Fatalf("Error encoding query: %s", err)
	  }
	  
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		//es.Search.WithIndex("velocity_notes_destructions_print_batch_id*"),
		es.Search.WithIndex("velocity_notes_destructions_all*"),
		es.Search.WithBody(&buf),
		//es.Search.WithSort("PRINT_BATCH_ID.keyword"),
		es.Search.WithSize(500),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	if err != nil {
	log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	} else {
		// Print the response status and error information.
		log.Fatalf("[%s] %s: %s",
		res.Status(),
		e["error"].(map[string]interface{})["type"],
		e["error"].(map[string]interface{})["reason"],
		)
	}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	// Print the PRINT_BATCH_ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	log.Printf(" PRINT_BATCH_ID=%s, SUM_ROW_COUNTER_NB=%d", hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"], int(hit.(map[string]interface{})["_source"].(map[string]interface{})["SUM_ROW_COUNTER_NB"].(float64)))
	resp = append(resp, GetNotesDestructionAgg{hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_DENOM_EN_NM"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["T_EXCEED_STRING_TXT"].(string),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["YEAR_NB"].(float64)),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["QUARTER_EN_NM"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["MONTH_EN_NM"].(string),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["SUM_ROW_COUNTER_NB"].(float64)),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["rank"].(float64))})
	}

	log.Println(strings.Repeat("=", 37))

	return &resp
}

//func GetNotesDestructionDetails_search_request(print_batch_id string, year int, quarter string, month string, denomination string) *GetNotesDestructionDetailsResponse {
func GetNotesDestructionDetails_search_request(print_batch_id string, year int, quarter string, month string, denomination string, from int, scroll_ID string) *GetNotesDestructionDetailsResponse {

	var (
		batchNum int
		scrollID string
		pbi string
		denom_en_nm string
		exceed_string_txt string
		yr_nb int
		qtr_en_nm string
		mnth_en_nm string
		row_cnter_nb int
	)
	resp := GetNotesDestructionDetailsResponse{}
	es,_ := es_connect()

	if scroll_ID == "None" {
		// Build the request body.
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"PRINT_BATCH_ID": print_batch_id,
							},
						},
						{
							"match": map[string]interface{}{
								"YEAR_NB": year,
							},
						},
						{
							"match": map[string]interface{}{
								"QUARTER_EN_NM": quarter,
							},
						},
						{
							"match": map[string]interface{}{
								"MONTH_EN_NM": month,
							},
						},
						{
							"match": map[string]interface{}{
								"BN_DENOM_EN_NM": denomination,
							},
						},
					},
				},
			},"_source": false,
			"fields": []string {"YEAR_NB","QUARTER_EN_NM","MONTH_EN_NM","T_EXCEED_STRING_TXT","BN_DENOM_EN_NM","PRINT_BATCH_ID","ROW_COUNTER_NB"},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("velocity_print_batch_id_" + print_batch_id + "*"),
			es.Search.WithBody(&buf),
			es.Search.WithFrom(from),
			es.Search.WithSize(20),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
			es.Search.WithScroll(time.Minute),
		)
		if err != nil {
		log.Fatalf("Error getting response: %s", err)
		}

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
				)
			}
			}

		// Handle the first batch of data and extract the scrollID
		//
		json := read(res.Body)
		res.Body.Close()

		scrollID = gjson.Get(json, "_scroll_id").String()
		content := gjson.Get(json, "hits.hits.#.fields").String()
		contentResult := gjson.Get(json, "hits.hits.#.fields")

		log.Println("Batch   ", batchNum)
		log.Println("ScrollID", scrollID)
		log.Println("Content", content)

		/* for idx, name := range contentResult.Array(){
			log.Println("Printing name...", idx)
			log.Println(name.String())	
		} */

		
		contentResult.ForEach(func(key, value gjson.Result) bool {
			log.Println(value)
			batchNum++
			log.Println(batchNum)
			value.ForEach(func(key1, value1 gjson.Result) bool {
				//log.Println("Key1/Value1")
				//log.Println(key1.String())
				//log.Println(value1.Array()[0])
				
				if key1.String() == "PRINT_BATCH_ID" {
					pbi = value1.Array()[0].String()
				}
				if key1.String() == "BN_DENOM_EN_NM" {
					denom_en_nm = value1.Array()[0].String()
				}
				if key1.String() == "T_EXCEED_STRING_TXT" {
					exceed_string_txt = value1.Array()[0].String()
				}
				if key1.String() == "YEAR_NB" {
					yr_nb,_ = strconv.Atoi(value1.Array()[0].String())
				}
				if key1.String() == "QUARTER_EN_NM" {
					qtr_en_nm = value1.Array()[0].String()
				}
				if key1.String() == "MONTH_EN_NM" {
					mnth_en_nm = value1.Array()[0].String()
				}
				if key1.String() == "ROW_COUNTER_NB" {
					row_cnter_nb,_ = strconv.Atoi(value1.Array()[0].String())
				}
				
				return true
			})
			resp = append(resp, GetNotesDestructionDetails{scrollID, pbi, denom_en_nm, exceed_string_txt, yr_nb, qtr_en_nm, mnth_en_nm, row_cnter_nb})
			return true // keep iterating
		})

		/* log.Println("PRINT_BATCH_ID = ", gjson.Get(json, "hits.hits.#.fields.PRINT_BATCH_ID"))
		log.Println("BN_DENOM_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.BN_DENOM_EN_NM"))
		log.Println("T_EXCEED_STRING_TXT = ", gjson.Get(json, "hits.hits.#.fields.T_EXCEED_STRING_TXT"))
		log.Println("YEAR_NB = ", gjson.Get(json, "hits.hits.#.fields.YEAR_NB"))
		log.Println("QUARTER_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.QUARTER_EN_NM"))
		log.Println("MONTH_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.MONTH_EN_NM"))
		log.Println("ROW_COUNTER_NB = ", gjson.Get(json, "hits.hits.#.fields.ROW_COUNTER_NB")) */
		log.Println(strings.Repeat("-", 80))

		//resp = append(resp, GetNotesDestructionDetails{scrollID, content})

	} else {
		scrollID = scroll_ID
		// for {
		//	batchNum++

			// Perform the scroll request and pass the scrollID and scroll duration
			res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
			if err != nil {
				log.Fatalf("Error: %s", err)
			}
			if res.IsError() {
				log.Fatalf("Error response: %s", res)
			}

			json := read(res.Body)
			res.Body.Close()
			// Extract the scrollID from response
			//
			scrollID = gjson.Get(json, "_scroll_id").String()
			// Extract the search results
			//
			hits := gjson.Get(json, "hits.hits")
			content := gjson.Get(hits.Raw, "#.fields")

			// Break out of the loop when there are no results
			//
			if len(hits.Array()) < 1 {
				log.Println("Finished scrolling")
				//break
			} else {
				log.Println("Batch   ", batchNum)
				log.Println("ScrollID", scrollID)
				log.Println("content ", content)
				content.ForEach(func(key, value gjson.Result) bool {
					value.ForEach(func(key1, value1 gjson.Result) bool {
						//log.Println("Key1/Value1")
						//log.Println(key1.String())
						//log.Println(value1.Array()[0])
						if key1.String() == "PRINT_BATCH_ID" {
							pbi = value1.Array()[0].String()
						}
						if key1.String() == "BN_DENOM_EN_NM" {
							denom_en_nm = value1.Array()[0].String()
						}
						if key1.String() == "T_EXCEED_STRING_TXT" {
							exceed_string_txt = value1.Array()[0].String()
						}
						if key1.String() == "YEAR_NB" {
							yr_nb,_ = strconv.Atoi(value1.Array()[0].String())
						}
						if key1.String() == "QUARTER_EN_NM" {
							qtr_en_nm = value1.Array()[0].String()
						}
						if key1.String() == "MONTH_EN_NM" {
							mnth_en_nm = value1.Array()[0].String()
						}
						if key1.String() == "ROW_COUNTER_NB" {
							row_cnter_nb,_ = strconv.Atoi(value1.Array()[0].String())
						}
						
						return true
					})
					resp = append(resp, GetNotesDestructionDetails{scrollID, pbi, denom_en_nm, exceed_string_txt, yr_nb, qtr_en_nm, mnth_en_nm, row_cnter_nb})
					return true // keep iterating
				})
				/* log.Println("PRINT_BATCH_ID = ", gjson.Get(hits.Raw, "#.fields.PRINT_BATCH_ID"))
				log.Println("BN_DENOM_EN_NM = ", gjson.Get(hits.Raw, "#.fields.BN_DENOM_EN_NM"))
				log.Println("T_EXCEED_STRING_TXT = ", gjson.Get(hits.Raw, "#.fields.T_EXCEED_STRING_TXT"))
				log.Println("YEAR_NB = ", gjson.Get(hits.Raw, "#.fields.YEAR_NB"))
				log.Println("QUARTER_EN_NM = ", gjson.Get(hits.Raw, "#.fields.QUARTER_EN_NM"))
				log.Println("MONTH_EN_NM = ", gjson.Get(hits.Raw, "#.fields.MONTH_EN_NM"))
				log.Println("ROW_COUNTER_NB = ", gjson.Get(hits.Raw, "#.fields.ROW_COUNTER_NB")) */
				log.Println(strings.Repeat("-", 80))
				//resp = append(resp, GetNotesDestructionDetails{scrollID, content})
			}

		//}
	
	}
	/* if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		fields := doc["fields"]

		fmt.Println(fields)
		fmt.Println(reflect.TypeOf(fields))
		pbi := fields.(map[string]interface{})["PRINT_BATCH_ID"]
		fmt.Println(reflect.TypeOf(pbi))
		//fmt.Println(reflect.TypeOf(reflect.ValueOf(pbi)))
		fmt.Println(pbi)
		pbi_s := make([]string, len(pbi.([]interface{})))
		for i, v := range pbi.([]interface{}) {
			pbi_s[i] = fmt.Sprint(v)
		}
		fmt.Println(pbi_s[0])

		denom_en := fields.(map[string]interface{})["BN_DENOM_EN_NM"]
		denom_en_s := make([]string, len(denom_en.([]interface{})))
		for i, v := range denom_en.([]interface{}) {
			denom_en_s[i] = fmt.Sprint(v)
		}

		t_exceed := fields.(map[string]interface{})["T_EXCEED_STRING_TXT"]
		t_exceed_s := make([]string, len(t_exceed.([]interface{})))
		for i, v := range t_exceed.([]interface{}) {
			t_exceed_s[i] = fmt.Sprint(v)
		}

		year_nb := fields.(map[string]interface{})["YEAR_NB"]
		year_nb_s := make([]string, len(year_nb.([]interface{})))
		for i, v := range year_nb.([]interface{}) {
			year_nb_s[i] = fmt.Sprint(v)
		}

		year_nb_i, err := strconv.Atoi(year_nb_s[0])
		if err != nil {
			fmt.Printf("%v", err)
			
		}

		qtr_en := fields.(map[string]interface{})["QUARTER_EN_NM"]
		qtr_en_s := make([]string, len(qtr_en.([]interface{})))
		for i, v := range qtr_en.([]interface{}) {
			qtr_en_s[i] = fmt.Sprint(v)
		}

		month_en := fields.(map[string]interface{})["MONTH_EN_NM"]
		month_en_s := make([]string, len(month_en.([]interface{})))
		for i, v := range month_en.([]interface{}) {
			month_en_s[i] = fmt.Sprint(v)
		}

		row_ctr := fields.(map[string]interface{})["ROW_COUNTER_NB"]
		row_ctr_s := make([]string, len(row_ctr.([]interface{})))
		for i, v := range row_ctr.([]interface{}) {
			row_ctr_s[i] = fmt.Sprint(v)
		}

		row_ctr_i, err := strconv.Atoi(row_ctr_s[0])
		if err != nil {
			fmt.Printf("%v", err)
			
		}  

		fmt.Println("printed fields")
		resp = append(resp, GetNotesDestructionDetails{pbi_s[0], denom_en_s[0],t_exceed_s[0],year_nb_i,qtr_en_s[0],month_en_s[0],row_ctr_i})
	} */

	return &resp
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}