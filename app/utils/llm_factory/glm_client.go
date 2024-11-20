package llm_factory

import (
	"catface/app/global/errcode"
	"time"

	"github.com/yankeguo/zhipu"
)

// INFO 维护 GLM Client 与用户之间的客户端消息队列，也就是在 "github.com/yankeguo/zhipu" 的基础上实现一层封装。

type GlmClientHub struct {
	Idle             int // 最大连接数
	Active           int // 最大活跃数
	ApiKey           string
	DefaultModelName string
	InitPrompt       string
	Clients          map[string]*ClientInfo
	LifeTime         time.Duration // 最长待机周期
}

type ClientInfo struct {
	Client     *zhipu.ChatCompletionService
	UserQuerys []string
	LastUsed   time.Time
}

func InitGlmClientHub(maxIdle, maxActive, lifetime int, apiKey, defaultModelName, initPrompt string) *GlmClientHub {
	hub := &GlmClientHub{
		Idle:             maxIdle,
		Active:           maxActive,
		ApiKey:           apiKey,
		DefaultModelName: defaultModelName,
		InitPrompt:       initPrompt,
		Clients:          make(map[string]*ClientInfo),
		LifeTime:         time.Duration(lifetime) * time.Second,
	}
	go hub.cleanupRoutine() // 启动定时器清理过期会话。
	return hub
}

const (
	GlmModeSimple = iota
	GlmModeKnowledgeHub
)

/**
 * @description: 鉴权用户之后，根据其 ID 来从 map池 里获取之前的连接。
 * // UPDATE 现在只是单用户单连接（也就是只支持“同时只有一个对话”），之后可以考虑扩展【消息队列】的封装方式。
 * 默认启用的是 没有预设的 prompt 的空。
 * @param {string} token： // TODO 如何在 token 中保存信息？
 * @return {*}
 */
func (g *GlmClientHub) GetOneGlmClientInfo(token string, mode int) (clientInfo *ClientInfo, code int) {
	if info, ok := g.Clients[token]; ok {
		info.LastUsed = time.Now() // INFO 刷新生命周期
		return info, 0
	}

	// 空闲数检查
	if g.Idle > 0 && g.Active > 0 {
		g.Idle -= 1
		g.Active -= 1
	} else {
		code = errcode.ErrGlmBusy
		return
	}

	// Client Init
	preClient, err := zhipu.NewClient(zhipu.WithAPIKey(g.ApiKey))
	if err != nil {
		code = errcode.ErrGlmNewClientFail
		return
	}
	client := preClient.ChatCompletion(g.DefaultModelName)

	if mode == GlmModeKnowledgeHub {
		client.AddMessage(zhipu.ChatCompletionMessage{
			Role:    zhipu.RoleSystem, // TIP 使用 System 角色来初始化对话
			Content: g.InitPrompt,
		})
	}

	clientInfo = &ClientInfo{
		Client:   client,
		LastUsed: time.Now(),
	}
	g.Clients[token] = clientInfo
	return
}

/**
 * @description: 获取并返回 ClientInfo 的 Client 和 code。
 * @param {string} token
 * @param {int} mode
 * @return {(*zhipu.ChatCompletionService, int)}
 */
func (g *GlmClientHub) GetOneGlmClient(token string, mode int) (*zhipu.ChatCompletionService, int) {
	clientInfo, code := g.GetOneGlmClientInfo(token, mode)
	if clientInfo == nil || code != 0 {
		return nil, code
	}
	return clientInfo.Client, code
}

// cleanupRoutine 定期检查并清理超过 1 小时未使用的 Client
func (g *GlmClientHub) cleanupRoutine() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		g.cleanupClients()
	}
}

// cleanupClients 清理超过 1 小时未使用的 Client
func (g *GlmClientHub) cleanupClients() {
	now := time.Now()
	for token, info := range g.Clients {
		if now.Sub(info.LastUsed) > g.LifeTime {
			delete(g.Clients, token)
			g.Idle += 1
		}
	}
}

/**
 * @description: ws 服务完毕，进入待机状态。
 * @param {string} token
 * @return {*}
 * @Tip 对于临时使用的小功能，需要依次 defer 下面两个函数。
 */
func (g *GlmClientHub) UnavtiveOneGlmClient(token string) bool {
	if clientInfo, exists := g.Clients[token]; exists {
		g.Active -= 1
		clientInfo.LastUsed = time.Now()
		return true
	}
	return false
}

/**
 * @description: 显式地释放资源。
 * @param {string} token
 * @return {*}
 */
func (g *GlmClientHub) ReleaseOneGlmClient(token string) bool {
	if _, exists := g.Clients[token]; exists {
		delete(g.Clients, token)
		g.Idle += 1
		return true
	}
	return false
}

// TAG ClientInfo
func (c *ClientInfo) AddQuery(query string) {
	c.UserQuerys = append(c.UserQuerys, query)
}
