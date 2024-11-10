package oss

import "catface/app/global/variable"

const qiNiuBucket = "sakura-fox-asia"
const qiNiuDomain = "https://images.cengkehepler.top"

// 12. 初始化全局变量
var (
	qiNiuAccessKey string
	qiNiuSecretKey string
)

func init() {
	// 如果没有配置七牛云的key和secret，那么直接panic
	qiNiuAccessKey = variable.ConfigYml.GetString("QiNiu.AccessKey")
	qiNiuSecretKey = variable.ConfigYml.GetString("QiNiu.SecretKey")
	if qiNiuAccessKey == "" || qiNiuSecretKey == "" {
		panic("七牛云配置文件中 AccessKey 或 SecretKey 为空")
	}
}
