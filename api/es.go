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
		resp.Count = int(hit.(map[string]interface{})["_source"].(map[string]interface{})["count"].(float64))
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


	for _, agg := range r["aggregations"].(map[string]interface{})["ALL_PROCESS_RUN_DATE_ID"].(map[string]interface{})["buckets"].([]interface{}){
		log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64)))
		//log.Printf(" PROCESS_RUN_DATE_ID=%d, count=%d", int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)))
		resp = append(resp, GetBetweenProcessRunDateIds{int(agg.(map[string]interface{})["key"].(map[string]interface{})["PROCESS_RUN_DATE_ID"].(float64)), int(agg.(map[string]interface{})["count"].(map[string]interface{})["value"].(float64))})
	}

	log.Println(strings.Repeat("=", 37))

	return &resp
}