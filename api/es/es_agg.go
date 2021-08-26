package api

import (
	"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	"path"
)

// Bucket contains how often a specific key was found in a term aggregation.
type Bucket struct {
	Key   interface{} `json:"key"`
	Count int         `json:"doc_count"`
}

func (c *Client) CompositeAggregate(index, doctype string, query map[string]interface{}, field string) ([]*Bucket, error) {
	return c.compositeAggregateAfter(index, doctype, query, field, nil)
}

var compositeSize = 500

func (c *Client) compositeAggregateAfter(index, doctype string, query map[string]interface{}, field string, after interface{}) ([]*Bucket, error) {
	var compositeResult []*Bucket
	request := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{
			"my_buckets": map[string]interface{}{
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
			},
		},
	}
	if after != nil {
		request["aggs"].(map[string]interface{})["my_buckets"].(map[string]interface{})["composite"].(map[string]interface{})["after"] = after
	}
	if query != nil {
		request["query"] = query
	}
	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %s", err)
	}
	apipath := path.Join(index, doctype) + "/_search"
	res, err := c.post(apipath, b)
	if err != nil {
		return nil, fmt.Errorf("could not get aggregations: %s", err)
	}
	result := struct {
		Aggregations struct {
			MyBuckets struct {
				Buckets []*struct {
					Key   map[string]interface{} `json:"key"`
					Count int                    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"my_buckets"`
		} `json:"aggregations"`
	}{}
	decoder := json.NewDecoder(bytes.NewReader(res))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode result: %s", err)
	}
	for _, bucket := range result.Aggregations.MyBuckets.Buckets {
		compositeResult = append(compositeResult, &Bucket{Key: bucket.Key[field], Count: bucket.Count})
	}
	if bucketLength := len(result.Aggregations.MyBuckets.Buckets); bucketLength > 0 {
		nextResult, err := c.compositeAggregateAfter(index, doctype, query, field, map[string]interface{}{
			field: result.Aggregations.MyBuckets.Buckets[bucketLength-1].Key[field],
		})
		if err != nil {
			return nil, err
		}
		compositeResult = append(compositeResult, nextResult...)
	}
	return compositeResult, nil
}