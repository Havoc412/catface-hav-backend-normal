package gorm_v2

import (
	"regexp"

	"gorm.io/gorm"
)

/**
 * @description: 通过检查字段的方式构建 Where 函数。
 * @param {*gorm.DB} db
 * @param {map[string][]uint8} conditions
 * @return {*}
 */
// INFO 特性，源于 MySQL 键值 index from 1，
// 同时 go 在解析参数之时，对于 Query 为空的情况会得到 [0] 的结果，
// 所以就可以用这种方式简单的过滤掉。
func BuildWhere(db *gorm.DB, conditions map[string][]uint8) *gorm.DB {
	for field, values := range conditions {
		if len(values) == 0 || len(values) == 1 && values[0] == 0 {
			continue
		}
		db = db.Where(field+" in (?)", values)
	}
	return db
}

// isLikePatternMatch 检查字符串是否匹配 LIKE '%name%' 模式
func IsLikePatternMatch(input, pattern string) bool {
	// 构建正则表达式
	regexPattern := ".*" + regexp.QuoteMeta(pattern) + ".*"

	// 编译正则表达式
	re := regexp.MustCompile(regexPattern)

	// 检查输入字符串是否匹配正则表达式
	return re.MatchString(input)
}
