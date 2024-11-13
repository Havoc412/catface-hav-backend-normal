package model_handler

import (
	"catface/app/global/variable"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/**
 * @description: 用于处理 ES-highlight 模块分析出来的 []String.
 * @param {[]interface{}} strs
 * @return {*} 将 []String 还原为 String，这里就是简单的拼接了一下。
 */
func TransStringSliceToString(strs []interface{}) string {
	var result string
	for _, str := range strs {
		if s, ok := str.(string); ok {
			result += s
		}
	}
	return result
}

// func concatKeywordsToSlice(highlightKeyword string, oriKeywords []string) []string {
// 	return append(oriKeywords, highlightKeyword)
// }

func MergeSouceWithHighlight(hit map[string]interface{}) map[string]interface{} {
	// 1. Get data
	source := hit["_source"].(map[string]interface{})
	highlight := hit["highlight"].(map[string]interface{})

	// 2. Merge data
	for k, v := range highlight {
		if _, ok := source[k]; ok {
			switch source[k].(type) {
			case string:
				source[k] = TransStringSliceToString(v.([]interface{}))
			case []interface{}:
				source[k+"_highlight"] = v // TODO 过滤，交给前端？
			}
		}
	}
	return source
}

/**
 * @description: 对 ES 发送的请求，返回结果。
 * @param {string} body
 * @param {string} index
 * @return {*}
 */
func SearchRequest(body string, index string) ([]interface{}, error) {
	ctx := context.Background()

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(body),
	}

	res, err := req.Do(ctx, variable.ElasticClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var k map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&k); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				k["error"].(map[string]interface{})["type"],
				k["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error extracting hits from response")
	}
	return hits, nil
}
