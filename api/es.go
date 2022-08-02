package api

import (
	"os"
	"io"
	//"fmt"
	"fmt"
	"encoding/json"
	"sync"
	"bytes"
	"context"
	"reflect"
	"strings"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	//"github.com/sig7493/bocAPI/api/es"

	"github.com/elastic/go-elasticsearch/v8"
	// "github.com/elastic/go-elasticsearch/v8/estransport"
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
	  fmt.Printf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func es_connect() (*elasticsearch.Client, error) {
	eshosts := getEnvVariable("ES_HOST")
	fmt.Printf("%v\n", eshosts)
	cfg := elasticsearch.Config{
		Addresses: []string{eshosts},
		/* Logger: &estransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}, */
	}
	es, err := elasticsearch.NewClient(cfg)

	if err!= nil {
		fmt.Printf("Unable to connect to elasticsearch \n %s", err)
	}
	return es, err
}

func get_es_cluster_info() {

	es,_ := es_connect()
	res, err := es.Info()
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		fmt.Printf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	fmt.Printf("Client: %s\n", elasticsearch.Version)
	fmt.Printf("Server: %s\n", r["version"].(map[string]interface{})["number"])
	fmt.Println(strings.Repeat("~", 37))
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
		fmt.Printf("Error encoding query: %s", err)
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
		fmt.Printf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  fmt.Printf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  fmt.Printf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
	  fmt.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	  // Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		fmt.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64)))
		resp.ProcessRunDateID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		resp.Count = resp.Count + int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))
	}

	fmt.Println(strings.Repeat("=", 37))

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
		fmt.Printf("Error encoding query: %s", err)
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
		fmt.Printf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  fmt.Printf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  fmt.Printf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
	  fmt.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	  // Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		fmt.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64)))
		//resp.ProcessRunDateID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		//resp.Count = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))
		resp = append(resp, GetBetweenProcessRunDateIds{int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))})
	}

	fmt.Println(strings.Repeat("=", 37))

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
		 fmt.Printf("Error encoding query: %s", err)
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
		 fmt.Printf("Error getting response: %s", err)
	   }
	   defer res.Body.Close()
 
	   if res.IsError() {
		 var e map[string]interface{}
		 if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		   fmt.Printf("Error parsing the response body: %s", err)
		 } else {
		   // Print the response status and error information.
		   fmt.Printf("[%s] %s: %s",
			 res.Status(),
			 e["error"].(map[string]interface{})["type"],
			 e["error"].(map[string]interface{})["reason"],
		   )
		 }
	   }
 
	   if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		 fmt.Printf("Error parsing the response body: %s", err)
	   }
	   // Print the response status, number of results, and request duration.
	   fmt.Printf(
		 "[%s] %d hits; took: %dms",
		 res.Status(),
		 int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		 int(r["took"].(float64)),
	   )
 
	   // Print the ID and document source for each hit.
	 for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		 fmt.Printf(" serial_number=%s, image_path=%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"], hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"])
		 
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
		 
		//  if hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"] == nil {
		// 	Process_Run_Date_ID = 0
		//  } else {
		// 	Process_Run_Date_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64))
		//  }

		//  if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"] == nil {
		// 	Bps_Shift_ID = 0
		//  } else {
		// 	 Bps_Shift_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_ID"].(float64))
		//  }

		//  if hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"] == nil {
		// 	Machine_ID = 0
		// } else {
		// 	Machine_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["MACHINE_ID"].(float64))
		// }
		//  if hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"] == nil {
		// 	Print_Batch_ID = 0
		// 	} else {
		// 	Print_Batch_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"] == nil {
		// 		Rdp_ID = 0
		// 	} else {
		// 		Rdp_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["RDP_ID"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"] == nil {
		// 		Bn_Status_ID = 0
		// 	} else {
		// 		Bn_Status_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_STATUS_ID"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"] == nil {
		// 		Output_Stacker_ID = 0
		// 	} else {
		// 		Output_Stacker_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["OUTPUT_STACKER_ID"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"] == nil {
		// 		Circ_Trial_ID = 0
		// 	} else {
		// 		Circ_Trial_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["CIRC_TRIAL_ID"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"] == nil {
		// 		Bps_Shift_Nb = 0
		// 	} else {
		// 		Bps_Shift_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["BPS_SHIFT_NB"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"] == nil {
		// 		Deposit_Nb = 0
		// 	} else {
		// 		Deposit_Nb = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["DEPOSIT_NB"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"] == nil {
		// 		Row_Counter_NB = 0
		// 	} else {
		// 		Row_Counter_NB = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["ROW_COUNTER_NB"].(float64))
		// 	}
		// if hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"] == nil {
		// 		Load_ID = 0
		// 	} else {
		// 		Load_ID = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["LOAD_ID"].(float64))
		// 	}	
		 
		
		
		 resp = append(resp, GetNotesValidityDetails{Denomination, Image_path, Rgb_color, Rgb_val, Serial_number, Process_Run_Date_ID, Bps_Shift_ID, Machine_ID, Print_Batch_ID, Rdp_ID, Bn_Status_ID, Output_Stacker_ID, Circ_Trial_ID, Bps_Shift_Nb, Deposit_Nb, Row_Counter_NB, Load_ID})
	 }
 
	 fmt.Println(strings.Repeat("=", 37))
 
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
		 fmt.Printf("Error encoding query: %s", err)
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
		 fmt.Printf("Error getting response: %s", err)
	   }
	   defer res.Body.Close()
 
	   if res.IsError() {
		 var e map[string]interface{}
		 if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		   fmt.Printf("Error parsing the response body: %s", err)
		 } else {
		   // Print the response status and error information.
		   fmt.Printf("[%s] %s: %s",
			 res.Status(),
			 e["error"].(map[string]interface{})["type"],
			 e["error"].(map[string]interface{})["reason"],
		   )
		 }
	   }
 
	   if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		 fmt.Printf("Error parsing the response body: %s", err)
	   }
	   // Print the response status, number of results, and request duration.
	   fmt.Printf(
		 "[%s] %d hits; took: %dms",
		 res.Status(),
		 int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		 int(r["took"].(float64)),
	   )
 
	   // Print the ID and document source for each hit.
	 for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		 fmt.Printf(" serial_number=%s, image_path=%s", hit.(map[string]interface{})["_source"].(map[string]interface{})["serial_number"], hit.(map[string]interface{})["_source"].(map[string]interface{})["image_path"])
		 
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
 
	 fmt.Println(strings.Repeat("=", 37))
 
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
 	// 	query := map[string]interface{}{
	// 		"size": 0,
	// 		"aggs": map[string]interface{}{
	// 			"ALL_PROCESS_RUN_DATE_ID": map[string]interface{}{
	// 				"composite": map[string]interface{}{
	// 					"size": compositeSize,
	// 					"sources": map[string]interface{}{
	// 						field: map[string]interface{}{
	// 							"terms": map[string]interface{}{
	// 								"field": field,
	// 								},
	// 							},
	// 						"count": map[string]interface{}{
	// 								"terms": map[string]interface{}{
	// 									"field": "count",
	// 									},
	// 								},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
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
		fmt.Printf("Error encoding query: %s", err)
	  }
	
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("banknotes_aggregate-*"),
		es.Search.WithBody(&buf),
		es.Search.WithPretty(),
	  )
	  if err != nil {
		fmt.Printf("Error getting response: %s", err)
	  }
	  defer res.Body.Close()

	  if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		  fmt.Printf("Error parsing the response body: %s", err)
		} else {
		  // Print the response status and error information.
		  fmt.Printf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		  )
		}
	  }

	  if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	  }
	  // Print the response status, number of results, and request duration.
		  fmt.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	  )

	if int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)) > 0 {  
	for _, agg := range r["aggregations"].(map[string]interface{})["ALL_PROCESS_RUN_DATE_ID"].(map[string]interface{})["buckets"].([]interface{}){
		////fmt.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64)))
		//fmt.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)))
		resp = append(resp, GetBetweenProcessRunDateIds{int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64))})
		}
	}
	

	fmt.Println(strings.Repeat("=", 37))

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
	fmt.Printf("Error encoding query: %s", err)
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
	fmt.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	} else {
		// Print the response status and error information.
		fmt.Printf("[%s] %s: %s",
		res.Status(),
		e["error"].(map[string]interface{})["type"],
		e["error"].(map[string]interface{})["reason"],
		)
	}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	fmt.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	// Print the PRINT_BATCH_ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	fmt.Printf(" PRINT_BATCH_ID=%s, SUM_ROW_COUNTER_NB=%d", hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"], int(hit.(map[string]interface{})["_source"].(map[string]interface{})["SUM_ROW_COUNTER_NB"].(float64)))
	resp = append(resp, GetNotesDestructionAgg{hit.(map[string]interface{})["_source"].(map[string]interface{})["PRINT_BATCH_ID"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["BN_DENOM_EN_NM"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["T_EXCEED_STRING_TXT"].(string),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["YEAR_NB"].(float64)),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["QUARTER_EN_NM"].(string),
	 hit.(map[string]interface{})["_source"].(map[string]interface{})["MONTH_EN_NM"].(string),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["SUM_ROW_COUNTER_NB"].(float64)),
	 int(hit.(map[string]interface{})["_source"].(map[string]interface{})["rank"].(float64))})
	}

	fmt.Println(strings.Repeat("=", 37))

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
		fmt.Println("<----------------------------------------------->")
		fmt.Println(query)
		fmt.Println("<----------------------------------------------->")
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("velocity_print_batch_id_" + print_batch_id + "*"),
			es.Search.WithBody(&buf),
			es.Search.WithFrom(from),
			es.Search.WithSize(500),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
			es.Search.WithScroll(time.Minute * 10),
		)
		if err != nil {
		fmt.Printf("Error getting response: %s", err)
		}

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				fmt.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				fmt.Printf("[%s] %s: %s",
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

		fmt.Println("Batch   ", batchNum)
		fmt.Println("ScrollID", scrollID)
		fmt.Println("Content", content)

		// for idx, name := range contentResult.Array(){
		// 	fmt.Println("Printing name...", idx)
		// 	fmt.Println(name.String())	
		// }

		
		contentResult.ForEach(func(key, value gjson.Result) bool {
			fmt.Println(value)
			batchNum++
			fmt.Println(batchNum)
			value.ForEach(func(key1, value1 gjson.Result) bool {
				//fmt.Println("Key1/Value1")
				//fmt.Println(key1.String())
				//fmt.Println(value1.Array()[0])
				
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

		// fmt.Println("PRINT_BATCH_ID = ", gjson.Get(json, "hits.hits.#.fields.PRINT_BATCH_ID"))
		// fmt.Println("BN_DENOM_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.BN_DENOM_EN_NM"))
		// fmt.Println("T_EXCEED_STRING_TXT = ", gjson.Get(json, "hits.hits.#.fields.T_EXCEED_STRING_TXT"))
		// fmt.Println("YEAR_NB = ", gjson.Get(json, "hits.hits.#.fields.YEAR_NB"))
		// fmt.Println("QUARTER_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.QUARTER_EN_NM"))
		// fmt.Println("MONTH_EN_NM = ", gjson.Get(json, "hits.hits.#.fields.MONTH_EN_NM"))
		// fmt.Println("ROW_COUNTER_NB = ", gjson.Get(json, "hits.hits.#.fields.ROW_COUNTER_NB"))
		fmt.Println(strings.Repeat("-", 80))

		//resp = append(resp, GetNotesDestructionDetails{scrollID, content})

	} else {
		scrollID = scroll_ID
		// for {
		//	batchNum++

			// Perform the scroll request and pass the scrollID and scroll duration
			res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
			if err != nil {
				fmt.Printf("Error: %s", err)
			}
			if res.IsError() {
				fmt.Printf("Error response: %s", res)
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
				fmt.Println("Finished scrolling")
				//break
			} else {
				fmt.Println("Batch   ", batchNum)
				fmt.Println("ScrollID", scrollID)
				fmt.Println("content ", content)
				content.ForEach(func(key, value gjson.Result) bool {
					value.ForEach(func(key1, value1 gjson.Result) bool {
						//fmt.Println("Key1/Value1")
						//fmt.Println(key1.String())
						//fmt.Println(value1.Array()[0])
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
				// fmt.Println("PRINT_BATCH_ID = ", gjson.Get(hits.Raw, "#.fields.PRINT_BATCH_ID"))
				// fmt.Println("BN_DENOM_EN_NM = ", gjson.Get(hits.Raw, "#.fields.BN_DENOM_EN_NM"))
				// fmt.Println("T_EXCEED_STRING_TXT = ", gjson.Get(hits.Raw, "#.fields.T_EXCEED_STRING_TXT"))
				// fmt.Println("YEAR_NB = ", gjson.Get(hits.Raw, "#.fields.YEAR_NB"))
				// fmt.Println("QUARTER_EN_NM = ", gjson.Get(hits.Raw, "#.fields.QUARTER_EN_NM"))
				// fmt.Println("MONTH_EN_NM = ", gjson.Get(hits.Raw, "#.fields.MONTH_EN_NM"))
				// fmt.Println("ROW_COUNTER_NB = ", gjson.Get(hits.Raw, "#.fields.ROW_COUNTER_NB"))
				fmt.Println(strings.Repeat("-", 80))
				//resp = append(resp, GetNotesDestructionDetails{scrollID, content})
			}

		//}
	
	}
	// if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 	fmt.Printf("Error parsing the response body: %s", err)
	// }
	// // Print the response status, number of results, and request duration.
	// fmt.Printf(
	// 	"[%s] %d hits; took: %dms",
	// 	res.Status(),
	// 	int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
	// 	int(r["took"].(float64)),
	// )

	// for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	doc := hit.(map[string]interface{})
	// 	fields := doc["fields"]

	// 	fmt.Println(fields)
	// 	fmt.Println(reflect.TypeOf(fields))
	// 	pbi := fields.(map[string]interface{})["PRINT_BATCH_ID"]
	// 	fmt.Println(reflect.TypeOf(pbi))
	// 	//fmt.Println(reflect.TypeOf(reflect.ValueOf(pbi)))
	// 	fmt.Println(pbi)
	// 	pbi_s := make([]string, len(pbi.([]interface{})))
	// 	for i, v := range pbi.([]interface{}) {
	// 		pbi_s[i] = fmt.Sprint(v)
	// 	}
	// 	fmt.Println(pbi_s[0])

	// 	denom_en := fields.(map[string]interface{})["BN_DENOM_EN_NM"]
	// 	denom_en_s := make([]string, len(denom_en.([]interface{})))
	// 	for i, v := range denom_en.([]interface{}) {
	// 		denom_en_s[i] = fmt.Sprint(v)
	// 	}

	// 	t_exceed := fields.(map[string]interface{})["T_EXCEED_STRING_TXT"]
	// 	t_exceed_s := make([]string, len(t_exceed.([]interface{})))
	// 	for i, v := range t_exceed.([]interface{}) {
	// 		t_exceed_s[i] = fmt.Sprint(v)
	// 	}

	// 	year_nb := fields.(map[string]interface{})["YEAR_NB"]
	// 	year_nb_s := make([]string, len(year_nb.([]interface{})))
	// 	for i, v := range year_nb.([]interface{}) {
	// 		year_nb_s[i] = fmt.Sprint(v)
	// 	}

	// 	year_nb_i, err := strconv.Atoi(year_nb_s[0])
	// 	if err != nil {
	// 		fmt.Printf("%v", err)
			
	// 	}

	// 	qtr_en := fields.(map[string]interface{})["QUARTER_EN_NM"]
	// 	qtr_en_s := make([]string, len(qtr_en.([]interface{})))
	// 	for i, v := range qtr_en.([]interface{}) {
	// 		qtr_en_s[i] = fmt.Sprint(v)
	// 	}

	// 	month_en := fields.(map[string]interface{})["MONTH_EN_NM"]
	// 	month_en_s := make([]string, len(month_en.([]interface{})))
	// 	for i, v := range month_en.([]interface{}) {
	// 		month_en_s[i] = fmt.Sprint(v)
	// 	}

	// 	row_ctr := fields.(map[string]interface{})["ROW_COUNTER_NB"]
	// 	row_ctr_s := make([]string, len(row_ctr.([]interface{})))
	// 	for i, v := range row_ctr.([]interface{}) {
	// 		row_ctr_s[i] = fmt.Sprint(v)
	// 	}

	// 	row_ctr_i, err := strconv.Atoi(row_ctr_s[0])
	// 	if err != nil {
	// 		fmt.Printf("%v", err)
			
	// 	}  

	// 	fmt.Println("printed fields")
	// 	resp = append(resp, GetNotesDestructionDetails{pbi_s[0], denom_en_s[0],t_exceed_s[0],year_nb_i,qtr_en_s[0],month_en_s[0],row_ctr_i})
	// }

	return &resp
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

func GetVelocityData_Search_Request(sr string, dn string, stacker string, sn string, day string, wk string, mn string, yr string, startdate string, enddate string, from string, scroll_ID string) *GetVelocityDataAllResponse {
	var (
		
		scrollID string
		series string
		denomination string
		sno string
		mtr_closed_tears_sum string
		fhlg_foil_hlgraphic_effect_ic string
		o_orientation_ce string
		hl_hole_areas_sum string
		skew_note_skew_nb string
		rmed_max_miss_edge_region string
		foreign_marks_total_front string
		gcwf_graffiti_over_win_foil string
		riwf_max_ink_wear_front_region string
		rcrn_max_miss_corner_region string
		tacf_tactile_feature_ic string
		gtf_graffiti_on_front_sum string
		rotr_max_open_tear_region string
		len_note_length_nb string
		rtap_tape_region string
		wid_note_width_nb string
		rmtr_max_closed_tear_region string
		iwb_ink_wear_on_back_nb string
		med_miss_edge_areas_sum string
		rgtf_max_graffiti_front_region string
		fol_foil_area_miss_sum string
		gtb_graffiti_on_back_sum string
		slf_soil_on_front_sum string
		hilo_hi_low_note_ride_nb string
		otr_open_tear_lengths_sum string
		slb_soil_on_back_sum string
		retr_max_clsd_edge_tear_region string
		miwb_max_ink_wear_on_back string
		rfed_max_folded_edge_region string
		optically_variable_ink_presence string
		fld_fold_corner_areas_sum string
		owl_opacification_wear_level string
		swf_small_wind_feature_ic string
		mfol_max_foil_scratch_length string
		mfed_max_folded_edge_area string
		sfol_foil_scratch_lengths_sum string
		tape_tape_areas_sum string
		fed_folded_edge_areas_sum string
		foreign_marks_total_back string
		gwb_max_graffiti_on_back string
		mcrn_max_miss_corner_area string
		mmed_max_miss_edge_area string
		stnb_staining_discolor_back string
		mfld_max_folded_corner_area string
		area_note_area_rt string
		motr_max_open_tear_length string
		rhl_max_hole_region string
		creases_crumple_score string
		crn_miss_corner_areas_sum string
		riwb_max_ink_wear_back_region string
		mhl_max_hole_area string
		optically_variable_ink_score string
		// machine_ce string
		ocis_process_run_date_id string
		gwf_max_graffiti_on_front string
		mmtr_max_closed_tear string
		miwf_max_ink_wear_on_front string
		iwf_ink_wear_on_front_nb string
		rfld_max_fold_corner_region string
		etr_closed_edge_tears_sum string
		rgtb_max_graffiti_back_region string
		boc_loc_id string
		foil_fitness string
		metr_max_closed_edge_tear string
		stnf_staining_discolor_front string
		adi_process_run_date_id string
		output_stacker_en_nm string
		fi_ce string
		bps_shift_nb string
		adi string
		adi_wk int
		adi_mnth int
		adi_day int
		rdp_ce string
		deposit_nb string
		adi_yr int
		adi_qtr int
		machine_ce string 

	)

	resp := GetVelocityDataResponse{}
	resp_all := GetVelocityDataAllResponse{}
	i_from,err := strconv.Atoi(from)
	if err != nil {
		fmt.Printf("Error with from value: %s", err)
		}

	es,_ := es_connect()
	if scroll_ID == "None" {
		// Build the request body.
		var buf bytes.Buffer
		query := map[string]interface{}{
			"size": 5000,
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"wildcard": map[string]interface{}{
								"E_EMISSION_CE": sr,
							},
						},
						{
							"wildcard": map[string]interface{}{
								"DNM_BANK_NOTE_DENOM_AM": dn,
							},
						},
						{
							"wildcard": map[string]interface{}{
								"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.OUTPUT_STACKER_EN_NM": stacker,
							},
						},
						{
							"wildcard": map[string]interface{}{
								"SN_SERIAL_NUMBER": sn,
							},
						},
						{
							"range": map[string]interface{}{
								"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.ADI_PROCESS_RUN_DATE_ID": map[string]interface{}{
									"gte": startdate,
									"lte": enddate,
								},
							},
						},
						// {
						// 	"wildcard": map[string]interface{}{
						// 		"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.ADI_DAY": day,
						// 	},
						// },
						// {
						// 	"wildcard": map[string]interface{}{
						// 		"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.ADI_WK": wk,
						// 	},
						// },
						// {
						// 	"wildcard": map[string]interface{}{
						// 		"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.ADI_MNTH": mn,
						// 	},
						// },
						
						
					},
				},
			},"_source": map[string]interface{}{
				"includes": []string {"E_EMISSION_CE", "DNM_BANK_NOTE_DENOM_AM", "SN_SERIAL_NUMBER", "merged_SN_DETAILS.transformed_ADI_SN_DETAILS.*", "merged_SN_DETAILS.transformed_OCIS_SN_DETAILS.*"}, 
				"excludes": []string {"merged_SN_DETAILS.transformed_ADI_SN_DETAILS.*.keyword", "merged_SN_DETAILS.transformed_OCIS_SN_DETAILS.*.keyword"},
			  },
		}

		
		// fmt.Println("<----------------------------------------------->")
		// fmt.Println(query)
		// fmt.Println("<----------------------------------------------->")
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			fmt.Printf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("velocity_valid_all_" + sr + "_" + dn + "*"),
			es.Search.WithBody(&buf),
			es.Search.WithFrom(i_from),
			// es.Search.WithSize(500),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
			es.Search.WithScroll(time.Minute),
		)
		if err != nil {
		fmt.Printf("Error getting response: %s", err)
		}	

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				fmt.Printf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				fmt.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
				)
			}
			}
		// defer res.Body.Close()
		// Handle the first batch of data and extract the scrollID
		//
		json := read(res.Body)
		res.Body.Close()

		fmt.Println("json TYPE:", reflect.TypeOf(json), "\n")

		////scrollID = gjson.Get(json, "_scroll_id").String()
		total,_ := strconv.Atoi(gjson.Get(json, "hits.total.value").String())
		content := gjson.Get(json, "hits.hits.#._source")
		scroll_ID = gjson.Get(json, "_scroll_id").String()
		content.ForEach(func(key, value gjson.Result) bool {
			value.ForEach(func(key1, value1 gjson.Result) bool {
				// fmt.Println(key1.String())
					// if key1.String()  == "merged_SN_DETAILS.transformed_OCIS_SN_DETAILS.MTR_CLOSED_TEARS_SUM" {
					// 	mcts = value1.Array()[0].String()
					// }
					if key1.String()  == "E_EMISSION_CE" {
						series = value1.Array()[0].String()
					}
					if key1.String()  == "SN_SERIAL_NUMBER" {
						sno = value1.Array()[0].String()
					}
					if key1.String()  == "DNM_BANK_NOTE_DENOM_AM" {
						denomination = value1.Array()[0].String()
					}
					if key1.String()  == "merged_SN_DETAILS" {
						value1.ForEach(func(key2, value2 gjson.Result) bool {
							value2.ForEach(func(key3, value3 gjson.Result) bool {
								// fmt.Println(key3.String())
								if key3.String()  == "MTR_CLOSED_TEARS_SUM" { 
									mtr_closed_tears_sum = value3.String()
								}

								if key3.String()  == "FHLG_FOIL_HLGRAPHIC_EFFECT_IC" { 
									fhlg_foil_hlgraphic_effect_ic = value3.String() 
									}
								if key3.String()  == "O_ORIENTATION_CE" { 
									o_orientation_ce = value3.String() 
									}
								if key3.String()  == "HL_HOLE_AREAS_SUM" { 
									hl_hole_areas_sum = value3.String() 
									}
								if key3.String()  == "SKEW_NOTE_SKEW_NB" { 
									skew_note_skew_nb = value3.String() 
									}
								if key3.String()  == "RMED_MAX_MISS_EDGE_REGION" { 
									rmed_max_miss_edge_region = value3.String() 
									}
								if key3.String()  == "FOREIGN_MARKS_TOTAL_FRONT" { 
									foreign_marks_total_front = value3.String() 
									}
								if key3.String()  == "GCWF_GRAFFITI_OVER_WIN_FOIL" { 
									gcwf_graffiti_over_win_foil = value3.String() 
									}
								if key3.String()  == "RIWF_MAX_INK_WEAR_FRONT_REGION" { 
									riwf_max_ink_wear_front_region = value3.String() 
									}
								if key3.String()  == "RCRN_MAX_MISS_CORNER_REGION" { 
									rcrn_max_miss_corner_region = value3.String() 
									}
								if key3.String()  == "TACF_TACTILE_FEATURE_IC" { 
									tacf_tactile_feature_ic = value3.String() 
									}
								if key3.String()  == "GTF_GRAFFITI_ON_FRONT_SUM" { 
									gtf_graffiti_on_front_sum = value3.String() 
									}
								if key3.String()  == "ROTR_MAX_OPEN_TEAR_REGION" { 
									rotr_max_open_tear_region = value3.String() 
									}
								if key3.String()  == "LEN_NOTE_LENGTH_NB" { 
									len_note_length_nb = value3.String() 
									}
								if key3.String()  == "RTAP_TAPE_REGION" { 
									rtap_tape_region = value3.String() 
									}
								if key3.String()  == "WID_NOTE_WIDTH_NB" { 
									wid_note_width_nb = value3.String() 
									}
								if key3.String()  == "RMTR_MAX_CLOSED_TEAR_REGION" { 
									rmtr_max_closed_tear_region = value3.String() 
									}
								if key3.String()  == "IWB_INK_WEAR_ON_BACK_NB" { 
									iwb_ink_wear_on_back_nb = value3.String() 
									}
								if key3.String()  == "MED_MISS_EDGE_AREAS_SUM" { 
									med_miss_edge_areas_sum = value3.String() 
									}
								if key3.String()  == "RGTF_MAX_GRAFFITI_FRONT_REGION" { 
									rgtf_max_graffiti_front_region = value3.String() 
									}
								if key3.String()  == "FOL_FOIL_AREA_MISS_SUM" { 
									fol_foil_area_miss_sum = value3.String() 
									}
								if key3.String()  == "GTB_GRAFFITI_ON_BACK_SUM" { 
									gtb_graffiti_on_back_sum = value3.String() 
									}
								if key3.String()  == "SLF_SOIL_ON_FRONT_SUM" { 
									slf_soil_on_front_sum = value3.String() 
									}
								if key3.String()  == "HILO_HI_LOW_NOTE_RIDE_NB" { 
									hilo_hi_low_note_ride_nb = value3.String() 
									}
								if key3.String()  == "OTR_OPEN_TEAR_LENGTHS_SUM" { 
									otr_open_tear_lengths_sum = value3.String() 
									}
								if key3.String()  == "SLB_SOIL_ON_BACK_SUM" { 
									slb_soil_on_back_sum = value3.String() 
									}
								if key3.String()  == "RETR_MAX_CLSD_EDGE_TEAR_REGION" { 
									retr_max_clsd_edge_tear_region = value3.String() 
									}
								if key3.String()  == "MIWB_MAX_INK_WEAR_ON_BACK" { 
									miwb_max_ink_wear_on_back = value3.String() 
									}
								if key3.String()  == "RFED_MAX_FOLDED_EDGE_REGION" { 
									rfed_max_folded_edge_region = value3.String() 
									}
								if key3.String()  == "OPTICALLY_VARIABLE_INK_PRESENCE" { 
									optically_variable_ink_presence = value3.String() 
									}
								if key3.String()  == "FLD_FOLD_CORNER_AREAS_SUM" { 
									fld_fold_corner_areas_sum = value3.String() 
									}
								if key3.String()  == "OWL_OPACIFICATION_WEAR_LEVEL" { 
									owl_opacification_wear_level = value3.String() 
									}
								if key3.String()  == "SWF_SMALL_WIND_FEATURE_IC" { 
									swf_small_wind_feature_ic = value3.String() 
									}
								if key3.String()  == "MFOL_MAX_FOIL_SCRATCH_LENGTH" { 
									mfol_max_foil_scratch_length = value3.String() 
									}
								if key3.String()  == "MFED_MAX_FOLDED_EDGE_AREA" { 
									mfed_max_folded_edge_area = value3.String() 
									}
								if key3.String()  == "SFOL_FOIL_SCRATCH_LENGTHS_SUM" { 
									sfol_foil_scratch_lengths_sum = value3.String() 
									}
								if key3.String()  == "TAPE_TAPE_AREAS_SUM" { 
									tape_tape_areas_sum = value3.String() 
									}
								if key3.String()  == "FED_FOLDED_EDGE_AREAS_SUM" { 
									fed_folded_edge_areas_sum = value3.String() 
									}
								if key3.String()  == "FOREIGN_MARKS_TOTAL_BACK" { 
									foreign_marks_total_back = value3.String() 
									}
								if key3.String()  == "GWB_MAX_GRAFFITI_ON_BACK" { 
									gwb_max_graffiti_on_back = value3.String() 
									}
								if key3.String()  == "MCRN_MAX_MISS_CORNER_AREA" { 
									mcrn_max_miss_corner_area = value3.String() 
									}
								if key3.String()  == "MMED_MAX_MISS_EDGE_AREA" { 
									mmed_max_miss_edge_area = value3.String() 
									}
								if key3.String()  == "STNB_STAINING_DISCOLOR_BACK" { 
									stnb_staining_discolor_back = value3.String() 
									}
								if key3.String()  == "MFLD_MAX_FOLDED_CORNER_AREA" { 
									mfld_max_folded_corner_area = value3.String() 
									}
								if key3.String()  == "AREA_NOTE_AREA_RT" { 
									area_note_area_rt = value3.String() 
									}
								if key3.String()  == "MOTR_MAX_OPEN_TEAR_LENGTH" { 
									motr_max_open_tear_length = value3.String() 
									}
								if key3.String()  == "RHL_MAX_HOLE_REGION" { 
									rhl_max_hole_region = value3.String() 
									}
								if key3.String()  == "CREASES_CRUMPLE_SCORE" { 
									creases_crumple_score = value3.String() 
									}
								if key3.String()  == "CRN_MISS_CORNER_AREAS_SUM" { 
									crn_miss_corner_areas_sum = value3.String() 
									}
								if key3.String()  == "RIWB_MAX_INK_WEAR_BACK_REGION" { 
									riwb_max_ink_wear_back_region = value3.String() 
									}
								if key3.String()  == "MHL_MAX_HOLE_AREA" { 
									mhl_max_hole_area = value3.String() 
									}
								if key3.String()  == "OPTICALLY_VARIABLE_INK_SCORE" { 
									optically_variable_ink_score = value3.String() 
									}
								if key3.String()  == "MACHINE_CE" { 
									machine_ce = value3.String() 
									}
								if key3.String()  == "OCIS_PROCESS_RUN_DATE_ID" { 
									ocis_process_run_date_id = value3.String() 
									}
								if key3.String()  == "GWF_MAX_GRAFFITI_ON_FRONT" { 
									gwf_max_graffiti_on_front = value3.String() 
									}
								if key3.String()  == "MMTR_MAX_CLOSED_TEAR" { 
									mmtr_max_closed_tear = value3.String() 
									}
								if key3.String()  == "MIWF_MAX_INK_WEAR_ON_FRONT" { 
									miwf_max_ink_wear_on_front = value3.String() 
									}
								if key3.String()  == "IWF_INK_WEAR_ON_FRONT_NB" { 
									iwf_ink_wear_on_front_nb = value3.String() 
									}
								if key3.String()  == "RFLD_MAX_FOLD_CORNER_REGION" { 
									rfld_max_fold_corner_region = value3.String() 
									}
								if key3.String()  == "ETR_CLOSED_EDGE_TEARS_SUM" { 
									etr_closed_edge_tears_sum = value3.String() 
									}
								if key3.String()  == "RGTB_MAX_GRAFFITI_BACK_REGION" { 
									rgtb_max_graffiti_back_region = value3.String() 
									}
								if key3.String()  == "BOC_LOC_ID" { 
									boc_loc_id = value3.String() 
									}
								if key3.String()  == "FOIL_FITNESS" { 
									foil_fitness = value3.String() 
									}
								if key3.String()  == "METR_MAX_CLOSED_EDGE_TEAR" { 
									metr_max_closed_edge_tear = value3.String() 
									}
								if key3.String()  == "STNF_STAINING_DISCOLOR_FRONT" { 
									stnf_staining_discolor_front = value3.String() 
									}
								if key3.String()  == "ADI_PROCESS_RUN_DATE_ID" { 
									adi_process_run_date_id = value3.String() 
									}
								if key3.String()  == "OUTPUT_STACKER_EN_NM" { 
									output_stacker_en_nm = value3.String() 
									}
								if key3.String()  == "FI_CE" { 
									fi_ce = value3.String() 
									}
								if key3.String()  == "BPS_SHIFT_NB" { 
									bps_shift_nb = value3.String() 
									}
								if key3.String()  == "ADI" { 
									adi = value3.String() 
									}
								if key3.String()  == "ADI_WK" { 
									adi_wk,_ = strconv.Atoi(value3.String())
									}
								if key3.String()  == "ADI_MNTH" { 
									adi_mnth,_ = strconv.Atoi(value3.String())
									}
								if key3.String()  == "ADI_DAY" { 
									adi_day,_ = strconv.Atoi(value3.String())
									}
								if key3.String()  == "RDP_CE" { 
									rdp_ce = value3.String() 
									}
								if key3.String()  == "DEPOSIT_NB" { 
									deposit_nb = value3.String() 
									}
								if key3.String()  == "ADI_YR" { 
									adi_yr,_ = strconv.Atoi(value3.String())
									}
								if key3.String()  == "ADI_QTR" { 
									adi_qtr,_ = strconv.Atoi(value3.String())
									}
								if key3.String()  == "MACHINE_CE" { 
									machine_ce = value3.String() 
									}

								return true
							})
							return true
						})
					}

					return true
				})
			/* fmt.Println("series", series)
			fmt.Println("sno", sno)
			fmt.Println("denomination", denomination)
			fmt.Println("mtr_closed_tears_sum", mtr_closed_tears_sum)
			fmt.Println("machine_ce", machine_ce) */

			resp = append(resp, GetVelocityData{series, denomination, sno, area_note_area_rt, boc_loc_id, creases_crumple_score, crn_miss_corner_areas_sum, etr_closed_edge_tears_sum, fed_folded_edge_areas_sum, fhlg_foil_hlgraphic_effect_ic, fld_fold_corner_areas_sum, foil_fitness, fol_foil_area_miss_sum, foreign_marks_total_back, foreign_marks_total_front, gcwf_graffiti_over_win_foil, gtb_graffiti_on_back_sum, gtf_graffiti_on_front_sum, gwb_max_graffiti_on_back, gwf_max_graffiti_on_front, hilo_hi_low_note_ride_nb, hl_hole_areas_sum, iwb_ink_wear_on_back_nb, iwf_ink_wear_on_front_nb, len_note_length_nb, machine_ce, mcrn_max_miss_corner_area, med_miss_edge_areas_sum, metr_max_closed_edge_tear, mfed_max_folded_edge_area, mfld_max_folded_corner_area, mfol_max_foil_scratch_length, mhl_max_hole_area, miwb_max_ink_wear_on_back, miwf_max_ink_wear_on_front, mmed_max_miss_edge_area, mmtr_max_closed_tear, motr_max_open_tear_length, mtr_closed_tears_sum, optically_variable_ink_presence, optically_variable_ink_score, otr_open_tear_lengths_sum, owl_opacification_wear_level, o_orientation_ce, rcrn_max_miss_corner_region, retr_max_clsd_edge_tear_region, rfed_max_folded_edge_region, rfld_max_fold_corner_region, rgtb_max_graffiti_back_region, rgtf_max_graffiti_front_region, rhl_max_hole_region, riwb_max_ink_wear_back_region, riwf_max_ink_wear_front_region, rmed_max_miss_edge_region, rmtr_max_closed_tear_region, rotr_max_open_tear_region, rtap_tape_region, sfol_foil_scratch_lengths_sum, skew_note_skew_nb, slb_soil_on_back_sum, slf_soil_on_front_sum, stnb_staining_discolor_back, stnf_staining_discolor_front, swf_small_wind_feature_ic, tacf_tactile_feature_ic, tape_tape_areas_sum, wid_note_width_nb, adi, adi_day, adi_mnth, adi_process_run_date_id, adi_qtr, adi_wk, adi_yr, bps_shift_nb, deposit_nb, fi_ce, output_stacker_en_nm, rdp_ce})
			// resp = append(resp, GetVelocityData{series, denomination, sno})
			return true
		})
		
		
		// fmt.Println("Content", content)
		fmt.Println("Total", total)
		fmt.Println("Scroll_id", scroll_ID)

		resp_all = GetVelocityDataAllResponse{scroll_ID, total, resp}
		
		// fmt.Println("resp_all -->", &resp_all)
			
	} else {
		scrollID = scroll_ID

		// Perform the scroll request and pass the scrollID and scroll duration
		res_scroll, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
		if res_scroll.IsError() {
			fmt.Printf("Error response: %s", res_scroll)
		}

		json_scroll := read(res_scroll.Body)
		res_scroll.Body.Close()

		total,_ := strconv.Atoi(gjson.Get(json_scroll, "hits.total.value").String())
		content := gjson.Get(json_scroll, "hits.hits.#._source")
		// scroll_ID = gjson.Get(json, "_scroll_id").String()
		hits := gjson.Get(json_scroll, "hits.hits")

		if len(hits.Array()) < 1 {
			fmt.Println("Finished scrolling")
			// return &GetVelocityDataAllResponse{scroll_ID, total, resp}
		} else {

					content.ForEach(func(key, value gjson.Result) bool {
						value.ForEach(func(key1, value1 gjson.Result) bool {
							// fmt.Println(key1.String())
							if key1.String()  == "E_EMISSION_CE" {
								series = value1.Array()[0].String()
							}
							if key1.String()  == "SN_SERIAL_NUMBER" {
								sno = value1.Array()[0].String()
							}
							if key1.String()  == "DNM_BANK_NOTE_DENOM_AM" {
								denomination = value1.Array()[0].String()
							}

							if key1.String()  == "merged_SN_DETAILS" {
								value1.ForEach(func(key2, value2 gjson.Result) bool {
									value2.ForEach(func(key3, value3 gjson.Result) bool {
										// fmt.Println(key3.String())
										if key3.String()  == "MTR_CLOSED_TEARS_SUM" { 
											mtr_closed_tears_sum = value3.String()
										}
		
										if key3.String()  == "FHLG_FOIL_HLGRAPHIC_EFFECT_IC" { 
											fhlg_foil_hlgraphic_effect_ic = value3.String() 
											}
										if key3.String()  == "O_ORIENTATION_CE" { 
											o_orientation_ce = value3.String() 
											}
										if key3.String()  == "HL_HOLE_AREAS_SUM" { 
											hl_hole_areas_sum = value3.String() 
											}
										if key3.String()  == "SKEW_NOTE_SKEW_NB" { 
											skew_note_skew_nb = value3.String() 
											}
										if key3.String()  == "RMED_MAX_MISS_EDGE_REGION" { 
											rmed_max_miss_edge_region = value3.String() 
											}
										if key3.String()  == "FOREIGN_MARKS_TOTAL_FRONT" { 
											foreign_marks_total_front = value3.String() 
											}
										if key3.String()  == "GCWF_GRAFFITI_OVER_WIN_FOIL" { 
											gcwf_graffiti_over_win_foil = value3.String() 
											}
										if key3.String()  == "RIWF_MAX_INK_WEAR_FRONT_REGION" { 
											riwf_max_ink_wear_front_region = value3.String() 
											}
										if key3.String()  == "RCRN_MAX_MISS_CORNER_REGION" { 
											rcrn_max_miss_corner_region = value3.String() 
											}
										if key3.String()  == "TACF_TACTILE_FEATURE_IC" { 
											tacf_tactile_feature_ic = value3.String() 
											}
										if key3.String()  == "GTF_GRAFFITI_ON_FRONT_SUM" { 
											gtf_graffiti_on_front_sum = value3.String() 
											}
										if key3.String()  == "ROTR_MAX_OPEN_TEAR_REGION" { 
											rotr_max_open_tear_region = value3.String() 
											}
										if key3.String()  == "LEN_NOTE_LENGTH_NB" { 
											len_note_length_nb = value3.String() 
											}
										if key3.String()  == "RTAP_TAPE_REGION" { 
											rtap_tape_region = value3.String() 
											}
										if key3.String()  == "WID_NOTE_WIDTH_NB" { 
											wid_note_width_nb = value3.String() 
											}
										if key3.String()  == "RMTR_MAX_CLOSED_TEAR_REGION" { 
											rmtr_max_closed_tear_region = value3.String() 
											}
										if key3.String()  == "IWB_INK_WEAR_ON_BACK_NB" { 
											iwb_ink_wear_on_back_nb = value3.String() 
											}
										if key3.String()  == "MED_MISS_EDGE_AREAS_SUM" { 
											med_miss_edge_areas_sum = value3.String() 
											}
										if key3.String()  == "RGTF_MAX_GRAFFITI_FRONT_REGION" { 
											rgtf_max_graffiti_front_region = value3.String() 
											}
										if key3.String()  == "FOL_FOIL_AREA_MISS_SUM" { 
											fol_foil_area_miss_sum = value3.String() 
											}
										if key3.String()  == "GTB_GRAFFITI_ON_BACK_SUM" { 
											gtb_graffiti_on_back_sum = value3.String() 
											}
										if key3.String()  == "SLF_SOIL_ON_FRONT_SUM" { 
											slf_soil_on_front_sum = value3.String() 
											}
										if key3.String()  == "HILO_HI_LOW_NOTE_RIDE_NB" { 
											hilo_hi_low_note_ride_nb = value3.String() 
											}
										if key3.String()  == "OTR_OPEN_TEAR_LENGTHS_SUM" { 
											otr_open_tear_lengths_sum = value3.String() 
											}
										if key3.String()  == "SLB_SOIL_ON_BACK_SUM" { 
											slb_soil_on_back_sum = value3.String() 
											}
										if key3.String()  == "RETR_MAX_CLSD_EDGE_TEAR_REGION" { 
											retr_max_clsd_edge_tear_region = value3.String() 
											}
										if key3.String()  == "MIWB_MAX_INK_WEAR_ON_BACK" { 
											miwb_max_ink_wear_on_back = value3.String() 
											}
										if key3.String()  == "RFED_MAX_FOLDED_EDGE_REGION" { 
											rfed_max_folded_edge_region = value3.String() 
											}
										if key3.String()  == "OPTICALLY_VARIABLE_INK_PRESENCE" { 
											optically_variable_ink_presence = value3.String() 
											}
										if key3.String()  == "FLD_FOLD_CORNER_AREAS_SUM" { 
											fld_fold_corner_areas_sum = value3.String() 
											}
										if key3.String()  == "OWL_OPACIFICATION_WEAR_LEVEL" { 
											owl_opacification_wear_level = value3.String() 
											}
										if key3.String()  == "SWF_SMALL_WIND_FEATURE_IC" { 
											swf_small_wind_feature_ic = value3.String() 
											}
										if key3.String()  == "MFOL_MAX_FOIL_SCRATCH_LENGTH" { 
											mfol_max_foil_scratch_length = value3.String() 
											}
										if key3.String()  == "MFED_MAX_FOLDED_EDGE_AREA" { 
											mfed_max_folded_edge_area = value3.String() 
											}
										if key3.String()  == "SFOL_FOIL_SCRATCH_LENGTHS_SUM" { 
											sfol_foil_scratch_lengths_sum = value3.String() 
											}
										if key3.String()  == "TAPE_TAPE_AREAS_SUM" { 
											tape_tape_areas_sum = value3.String() 
											}
										if key3.String()  == "FED_FOLDED_EDGE_AREAS_SUM" { 
											fed_folded_edge_areas_sum = value3.String() 
											}
										if key3.String()  == "FOREIGN_MARKS_TOTAL_BACK" { 
											foreign_marks_total_back = value3.String() 
											}
										if key3.String()  == "GWB_MAX_GRAFFITI_ON_BACK" { 
											gwb_max_graffiti_on_back = value3.String() 
											}
										if key3.String()  == "MCRN_MAX_MISS_CORNER_AREA" { 
											mcrn_max_miss_corner_area = value3.String() 
											}
										if key3.String()  == "MMED_MAX_MISS_EDGE_AREA" { 
											mmed_max_miss_edge_area = value3.String() 
											}
										if key3.String()  == "STNB_STAINING_DISCOLOR_BACK" { 
											stnb_staining_discolor_back = value3.String() 
											}
										if key3.String()  == "MFLD_MAX_FOLDED_CORNER_AREA" { 
											mfld_max_folded_corner_area = value3.String() 
											}
										if key3.String()  == "AREA_NOTE_AREA_RT" { 
											area_note_area_rt = value3.String() 
											}
										if key3.String()  == "MOTR_MAX_OPEN_TEAR_LENGTH" { 
											motr_max_open_tear_length = value3.String() 
											}
										if key3.String()  == "RHL_MAX_HOLE_REGION" { 
											rhl_max_hole_region = value3.String() 
											}
										if key3.String()  == "CREASES_CRUMPLE_SCORE" { 
											creases_crumple_score = value3.String() 
											}
										if key3.String()  == "CRN_MISS_CORNER_AREAS_SUM" { 
											crn_miss_corner_areas_sum = value3.String() 
											}
										if key3.String()  == "RIWB_MAX_INK_WEAR_BACK_REGION" { 
											riwb_max_ink_wear_back_region = value3.String() 
											}
										if key3.String()  == "MHL_MAX_HOLE_AREA" { 
											mhl_max_hole_area = value3.String() 
											}
										if key3.String()  == "OPTICALLY_VARIABLE_INK_SCORE" { 
											optically_variable_ink_score = value3.String() 
											}
										if key3.String()  == "MACHINE_CE" { 
											machine_ce = value3.String() 
											}
										if key3.String()  == "OCIS_PROCESS_RUN_DATE_ID" { 
											ocis_process_run_date_id = value3.String() 
											}
										if key3.String()  == "GWF_MAX_GRAFFITI_ON_FRONT" { 
											gwf_max_graffiti_on_front = value3.String() 
											}
										if key3.String()  == "MMTR_MAX_CLOSED_TEAR" { 
											mmtr_max_closed_tear = value3.String() 
											}
										if key3.String()  == "MIWF_MAX_INK_WEAR_ON_FRONT" { 
											miwf_max_ink_wear_on_front = value3.String() 
											}
										if key3.String()  == "IWF_INK_WEAR_ON_FRONT_NB" { 
											iwf_ink_wear_on_front_nb = value3.String() 
											}
										if key3.String()  == "RFLD_MAX_FOLD_CORNER_REGION" { 
											rfld_max_fold_corner_region = value3.String() 
											}
										if key3.String()  == "ETR_CLOSED_EDGE_TEARS_SUM" { 
											etr_closed_edge_tears_sum = value3.String() 
											}
										if key3.String()  == "RGTB_MAX_GRAFFITI_BACK_REGION" { 
											rgtb_max_graffiti_back_region = value3.String() 
											}
										if key3.String()  == "BOC_LOC_ID" { 
											boc_loc_id = value3.String() 
											}
										if key3.String()  == "FOIL_FITNESS" { 
											foil_fitness = value3.String() 
											}
										if key3.String()  == "METR_MAX_CLOSED_EDGE_TEAR" { 
											metr_max_closed_edge_tear = value3.String() 
											}
										if key3.String()  == "STNF_STAINING_DISCOLOR_FRONT" { 
											stnf_staining_discolor_front = value3.String() 
											}
										if key3.String()  == "ADI_PROCESS_RUN_DATE_ID" { 
											adi_process_run_date_id = value3.String() 
											}
										if key3.String()  == "OUTPUT_STACKER_EN_NM" { 
											output_stacker_en_nm = value3.String() 
											}
										if key3.String()  == "FI_CE" { 
											fi_ce = value3.String() 
											}
										if key3.String()  == "BPS_SHIFT_NB" { 
											bps_shift_nb = value3.String() 
											}
										if key3.String()  == "ADI" { 
											adi = value3.String() 
											}
										if key3.String()  == "ADI_WK" { 
											adi_wk,_ = strconv.Atoi(value3.String())
											}
										if key3.String()  == "ADI_MNTH" { 
											adi_mnth,_ = strconv.Atoi(value3.String())
											}
										if key3.String()  == "ADI_DAY" { 
											adi_day,_ = strconv.Atoi(value3.String())
											}
										if key3.String()  == "RDP_CE" { 
											rdp_ce = value3.String() 
											}
										if key3.String()  == "DEPOSIT_NB" { 
											deposit_nb = value3.String() 
											}
										if key3.String()  == "ADI_YR" { 
											adi_yr,_ = strconv.Atoi(value3.String())
											}
										if key3.String()  == "ADI_QTR" { 
											adi_qtr,_ = strconv.Atoi(value3.String())
											}
										if key3.String()  == "MACHINE_CE" { 
											machine_ce = value3.String() 
											}
										return true
									})
									return true
								})
							}
							
							return true
						})
						resp = append(resp, GetVelocityData{series, denomination, sno, area_note_area_rt, boc_loc_id, creases_crumple_score, crn_miss_corner_areas_sum, etr_closed_edge_tears_sum, fed_folded_edge_areas_sum, fhlg_foil_hlgraphic_effect_ic, fld_fold_corner_areas_sum, foil_fitness, fol_foil_area_miss_sum, foreign_marks_total_back, foreign_marks_total_front, gcwf_graffiti_over_win_foil, gtb_graffiti_on_back_sum, gtf_graffiti_on_front_sum, gwb_max_graffiti_on_back, gwf_max_graffiti_on_front, hilo_hi_low_note_ride_nb, hl_hole_areas_sum, iwb_ink_wear_on_back_nb, iwf_ink_wear_on_front_nb, len_note_length_nb, machine_ce, mcrn_max_miss_corner_area, med_miss_edge_areas_sum, metr_max_closed_edge_tear, mfed_max_folded_edge_area, mfld_max_folded_corner_area, mfol_max_foil_scratch_length, mhl_max_hole_area, miwb_max_ink_wear_on_back, miwf_max_ink_wear_on_front, mmed_max_miss_edge_area, mmtr_max_closed_tear, motr_max_open_tear_length, mtr_closed_tears_sum, optically_variable_ink_presence, optically_variable_ink_score, otr_open_tear_lengths_sum, owl_opacification_wear_level, o_orientation_ce, rcrn_max_miss_corner_region, retr_max_clsd_edge_tear_region, rfed_max_folded_edge_region, rfld_max_fold_corner_region, rgtb_max_graffiti_back_region, rgtf_max_graffiti_front_region, rhl_max_hole_region, riwb_max_ink_wear_back_region, riwf_max_ink_wear_front_region, rmed_max_miss_edge_region, rmtr_max_closed_tear_region, rotr_max_open_tear_region, rtap_tape_region, sfol_foil_scratch_lengths_sum, skew_note_skew_nb, slb_soil_on_back_sum, slf_soil_on_front_sum, stnb_staining_discolor_back, stnf_staining_discolor_front, swf_small_wind_feature_ic, tacf_tactile_feature_ic, tape_tape_areas_sum, wid_note_width_nb, adi, adi_day, adi_mnth, adi_process_run_date_id, adi_qtr, adi_wk, adi_yr, bps_shift_nb, deposit_nb, fi_ce, output_stacker_en_nm, rdp_ce})
						// resp = append(resp, GetVelocityData{series, denomination, sno})
						return true
					})
					fmt.Println("Total", total)
					fmt.Println("Scroll_id", scroll_ID)
					resp_all = GetVelocityDataAllResponse{scrollID, total, resp}
				}
			
		}
	// return &resp
	
	return &resp_all
}


