package data_transfer

import (
	"catface/app/global/consts"
	"catface/app/global/variable"
	"catface/app/http/validator/core/interf"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 将验证器成员(字段)绑定到数据传输上下文，方便控制器获取
/**
本函数参数说明：
validatorInterface 实现了验证器接口的结构体
extra_add_data_prefix  验证器绑定参数传递给控制器的数据前缀
context  gin上下文
*/

func DataAddContext(validatorInterface interf.ValidatorInterface, extraAddDataPrefix string, context *gin.Context) *gin.Context {
	var tempJson interface{}
	if tmpBytes, err1 := json.Marshal(validatorInterface); err1 == nil {
		if err2 := json.Unmarshal(tmpBytes, &tempJson); err2 == nil {
			if value, ok := tempJson.(map[string]interface{}); ok {
				for k, v := range value {
					context.Set(extraAddDataPrefix+k, v)
				}
				// 此外给上下文追加三个键：created_at  、 updated_at  、 deleted_at ，实际根据需要自己选择获取相关键值
				curDateTime := time.Now().Format(variable.DateFormat)
				context.Set(extraAddDataPrefix+"created_at", curDateTime)
				context.Set(extraAddDataPrefix+"updated_at", curDateTime)
				context.Set(extraAddDataPrefix+"deleted_at", curDateTime)
				return context
			}
		}
	}
	return nil
}

// getSlice 是一个通用的辅助函数，用于从 context 中获取切片。
func getSlice(context *gin.Context, ValidatorPrefix string, key string, elemType reflect.Type) interface{} {
	if val, ok := context.Get(ValidatorPrefix + key); ok && val != nil {
		if slice, ok := val.([]interface{}); ok {
			result := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(slice))
			for _, item := range slice {
				if reflect.TypeOf(item) == elemType {
					result = reflect.Append(result, reflect.ValueOf(item))
				}
			}
			return result.Interface()
		}
	}
	return nil
}

// GetStringSlice 从 context 中获取字符串切片。
func GetStringSlice(context *gin.Context, key string) (ss []string) {
	if val := getSlice(context, consts.ValidatorPrefix, key, reflect.TypeOf("")); val != nil {
		ss = val.([]string)
	}
	// INFO 同时重新装载。
	context.Set(consts.ValidatorPrefix+key, ss)
	return
}

// GetIntSlice 从 context 中获取整数切片。
func GetIntSlice(context *gin.Context, key string) (ss []int) {
	if val := getSlice(context, consts.ValidatorPrefix, key, reflect.TypeOf(0)); val != nil {
		ss = val.([]int)
	}
	context.Set(consts.ValidatorPrefix+key, ss)
	return
}

// ConvertSliceToString 是一个泛型函数，可以接受任何类型的切片
func ConvertSliceToString[T any](slice []T) (string, error) {
	var strBuilder strings.Builder
	for i, v := range slice {
		if i > 0 {
			strBuilder.WriteString(",")
		}
		strBuilder.WriteString(fmt.Sprintf("%v", v))
	}
	return strBuilder.String(), nil
}
