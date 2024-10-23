package test

import (
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ElasticClient *elasticsearch.Client
	StudentIndex  string = "student_index"
)

func InitElastic() {
	var err error
	ElasticClient, err = elasticsearch.NewClient(elasticsearch.Config{
		// Addresses: []string{"http://113.44.68.213:9200"},
		Addresses: []string{"http://127.0.0.1:9200"},
		// Username:  "elastic",
		// Password:  "U8n61yn*Sp4Kvbuqo_K8",
	})
	if err != nil {
		panic(err)
	}
}

type Student struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
	School  string `json:"school,omitempty"`
}
