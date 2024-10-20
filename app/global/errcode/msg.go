package errcode

type msg map[int]string

var ErrMsg msg

func init() {
	ErrMsg = make(msg)
	GeneralMsgInit(ErrMsg)
	AnimalMsgInit(ErrMsg)
	UserMsgInit(ErrMsg)
}

func GeneralMsgInit(m msg) {
	m[0] = ""
	m[ErrInvalidData] = "参数无效"
	m[ErrInternalError] = "内部服务器错误"
	m[ErrDataNoFound] = "无数据查询"
}
