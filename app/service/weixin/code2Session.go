package weixin

import (
	"catface/app/global/variable"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

func Code2Session(js_code string) (*Response, error) {
	// return &Response{ // TEST
	// 	OpenId:     "open_id",
	// 	SessionKey: "session_key",
	// 	Errcode:    0,
	// 	Errmsg:     "ok",
	// }, nil

	appid := variable.ConfigYml.GetString("Weixin.AppId")
	appSecret := variable.ConfigYml.GetString("Weixin.AppSecret")
	grantType := variable.ConfigYml.GetString("Weixin.Code2Session.GrantType")

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s", appid, appSecret, js_code, grantType)

	// 创建一个新的HTTP请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// 解析响应体。
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response body: %v", err)
	}
	// check Success
	if response.Errcode != 0 {
		return nil, fmt.Errorf("error code: %d, error message: %s", response.Errcode, response.Errmsg)
	}

	return &response, nil
}
