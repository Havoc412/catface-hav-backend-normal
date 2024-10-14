package routers

import (
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/http/middleware/cors"
	validatorFactory "catface/app/http/validator/core/factory"
	"catface/app/utils/gin_release"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 该路由主要设置门户类网站等前台路由  // INFO 写一些类似怕早冲的任务？

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	// 非调试模式（生产模式） 日志写到日志文件
	if variable.ConfigYml.GetBool("AppDebug") == false {
		//1.gin自行记录接口访问日志，不需要nginx，如果开启以下3行，那么请屏蔽第 34 行代码
		//gin.DisableConsoleColor()
		//f, _ := os.Create(variable.BasePath + variable.ConfigYml.GetString("Logs.GinLogName"))
		//gin.DefaultWriter = io.MultiWriter(f)

		//【生产模式】
		// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
		// 如果部署到生产环境，请使用以下模式：
		// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
		// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
		// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
		router = gin_release.ReleaseRouter()
	} else {
		// 调试模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	// 设置可信任的代理服务器列表,gin (2021-11-24发布的v1.7.7版本之后出的新功能)
	if variable.ConfigYml.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(variable.ConfigYml.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			variable.ZapLog.Error(consts.GinSetTrustProxyError, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	//根据配置进行设置跨域
	if variable.ConfigYml.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(cors.Next())
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello Havoc.")
	})

	//  创建一个门户类接口路由组
	vApi := router.Group("/api/v1/")
	{
		// // 模拟一个首页路由
		home := vApi.Group("home/")
		{
			// 第二个参数说明：
			// 1.它是一个表单参数验证器函数代码段，该函数从容器中解析，整个代码段略显复杂，但是对于使用者，您只需要了解用法即可，使用很简单，看下面 ↓↓↓
			// 2.编写该接口的验证器，位置：app/http/validator/api/home/news.go
			// 3.将以上验证器注册在容器：app/http/validator/common/register_validator/api_register_validator.go  18 行为注册时的键（consts.ValidatorPrefix + "HomeNews"）。那么获取的时候就用该键即可从容器获取
			home.GET("news", validatorFactory.Create(consts.ValidatorPrefix+"HomeNews"))
		}

	}
	return router
}
