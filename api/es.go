package api

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"sync"
	"bytes"
	"context"
	
	"strings"

	"github.com/joho/godotenv"

	//"github.com/sig7493/bocAPI/api/es"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/estransport"
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