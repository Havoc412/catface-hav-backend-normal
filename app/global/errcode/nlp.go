package errcode

const (
	ErrNoContent = ErrNlp + iota
	ErrNoDocFound
)

func NlpMsgInit(m msg) {
	m[ErrNoContent] = "内容为空"
	m[ErrNoDocFound] = "没有找到相关文档"
}

func NlpMsgUserInit(m msg) {
	m[ErrNoContent] = "请输入内容"
	m[ErrNoDocFound] = "没有找到相关文档"
}