func GetVelocityAggWearCategoryData_Search_Request() *GetVelocityAggWearCategoryDataResponse {
	var (
		series int
		year int 
		qtr int 
		denom int 
		fld_notes_with_zero_count float64 
		fld_average_excluding_zero float64 
		fld_non_zero_pct float64 
		fld_gt_eq_threshold_pct float64 
		fld_average_reading float64 
		fed_notes_with_zero_count float64 
		fed_average_excluding_zero float64 
		fed_non_zero_pct float64 
		fed_gt_eq_threshold_pct float64 
		fed_average_reading float64 
		crs_notes_with_zero_count float64 
		crs_excluding_zero_mean float64 
		crs_non_zero_pct float64 
		crs_gt_eq_threshold_pct float64 
		crs_average_reading float64 
		hl_notes_with_zero_count float64 
		hl_excluding_zero_mean float64 
		hl_non_zero_pct float64 
		hl_gt_eq_threshold_pct float64 
		hl_average_reading float64 
		otr_notes_with_zero_count float64 
		otr_excluding_zero_mean float64 
		otr_non_zero_pct float64 
		otr_gt_eq_threshold_pct float64 
		otr_average_reading float64 
		motr_notes_with_zero_count float64 
		motr_excluding_zero_mean float64 
		motr_non_zero_pct float64 
		motr_gt_eq_threshold_pct float64 
		motr_average_reading float64 
		etr_mtr_notes_with_zero_count float64 
		etr_mtr_excluding_zero_mean float64 
		etr_mtr_non_zero_pct float64 
		etr_mtr_gt_eq_threshold_pct float64 
		etr_mtr_average_reading float64 
		metr_mmtr_notes_with_zero_count float64 
		metr_mmtr_excluding_zero_mean float64 
		metr_mmtr_non_zero_pct float64 
		metr_mmtr_gt_eq_threshold_pct float64 
		metr_mmtr_average_reading float64 
		crn_notes_with_zero_count float64 
		crn_excluding_zero_mean float64 
		crn_non_zero_pct float64 
		crn_gt_eq_threshold_pct float64 
		crn_average_reading float64 
		med_notes_with_zero_count float64 
		med_excluding_zero_mean float64 
		med_non_zero_pct float64 
		med_gt_eq_threshold_pct float64 
		med_average_reading float64 
		tape_notes_with_zero_count float64 
		tape_excluding_zero_mean float64 
		tape_non_zero_pct float64 
		tape_gt_eq_threshold_pct float64 
		tape_average_reading float64 
		fol_notes_with_zero_count float64 
		fol_excluding_zero_mean float64 
		fol_non_zero_pct float64 
		fol_gt_eq_threshold_pct float64 
		fol_average_reading float64 
		sfol_notes_with_zero_count float64 
		sfol_excluding_zero_mean float64 
		sfol_non_zero_pct float64 
		sfol_gt_eq_threshold_pct float64 
		sfol_average_reading float64 
		mfol_notes_with_zero_count float64 
		mfol_excluding_zero_mean float64 
		mfol_non_zero_pct float64 
		mfol_gt_eq_threshold_pct float64 
		mfol_average_reading float64 
		iwb_notes_with_zero_count float64 
		iwb_excluding_zero_mean float64 
		iwb_non_zero_pct float64 
		iwb_gt_eq_threshold_pct float64 
		iwb_average_reading float64 
		iwf_notes_with_zero_count float64 
		iwf_excluding_zero_mean float64 
		iwf_non_zero_pct float64 
		iwf_gt_eq_threshold_pct float64 
		iwf_average_reading float64 
		gtf_notes_with_zero_count float64 
		gtf_excluding_zero_mean float64 
		gtf_non_zero_pct float64 
		gtf_gt_eq_threshold_pct float64 
		gtf_average_reading float64 
		gtb_notes_with_zero_count float64 
		gtb_excluding_zero_mean float64 
		gtb_non_zero_pct float64 
		gtb_gt_eq_threshold_pct float64 
		gtb_average_reading float64 
		gcwf_notes_with_zero_count float64 
		gcwf_excluding_zero_mean float64 
		gcwf_non_zero_pct float64 
		gcwf_gt_eq_threshold_pct float64 
		gcwf_average_reading float64 
		owl_notes_with_zero_count float64 
		owl_excluding_zero_mean float64 
		owl_non_zero_pct float64 
		owl_gt_eq_threshold_pct float64 
		owl_average_reading float64 
		stnf_notes_with_zero_count float64 
		stnf_excluding_zero_mean float64 
		stnf_non_zero_pct float64 
		stnf_gt_eq_threshold_pct float64 
		stnf_average_reading float64 
		stnb_notes_with_zero_count float64 
		stnb_excluding_zero_mean float64 
		stnb_non_zero_pct float64 
		stnb_gt_eq_threshold_pct float64 
		stnb_average_reading float64 
		fmtf_notes_with_zero_count float64 
		fmtf_excluding_zero_mean float64 
		fmtf_non_zero_pct float64 
		fmtf_gt_eq_threshold_pct float64 
		fmtf_average_reading float64 
		fmtb_notes_with_zero_count float64 
		fmtb_excluding_zero_mean float64 
		fmtb_non_zero_pct float64 
		fmtb_gt_eq_threshold_pct float64 
		fmtb_average_reading float64 
		ffit_notes_with_zero_count float64 
		ffit_excluding_zero_mean float64 
		ffit_non_zero_pct float64 
		ffit_gt_eq_threshold_pct float64 
		ffit_average_reading float64 
		ovip_notes_with_zero_count float64 
		ovip_excluding_zero_mean float64 
		ovip_non_zero_pct float64 
		ovip_gt_eq_threshold_pct float64 
		ovip_average_reading float64 
		ovis_notes_with_zero_count float64 
		ovis_excluding_zero_mean float64 
		ovis_non_zero_pct float64 
		ovis_gt_eq_threshold_pct float64 
		ovis_average_reading float64
	)
	
	resp_payload := GetVelocityAggWearCatDataResponse{}
	resp := GetVelocityAggWearCategoryDataResponse{}

	es,_ := es_connect()

	// Build the request body.
	var buf bytes.Buffer
	//query := map[string]interface{}{"query": map[string]interface{}{"match_all": map[string]interface{}{}}}

	query := map[string]interface{}{
		"sort": []map[string]interface{}{
				{
				"Series": map[string]interface{}{
					"order": "asc",
					},
				  "Year": map[string]interface{}{
					"order": "asc",
				  },
				 "Qtr": map[string]interface{}{
					"order": "asc",
				  },
				 "Denom": map[string]interface{}{
					"order": "asc",
				  }, 
				},
			  }}

   if err := json.NewEncoder(&buf).Encode(query); err != nil {
	fmt.Printf("Error encoding query: %s", err)
	  }
	  
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("velocity_valid_aggregate_all_properties_metric"),
		es.Search.WithBody(&buf),
		es.Search.WithSize(500),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	  )
	if err != nil {
	fmt.Printf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
	} else {
		// Print the response status and error information.
		fmt.Printf("[%s] %s: %s",
		res.Status(),
		e["error"].(map[string]interface{})["type"],
		e["error"].(map[string]interface{})["reason"],
		)
	}
	}

	json := read(res.Body)
	res.Body.Close()

	total,_ := strconv.Atoi(gjson.Get(json, "hits.total.value").String())
	content := gjson.Get(json, "hits.hits.#._source")

	// fmt.Println("Total = %d", total)
	// fmt.Println("Content = %s", content)

	content.ForEach(func(key, value gjson.Result) bool {
		value.ForEach(func(key1, value1 gjson.Result) bool {
			// fmt.Println(key1)
			// fmt.Println(value1)
			if key1.String()  == "Series" {
				series,_ = strconv.Atoi(value1.String())
				} 
			if key1.String()  == "Year" {
				year,_ = strconv.Atoi(value1.String())
				} 
				if key1.String()  == "Qtr" {
				qtr,_ = strconv.Atoi(value1.String())
				} 
				if key1.String()  == "Denom" {
				denom,_ = strconv.Atoi(value1.String())
				} 
				if key1.String()  == "FLD_Notes_with_Zero_Count" {
				fld_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FLD_Average_Excluding_Zero" {
				fld_average_excluding_zero,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FLD_Non_Zero_Pct" {
				fld_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FLD_gt_eq_Threshold_Pct" {
				fld_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FLD_Average_Reading" {
				fld_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FED_Notes_with_Zero_Count" {
				fed_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FED_Average_Excluding_Zero" {
				fed_average_excluding_zero,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FED_Non_Zero_Pct" {
				fed_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FED_gt_eq_Threshold_Pct" {
				fed_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "" {
				fed_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRS_Notes_with_Zero_Count" {
				crs_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRS_Excluding_Zero_Mean" {
				crs_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRS_Non_Zero_Pct" {
				crs_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRS_gt_eq_Threshold_Pct" {
				crs_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRS_Average_Reading" {
				crs_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "HL_Notes_with_Zero_Count" {
				hl_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "HL_Excluding_Zero_Mean" {
				hl_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "HL_Non_Zero_Pct" {
				hl_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "HL_gt_eq_Threshold_Pct" {
				hl_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "HL_Average_Reading" {
				hl_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OTR_Notes_with_Zero_Count" {
				otr_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OTR_Excluding_Zero_Mean" {
				otr_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OTR_Non_Zero_Pct" {
				otr_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OTR_gt_eq_Threshold_Pct" {
				otr_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OTR_Average_Reading" {
				otr_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MOTR_Notes_with_Zero_Count" {
				motr_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MOTR_Excluding_Zero_Mean" {
				motr_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MOTR_Non_Zero_Pct" {
				motr_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MOTR_gt_eq_Threshold_Pct" {
				motr_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MOTR_Average_Reading" {
				motr_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "ETR_MTR_Notes_with_Zero_Count" {
				etr_mtr_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "ETR_MTR_Excluding_Zero_Mean" {
				etr_mtr_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "ETR_MTR_Non_Zero_Pct" {
				etr_mtr_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "ETR_MTR_gt_eq_Threshold_Pct" {
				etr_mtr_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "ETR_MTR_Average_Reading" {
				etr_mtr_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "METR_MMTR_Notes_with_Zero_Count" {
				metr_mmtr_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "METR_MMTR_Excluding_Zero_Mean" {
				metr_mmtr_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "METR_MMTR_Non_Zero_Pct" {
				metr_mmtr_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "METR_MMTR_gt_eq_Threshold_Pct" {
				metr_mmtr_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "METR_MMTR_Average_Reading" {
				metr_mmtr_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRN_Notes_with_Zero_Count" {
				crn_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRN_Excluding_Zero_Mean" {
				crn_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRN_Non_Zero_Pct" {
				crn_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRN_gt_eq_Threshold_Pct" {
				crn_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "CRN_Average_Reading" {
				crn_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MED_Notes_with_Zero_Count" {
				med_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MED_Excluding_Zero_Mean" {
				med_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MED_Non_Zero_Pct" {
				med_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MED_gt_eq_Threshold_Pct" {
				med_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MED_Average_Reading" {
				med_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "TAPE_Notes_with_Zero_Count" {
				tape_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "TAPE_Excluding_Zero_Mean" {
				tape_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "TAPE_Non_Zero_Pct" {
				tape_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "TAPE_gt_eq_Threshold_Pct" {
				tape_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "TAPE_Average_Reading" {
				tape_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FOL_Notes_with_Zero_Count" {
				fol_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FOL_Excluding_Zero_Mean" {
				fol_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FOL_Non_Zero_Pct" {
				fol_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FOL_gt_eq_Threshold_Pct" {
				fol_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FOL_Average_Reading" {
				fol_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "SFOL_Notes_with_Zero_Count" {
				sfol_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "SFOL_Excluding_Zero_Mean" {
				sfol_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "SFOL_Non_Zero_Pct" {
				sfol_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "SFOL_gt_eq_Threshold_Pct" {
				sfol_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "SFOL_Average_Reading" {
				sfol_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MFOL_Notes_with_Zero_Count" {
				mfol_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MFOL_Excluding_Zero_Mean" {
				mfol_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MFOL_Non_Zero_Pct" {
				mfol_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MFOL_gt_eq_Threshold_Pct" {
				mfol_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "MFOL_Average_Reading" {
				mfol_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWB_Notes_with_Zero_Count" {
				iwb_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWB_Excluding_Zero_Mean" {
				iwb_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWB_Non_Zero_Pct" {
				iwb_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWB_gt_eq_Threshold_Pct" {
				iwb_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWB_Average_Reading" {
				iwb_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWF_Notes_with_Zero_Count" {
				iwf_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWF_Excluding_Zero_Mean" {
				iwf_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWF_Non_Zero_Pct" {
				iwf_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWF_gt_eq_Threshold_Pct" {
				iwf_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "IWF_Average_Reading" {
				iwf_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTF_Notes_with_Zero_Count" {
				gtf_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTF_Excluding_Zero_Mean" {
				gtf_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTF_Non_Zero_Pct" {
				gtf_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTF_gt_eq_Threshold_Pct" {
				gtf_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTF_Average_Reading" {
				gtf_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTB_Notes_with_Zero_Count" {
				gtb_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTB_Excluding_Zero_Mean" {
				gtb_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTB_Non_Zero_Pct" {
				gtb_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTB_gt_eq_Threshold_Pct" {
				gtb_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GTB_Average_Reading" {
				gtb_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GCWF_Notes_with_Zero_Count" {
				gcwf_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GCWF_Excluding_Zero_Mean" {
				gcwf_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GCWF_Non_Zero_Pct" {
				gcwf_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GCWF_gt_eq_Threshold_Pct" {
				gcwf_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "GCWF_Average_Reading" {
				gcwf_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OWL_Notes_with_Zero_Count" {
				owl_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OWL_Excluding_Zero_Mean" {
				owl_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OWL_Non_Zero_Pct" {
				owl_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OWL_gt_eq_Threshold_Pct" {
				owl_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OWL_Average_Reading" {
				owl_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNF_Notes_with_Zero_Count" {
				stnf_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNF_Excluding_Zero_Mean" {
				stnf_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNF_Non_Zero_Pct" {
				stnf_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNF_gt_eq_Threshold_Pct" {
				stnf_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNF_Average_Reading" {
				stnf_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNB_Notes_with_Zero_Count" {
				stnb_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNB_Excluding_Zero_Mean" {
				stnb_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNB_Non_Zero_Pct" {
				stnb_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNB_gt_eq_Threshold_Pct" {
				stnb_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "STNB_Average_Reading" {
				stnb_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTF_Notes_with_Zero_Count" {
				fmtf_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTF_Excluding_Zero_Mean" {
				fmtf_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTF_Non_Zero_Pct" {
				fmtf_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTF_gt_eq_Threshold_Pct" {
				fmtf_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTF_Average_Reading" {
				fmtf_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTB_Notes_with_Zero_Count" {
				fmtb_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTB_Excluding_Zero_Mean" {
				fmtb_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTB_Non_Zero_Pct" {
				fmtb_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTB_gt_eq_Threshold_Pct" {
				fmtb_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FMTB_Average_Reading" {
				fmtb_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FFIT_Notes_with_Zero_Count" {
				ffit_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FFIT_Excluding_Zero_Mean" {
				ffit_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FFIT_Non_Zero_Pct" {
				ffit_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FFIT_gt_eq_Threshold_Pct" {
				ffit_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "FFIT_Average_Reading" {
				ffit_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIP_Notes_with_Zero_Count" {
				ovip_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIP_Excluding_Zero_Mean" {
				ovip_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIP_Non_Zero_Pct" {
				ovip_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIP_gt_eq_Threshold_Pct" {
				ovip_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIP_Average_Reading" {
				ovip_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIS_Notes_with_Zero_Count" {
				ovis_notes_with_zero_count,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIS_Excluding_Zero_Mean" {
				ovis_excluding_zero_mean,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIS_Non_Zero_Pct" {
				ovis_non_zero_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIS_gt_eq_Threshold_Pct" {
				ovis_gt_eq_threshold_pct,_ = strconv.ParseFloat(value1.String(), 64)
				} 
				if key1.String()  == "OVIS_Average_Reading" {
				ovis_average_reading,_ = strconv.ParseFloat(value1.String(), 64)
				}
		
		return true
		})
		resp_payload = append(resp_payload, GetVelocityAggWearCatData{series, year, qtr, denom, fld_notes_with_zero_count, fld_average_excluding_zero, fld_non_zero_pct, fld_gt_eq_threshold_pct, fld_average_reading, fed_notes_with_zero_count, fed_average_excluding_zero, fed_non_zero_pct, fed_gt_eq_threshold_pct, fed_average_reading, crs_notes_with_zero_count, crs_excluding_zero_mean, crs_non_zero_pct, crs_gt_eq_threshold_pct, crs_average_reading, hl_notes_with_zero_count, hl_excluding_zero_mean, hl_non_zero_pct, hl_gt_eq_threshold_pct, hl_average_reading, otr_notes_with_zero_count, otr_excluding_zero_mean, otr_non_zero_pct, otr_gt_eq_threshold_pct, otr_average_reading, motr_notes_with_zero_count, motr_excluding_zero_mean, motr_non_zero_pct, motr_gt_eq_threshold_pct, motr_average_reading, etr_mtr_notes_with_zero_count, etr_mtr_excluding_zero_mean, etr_mtr_non_zero_pct, etr_mtr_gt_eq_threshold_pct, etr_mtr_average_reading, metr_mmtr_notes_with_zero_count, metr_mmtr_excluding_zero_mean, metr_mmtr_non_zero_pct, metr_mmtr_gt_eq_threshold_pct, metr_mmtr_average_reading, crn_notes_with_zero_count, crn_excluding_zero_mean, crn_non_zero_pct, crn_gt_eq_threshold_pct, crn_average_reading, med_notes_with_zero_count, med_excluding_zero_mean, med_non_zero_pct, med_gt_eq_threshold_pct, med_average_reading, tape_notes_with_zero_count, tape_excluding_zero_mean, tape_non_zero_pct, tape_gt_eq_threshold_pct, tape_average_reading, fol_notes_with_zero_count, fol_excluding_zero_mean, fol_non_zero_pct, fol_gt_eq_threshold_pct, fol_average_reading, sfol_notes_with_zero_count, sfol_excluding_zero_mean, sfol_non_zero_pct, sfol_gt_eq_threshold_pct, sfol_average_reading, mfol_notes_with_zero_count, mfol_excluding_zero_mean, mfol_non_zero_pct, mfol_gt_eq_threshold_pct, mfol_average_reading, iwb_notes_with_zero_count, iwb_excluding_zero_mean, iwb_non_zero_pct, iwb_gt_eq_threshold_pct, iwb_average_reading, iwf_notes_with_zero_count, iwf_excluding_zero_mean, iwf_non_zero_pct, iwf_gt_eq_threshold_pct, iwf_average_reading, gtf_notes_with_zero_count, gtf_excluding_zero_mean, gtf_non_zero_pct, gtf_gt_eq_threshold_pct, gtf_average_reading, gtb_notes_with_zero_count, gtb_excluding_zero_mean, gtb_non_zero_pct, gtb_gt_eq_threshold_pct, gtb_average_reading, gcwf_notes_with_zero_count, gcwf_excluding_zero_mean, gcwf_non_zero_pct, gcwf_gt_eq_threshold_pct, gcwf_average_reading, owl_notes_with_zero_count, owl_excluding_zero_mean, owl_non_zero_pct, owl_gt_eq_threshold_pct, owl_average_reading, stnf_notes_with_zero_count, stnf_excluding_zero_mean, stnf_non_zero_pct, stnf_gt_eq_threshold_pct, stnf_average_reading, stnb_notes_with_zero_count, stnb_excluding_zero_mean, stnb_non_zero_pct, stnb_gt_eq_threshold_pct, stnb_average_reading, fmtf_notes_with_zero_count, fmtf_excluding_zero_mean, fmtf_non_zero_pct, fmtf_gt_eq_threshold_pct, fmtf_average_reading, fmtb_notes_with_zero_count, fmtb_excluding_zero_mean, fmtb_non_zero_pct, fmtb_gt_eq_threshold_pct, fmtb_average_reading, ffit_notes_with_zero_count, ffit_excluding_zero_mean, ffit_non_zero_pct, ffit_gt_eq_threshold_pct, ffit_average_reading, ovip_notes_with_zero_count, ovip_excluding_zero_mean, ovip_non_zero_pct, ovip_gt_eq_threshold_pct, ovip_average_reading, ovis_notes_with_zero_count, ovis_excluding_zero_mean, ovis_non_zero_pct, ovis_gt_eq_threshold_pct, ovis_average_reading})

		return true
	})
	resp = GetVelocityAggWearCategoryDataResponse{total, resp_payload}
	return &resp

}