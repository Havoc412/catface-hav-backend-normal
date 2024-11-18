package errcode

const (
	ErrNoContent = ErrNlp + iota
	ErrNoDocFound
	ErrPythonServierDown
	ErrGlmBusy
	ErrGlmHistoryLoss
	ErrGlmNewClientFail
)

func NlpMsgInit(m msg) {
	m[ErrNoContent] = "内容为空"
	m[ErrNoDocFound] = "没有找到相关文档"
	m[ErrGlmNewClientFail] = "GLM 新建客户端失败"
}

func NlpMsgUserInit(m msg) {
	m[ErrNoContent] = "请输入内容"
	m[ErrNoDocFound] = "小护没有在知识库中找到相关文档。😿"
	m[ErrPythonServierDown] = "小护的🐍python服务挂了，此功能暂时无法使用。😿"
	m[ErrGlmBusy] = "现在有太多人咨询小护，请稍后再来。"
	m[ErrGlmHistoryLoss] = "抱歉！小护找不到之前的会话记录了，我们重新开始新的对话吧。"
	// m[ErrGlmNewClientFail] = "小护新建客户端失败了，请稍后再来。"
}
