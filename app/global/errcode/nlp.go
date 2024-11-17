package errcode

const (
	ErrNoContent = ErrNlp + iota
	ErrNoDocFound
	ErrPythonServierDown
)

func NlpMsgInit(m msg) {
	m[ErrNoContent] = "内容为空"
	m[ErrNoDocFound] = "没有找到相关文档"
}

func NlpMsgUserInit(m msg) {
	m[ErrNoContent] = "请输入内容"
	m[ErrNoDocFound] = "小护没有在知识库中找到相关文档。😿"
	m[ErrPythonServierDown] = "小护的🐍python服务挂了，请稍后再试。😿"
}
