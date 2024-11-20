package model_res

import (
	"catface/app/global/consts"
	"encoding/json"
)

func CreateNlpWebSocketResult(t string, data any) *NlpWebSocketResult {
	if t == "" {
		t = consts.AiMessageTypeText
	}

	return &NlpWebSocketResult{
		Type: t,
		Data: data,
	}
}

type NlpWebSocketResult struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (n *NlpWebSocketResult) JsonMarshal() []byte {
	data, _ := json.Marshal(n)
	return data
}
