package model_es

import (
	"bytes"
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/utils/data_bind"
	"catface/app/utils/model_handler"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func CreateAnimalESFactory(animal *model.Animal) *Animal {
	if animal == nil {
		return &Animal{}
	}
	return &Animal{
		Id:          animal.Id,
		Name:        animal.Name,
		NickNames:   animal.NickNamesList,
		Description: animal.Description,
	}
}

type Animal struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	NickNames   []string `json:"nick_names"`
	Description string   `json:"description"`

	// After handler
	NickNamesHighlight []string `json:"nick_names_highlight"`
}

func (a *Animal) IndexName() string {
	return "catface_animals"
}

func (a *Animal) InsertDocument() error {
	ctx := context.Background()

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index: a.IndexName(),
		// DocumentID: fmt.Sprintf("%d", a.Id),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	return nil
}

func (a *Animal) QueryDocumentsMatchAll(query string, num int) ([]Animal, error) {
	body := fmt.Sprintf(`{
  "size": %d, 
  "query": {
    "bool": {
      "should": [
        { "match": {"name": "%s" }},
        { "match": {"nick_names": "%s" }},
        { "match": {"description": "%s" }}
      ]
    }
  },
  "highlight": {
	"pre_tags": ["%v"],
	"post_tags": ["%v"],
    "fields": {
      "name": {},
      "nick_names": {
	 	"pre_tags": [""],
        "post_tags": [""] 
	  },
      "description": {
        "fragment_size" : 15
      }
    }
  }
}`, num, query, query, query, consts.PreTags, consts.PostTags)

	hits, err := model_handler.SearchRequest(body, a.IndexName())
	if err != nil {
		return nil, err
	}

	var animals []Animal
	for _, hit := range hits {
		data := model_handler.MergeSouceWithHighlight(hit.(map[string]interface{}))

		var animal Animal
		if err := data_bind.ShouldBindFormMapToModel(data, &animal); err != nil {
			continue
		}

		animals = append(animals, animal)
	}

	return animals, nil
}
