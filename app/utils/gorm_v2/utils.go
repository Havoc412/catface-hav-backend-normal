package gorm_v2

import "gorm.io/gorm"

/**
 * @description: 通过检查字段的方式构建 Where 函数。
 * @param {*gorm.DB} db
 * @param {map[string][]uint8} conditions
 * @return {*}
 */
func BuildWhere(db *gorm.DB, conditions map[string][]uint8) *gorm.DB {
	for field, values := range conditions {
		if len(values) == 0 || len(values) == 1 && values[0] == 0 {
			continue
		}
		db = db.Where(field+" in (?)", values)
	}
	return db
}
