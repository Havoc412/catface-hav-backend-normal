package errcode

const (
	ErrPythonService = ErrSubService + iota
)

func SubServiceMsgInit(m msg) {
	m[ErrPythonService] = "python微服务异常"
}
