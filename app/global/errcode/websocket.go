package errcode

const (
	ErrWebsocketUpgradeFail = ErrWebSocket + iota
)

func WsMsgInit(m msg) {
	m[ErrWebsocketUpgradeFail] = "websocket升级失败"
}
