package errcode

type msg map[int]string

var ErrMsg msg

var ErrMsgForUser msg

func init() {
	ErrMsg = make(msg)
	GeneralMsgInit(ErrMsg)
	AnimalMsgInit(ErrMsg)
	UserMsgInit(ErrMsg)
	EnocunterMsgInit(ErrMsg)
	NlpMsgInit(ErrMsg)
	KnowledgeMsgInit(ErrMsg)
	SubServiceMsgInit(ErrMsg)
	WsMsgInit(ErrMsg)

	// INGO
	ErrMsgForUser = make(msg)
	GeneralMsgUserInit(ErrMsgForUser)
	AnimalMsgUserInit(ErrMsgForUser)
	EncounterMsgUserInit(ErrMsgForUser)
	KnowledgeMsgUserInit(ErrMsgForUser)
	NlpMsgUserInit(ErrMsgForUser)
}

func GeneralMsgInit(m msg) {
	m[0] = ""
	m[ErrInvalidData] = "参数无效"
	m[ErrInternalError] = "内部服务器错误"
	m[ErrDataNoFound] = "无数据查询"
}

func GeneralMsgUserInit(m msg) {
	m[ErrServerDown] = "后端服务未启动，此功能暂时无法使用。😿" // 前端 rich-text 能够处理，有意思。
}
