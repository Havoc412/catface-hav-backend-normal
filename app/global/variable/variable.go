package variable

import (
	"catface/app/global/my_errors"
	"catface/app/utils/llm_factory"
	"catface/app/utils/snow_flake/snowflake_interf"
	"catface/app/utils/yml_config/ymlconfig_interf"
	"log"
	"os"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BasePath           string                  // 定义项目的根目录
	EventDestroyPrefix = "Destroy_"            //  程序退出时需要销毁的事件前缀
	ConfigKeyPrefix    = "Config_"             //  配置文件键值缓存时，键的前缀
	DateFormat         = "2006-01-02 15:04:05" //  设置全局日期时间格式

	// INFO 全局日志指针
	ZapLog *zap.Logger
	// 全局配置文件
	ConfigYml       ymlconfig_interf.YmlConfigInterf // 全局配置文件指针
	ConfigGormv2Yml ymlconfig_interf.YmlConfigInterf // 全局配置文件指针
	PromptsYml      ymlconfig_interf.YmlConfigInterf // 提示词配置文件

	//gorm 数据库客户端，如果您操作数据库使用的是gorm，请取消以下注释，在 bootstrap>init 文件，进行初始化即可使用
	GormDbMysql      *gorm.DB // 全局gorm的客户端连接
	GormDbSqlserver  *gorm.DB // 全局gorm的客户端连接
	GormDbPostgreSql *gorm.DB // 全局gorm的客户端连接

	//雪花算法全局变量
	SnowFlake snowflake_interf.InterfaceSnowFlake

	//websocket
	WebsocketHub              interface{}
	WebsocketHandshakeSuccess = `{"code":200,"msg":"ws连接成功","data":""}`
	WebsocketServerPingMsg    = "Server->Ping->Client"

	//casbin 全局操作指针
	Enforcer *casbin.SyncedEnforcer

	// GLM 全局客户端集中管理
	GlmClientHub *llm_factory.GlmClientHub

	// ES 全局客户端
	ElasticClient *elasticsearch.Client
)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}

func init() {
	// 1. 初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			// 替换 \ 为 /，然后移除 /test 及其后的内容
			curPath = strings.ReplaceAll(curPath, "\\", "/")
			parts := strings.Split(curPath, "/test")
			BasePath = parts[0]
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
