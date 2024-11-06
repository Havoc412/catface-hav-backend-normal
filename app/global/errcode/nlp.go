package errcode

const (
	ErrNoContent = ErrNlp + iota
)

func NlpMsgInit(m msg) {
	m[ErrNoContent] = "内容为空"
}
