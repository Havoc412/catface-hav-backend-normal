package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func TestAddIndex(t *testing.T) {
	InitElastic()

	stu := Student{
		Id:      "student_1",
		Name:    "Geek 1",
		Age:     22,
		Address: "Wuhan",
		School:  "Hebei School",
	}

	marshal, err := json.Marshal(stu)
	if err != nil {
		fmt.Println(err)
		return
	}

	indexReq := esapi.IndexRequest{
		Index:      StudentIndex,
		DocumentID: fmt.Sprintf("%s_%s", StudentIndex, stu.Id),
		Body:       bytes.NewReader(marshal),
		Refresh:    "true",
	}
	if _, err = indexReq.Do(context.Background(), ElasticClient); err != nil {
		fmt.Println(err)
		return
	}
}

func TestSearch(t *testing.T) {
	InitElastic()

	students := make([]*Student, 0)

	keyWord := "Yan"

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"match": map[string]interface{}{"name": keyWord}},
					{"match": map[string]interface{}{"address": keyWord}},
					{"match": map[string]interface{}{"school": keyWord}},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Println(err)
		return
	}

	res, err := ElasticClient.Search(
		ElasticClient.Search.WithContext(context.Background()),
		ElasticClient.Search.WithIndex(StudentIndex),
		ElasticClient.Search.WithBody(&buf),
		ElasticClient.Search.WithTrackTotalHits(true),
		ElasticClient.Search.WithPretty(),
	)
	if err != nil || res.IsError() {
		fmt.Println(err)
		return
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println(err)
		return
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		if _, v := hit.(map[string]interface{})["_source"]; v {
			var model Student
			body, err := json.Marshal(hit.(map[string]interface{})["_source"])
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := json.Unmarshal(body, &model); err != nil {
				fmt.Println(err)
				continue
			}
			students = append(students, &model)
		}
	}

	for _, s := range students {
		fmt.Printf("%#v \n", s)
	}
}
