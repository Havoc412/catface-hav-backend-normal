package web

import (
	"catface/app/global/consts"
	"catface/app/global/errcode"
	"catface/app/global/variable"
	"catface/app/model"
	"catface/app/utils/model_handler"
	"catface/app/utils/query_handler"
	"catface/app/utils/response"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Animals struct { // INFO 起到一个标记的作用，这样 web.xxx 的时候不同模块就不会命名冲突了。
}

func buildSelectAttrs(db *gorm.DB, attrs string) *gorm.DB {
	if len(attrs) > 0 {
		// 1. 获取空 Field
		fieldMap := model_handler.GetModelField(model.Animal{})

		// 2. 开始检查请求字段
		attrsArray := query_handler.StringToStringArray(attrs)
		for _, attr := range attrsArray {
			if attr == "*" { // 不需要过滤，直接返回
				return db
			} else if attr == "avatar" {
				fieldMap["avatar_height"] = true
				fieldMap["avatar_width"] = true
			}
			// 过滤 无效 的请求字段
			if _, ok := fieldMap[attr]; ok {
				fieldMap[attr] = true
				continue
			}
		}

		// 3. 装填字段，并 Select
		var validSelectedFields []string
		for key, value := range fieldMap {
			if value {
				validSelectedFields = append(validSelectedFields, key)
			}
		}
		db = db.Select(validSelectedFields)
	}
	return db
}

/**
 * @description: 通过检查字段的方式构建 Where 函数。
 * @param {*gorm.DB} db
 * @param {map[string][]uint8} conditions
 * @return {*}
 */
func buildQuery(db *gorm.DB, conditions map[string][]uint8) *gorm.DB {
	for field, values := range conditions {
		if len(values) == 0 || len(values) == 1 && values[0] == 0 {
			continue
		}
		db = db.Where(field+" in (?)", values)
	}
	return db
}

func (a *Animals) List(context *gin.Context) {
	// 1. Get Params
	attrs := context.GetString(consts.ValidatorPrefix + "attrs")
	gender := query_handler.StringToUint8Array(context.GetString(consts.ValidatorPrefix + "gender"))
	breed := query_handler.StringToUint8Array(context.GetString(consts.ValidatorPrefix + "breed"))
	sterilzation := query_handler.StringToUint8Array(context.GetString(consts.ValidatorPrefix + "sterilzation"))
	status := query_handler.StringToUint8Array(context.GetString(consts.ValidatorPrefix + "status"))
	num := context.GetFloat64(consts.ValidatorPrefix + "num")
	skip := context.GetFloat64(consts.ValidatorPrefix + "skip")

	// 创建条件映射
	conditions := map[string][]uint8{
		"gender":        gender,
		"breed":         breed,
		"sterilization": sterilzation,
		"status":        status,
	}

	// 2. Select & Filter
	if num == 0 {
		num = 10
	}
	db := variable.GormDbMysql.Table("animals").Limit(int(num)).Offset(int(skip))
	db = buildSelectAttrs(db, attrs)
	db = buildQuery(db, conditions)

	// 3. Find
	var animals []model.Animal
	err := db.Find(&animals).Error
	if err != nil {
		response.Fail(context, errcode.ErrAnimalSqlFind, errcode.ErrMsg[errcode.ErrAnimalSqlFind], err) // UPDATE consts ?
	} else {
		response.Success(context, consts.CurdStatusOkMsg, animals)
	}
}

func (a *Animals) Detail(context *gin.Context) {
	// 1. Get Params
	anmId, err := strconv.Atoi(context.Param("anm_id"))
	if err != nil {
		response.Fail(context, errcode.ErrAnimalSqlFind, errcode.ErrMsg[errcode.ErrAnimalSqlFind], err)
		return
	}
	fmt.Println("anmId:", anmId)

	// 2. Select & Filter
	var animal model.Animal
	err = variable.GormDbMysql.Table("animals").Model(&animal).Where("id = ?", anmId).Scan(&animal).Error // TIP GORM.First 采取默认的
	if err != nil {
		response.Fail(context, errcode.ErrAnimalSqlFind, errcode.ErrMsg[errcode.ErrAnimalSqlFind], err) // UPDATE consts ?
	} else {
		response.Success(context, consts.CurdStatusOkMsg, animal)
	}

}
