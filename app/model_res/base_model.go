package model_res

type DocInterface interface {
	ToString() string
}

type DocBase struct {
	Type string `json:"type"`
}

func (d DocBase) ToString() string {
	return ""
}
