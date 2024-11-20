package model_res

// TODO DELETE 这个文件或许没什么意义了。
type DocInterface interface {
	ToString() string
}

type DocBase struct {
	Type string `json:"type"`
}

func (d DocBase) ToString() string {
	return ""
}
