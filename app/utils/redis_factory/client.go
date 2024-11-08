package redis_factory

import (
	"catface/app/core/event_manage"
	"catface/app/global/my_errors"
	"catface/app/global/variable"
	"catface/app/utils/yml_config"
	"catface/app/utils/yml_config/ymlconfig_interf"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

var redisPool *redis.Pool
var configYml ymlconfig_interf.YmlConfigInterf

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	configYml = yml_config.CreateYamlFactory()
	redisPool = initRedisClientPool()
}
func initRedisClientPool() *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     configYml.GetInt("Redis.MaxIdle"),                        //最大空闲数
		MaxActive:   configYml.GetInt("Redis.MaxActive"),                      //最大活跃数
		IdleTimeout: configYml.GetDuration("Redis.IdleTimeout") * time.Second, //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redis.Conn, error) {
			//此处对应redis ip及端口号
			conn, err := redis.Dial("tcp", configYml.GetString("Redis.Host")+":"+configYml.GetString("Redis.Port"))
			if err != nil {
				variable.ZapLog.Error(my_errors.ErrorsRedisInitConnFail + err.Error())
				return nil, err
			}
			auth := configYml.GetString("Redis.Auth") //通过配置项设置redis密码
			if len(auth) >= 1 {
				if _, err := conn.Do("AUTH", auth); err != nil {
					_ = conn.Close()
					variable.ZapLog.Error(my_errors.ErrorsRedisAuthFail + err.Error())
				}
			}
			_, _ = conn.Do("select", configYml.GetInt("Redis.IndexDb"))
			return conn, err
		},
	}
	// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
	eventManageFactory := event_manage.CreateEventManageFactory()
	if _, exists := eventManageFactory.Get(variable.EventDestroyPrefix + "Redis"); exists == false {
		eventManageFactory.Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
			_ = redisPool.Close()
		})
	}
	return redisPool
}

// 从连接池获取一个redis连接
func GetOneRedisClient() *RedisClient {
	maxRetryTimes := configYml.GetInt("Redis.ConnFailRetryTimes")
	var oneConn redis.Conn
	for i := 1; i <= maxRetryTimes; i++ {
		oneConn = redisPool.Get()
		// 首先通过执行一个获取时间的命令检测连接是否有效，如果已有的连接无法执行命令，则重新尝试连接到redis服务器获取新的连接池地址
		// 连接不可用可能会发生的场景主要有：服务端redis重启、客户端网络在有线和无线之间切换等
		if _, replyErr := oneConn.Do("time"); replyErr != nil {
			//fmt.Printf("连接已经失效(出错)：%+v\n", replyErr.Error())
			// 如果已有的redis连接池获取连接出错(官方库的说法是连接不可用)，那么继续使用从新初始化连接池
			initRedisClientPool()
			oneConn = redisPool.Get()
		}

		if err := oneConn.Err(); err != nil {
			//variable.ZapLog.Error("Redis：网络中断,开始重连进行中..." , zap.Error(oneConn.Err()))
			if i == maxRetryTimes {
				variable.ZapLog.Error(my_errors.ErrorsRedisGetConnFail, zap.Error(oneConn.Err()))
				return nil
			}
			//如果出现网络短暂的抖动，短暂休眠后，支持自动重连
			time.Sleep(time.Second * configYml.GetDuration("Redis.ReConnectInterval"))
		} else {
			break
		}
	}
	return &RedisClient{oneConn}
}

// 定义一个redis客户端结构体
type RedisClient struct {
	client redis.Conn
}

// 为redis-go 客户端封装统一操作函数入口
func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
	return r.client.Do(cmd, args...)
}

// 释放连接到连接池
func (r *RedisClient) ReleaseOneRedisClient() {
	_ = r.client.Close()
}

//  封装几个数据类型转换的函数

// bool 类型转换
func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

// string 类型转换
func (r *RedisClient) String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// string map 类型转换
func (r *RedisClient) StringMap(reply interface{}, err error) (map[string]string, error) {
	return redis.StringMap(reply, err)
}

// strings 类型转换
func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

// Float64 类型转换
func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

// int 类型转换
func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// int64 类型转换
func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// int map 类型转换
func (r *RedisClient) IntMap(reply interface{}, err error) (map[string]int, error) {
	return redis.IntMap(reply, err)
}

// Int64Map 类型转换
func (r *RedisClient) Int64Map(reply interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(reply, err)
}

// int64s 类型转换
func (r *RedisClient) Int64s(reply interface{}, err error) ([]int64, error) {
	return redis.Int64s(reply, err)
}

// uint64 类型转换
func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

// Bytes 类型转换
func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

// 以上封装了很多最常见类型转换函数，其他您可以参考以上格式自行封装

// 配合 List 类型使用; 配合 lpush 来获取返回的信息。 // UPDATE 感觉 前后的 [] 机制，应该是 go 导致的，em，肯定能优化，
func (r *RedisClient) Int64sFromList(reply interface{}, err error) ([]int64, error) {
	values, err := redis.Values(reply, err)
	if err != nil {
		return nil, err
	}
	// fmt.Println("len", len(values))

	if len(values) == 0 {
		return nil, fmt.Errorf("empty values")
	}

	// 假设 values 只有一个元素，且该元素是一个字符串 // INFO 这里返回的是一个 []uint8 的类型，但是处理的时候 string化 就会 导致前后的 '[]'
	strValue := string(values[0].([]byte))
	// fmt.Println("Original string:", strValue)

	// 去掉前后的括号
	strValue = strings.TrimPrefix(strValue, "[")
	strValue = strings.TrimSuffix(strValue, "]")
	// fmt.Println("Trimmed string:", strValue)

	// 拆分成单独的数字字符串
	strValues := strings.Split(strValue, " ")
	// fmt.Println("Split values:", strValues)

	var result []int64
	for _, str := range strValues {
		// 将字符串转换为 int64
		intValue, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}

		// 将转换后的 int64 值添加到结果切片中
		result = append(result, intValue)
	}

	return result, nil
}
