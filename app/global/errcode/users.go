package errcode

const (
	ErrWeixinApi = iota + ErrUser
)

func UserMsgInit(m msg) {
	m[ErrWeixinApi] = "微信接口调用失败"
}
