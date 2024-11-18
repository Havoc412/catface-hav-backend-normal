package llm_factory

import (
	"catface/app/global/errcode"
	"time"

	"github.com/yankeguo/zhipu"
)

// INFO 维护 GLM Client 与用户之间的客户端消息队列，也就是在 "github.com/yankeguo/zhipu" 的基础上实现一层封装。

type GlmClientHub struct {
	MaxIdle          int
	MaxActive        int
	ApiKey           string
	DefaultModelName string
	InitPrompt       string
	Clients          map[string]*ClientInfo
	LifeTime         time.Duration
}

type ClientInfo struct {
	Client   *zhipu.ChatCompletionService
	LastUsed time.Time
}

func InitGlmClientHub(maxIdle, maxActive, lifetime int, apiKey, defaultModelName, initPrompt string) *GlmClientHub {
	hub := &GlmClientHub{
		MaxIdle:          maxIdle,
		MaxActive:        maxActive,
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
func (g *GlmClientHub) GetOneGlmClient(token string, mode int) (client *zhipu.ChatCompletionService, code int) {
	if info, ok := g.Clients[token]; ok {
		info.LastUsed = time.Now() // INFO 刷新生命周期
		return info.Client, 0
	}

	// 空闲数检查
	if g.MaxIdle > 0 {
		g.MaxIdle -= 1
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
	client = preClient.ChatCompletion(g.DefaultModelName)

	if mode == GlmModeKnowledgeHub {
		client.AddMessage(zhipu.ChatCompletionMessage{
			Role:    zhipu.RoleSystem, // TIP 使用 System 角色来初始化对话
			Content: g.InitPrompt,
		})
	}

	g.Clients[token] = &ClientInfo{
		Client:   client,
		LastUsed: time.Now(),
	}
	return
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
			g.MaxIdle += 1
		}
	}
}

/**
 * @description: 显式地释放资源。
 * @param {string} token
 * @return {*}
 */
func (g *GlmClientHub) ReleaseOneGlmClient(token string) {
	delete(g.Clients, token)
	g.MaxIdle += 1
}
