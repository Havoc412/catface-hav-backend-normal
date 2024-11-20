package curd

import (
	"catface/app/model"
	"catface/app/model_es"
	"catface/app/utils/gorm_v2"
	"catface/app/utils/model_handler"
	"catface/app/utils/query_handler"
	"fmt"
)

func CreateAnimalsCurdFactory() *AnimalsCurd {
	return &AnimalsCurd{model.CreateAnimalFactory("")}
}

type AnimalsCurd struct {
	animals *model.Animal // INFO 难道数据就是存储到这里？
}

func getSelectAttrs(attrs string) (validSelectedFields []string) {
	if len(attrs) == 0 {
		return nil
	}
	// 1. 获取空 Field
	fieldMap := model_handler.GetModelField(model.Animal{})

	// 2. 开始检查请求字段
	attrsArray := query_handler.StringToStringArray(attrs)
	for _, attr := range attrsArray {
		if attr == "*" { // 不需要过滤，直接返回
			return nil
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

	// 3. 装填字段
	for key, value := range fieldMap {
		if value {
			validSelectedFields = append(validSelectedFields, key)
		}
	}
	return
}

func (a *AnimalsCurd) List(attrs string, gender string, breed string, sterilization string, status string, department string, notInIds []int64, num int, skip int, userId int) (temp []model.AnimalWithLikeList) {
	validSelectedFields := getSelectAttrs(attrs)
	genderArray := query_handler.StringToUint8Array(gender)
	breedArray := query_handler.StringToUint8Array(breed)
	sterilizationArray := query_handler.StringToUint8Array(sterilization)
	statusArray := query_handler.StringToUint8Array(status)
	departmentArray := query_handler.StringToUint8Array(department)

	if num == 0 {
		num = 10
	}

	animals := model.CreateAnimalFactory("").Show(validSelectedFields, genderArray, breedArray, sterilizationArray, statusArray, departmentArray, notInIds, num, skip)

	// 状态记录
	var likeRes []bool
	var err error

	if userId > 0 {
		// Like：批量查询
		var animalIds []int
		for _, animal := range animals {
			animalIds = append(animalIds, int(animal.Id))
		}
		likeRes, err = model.CreateAnimalLikeFactory("").LikedBatch(userId, animalIds)
	}

	if err == nil && userId > 0 {
		for i := range animals {
			animalWithLike := model.AnimalWithLikeList{
				Animal: animals[i],
				Like:   likeRes[i],
			}
			temp = append(temp, animalWithLike)
		}
	} else {
		for i := range animals {
			animalWithLike := model.AnimalWithLikeList{
				Animal: animals[i],
			}
			temp = append(temp, animalWithLike)
		}
	}

	return
}

func (a *AnimalsCurd) ShowByName(attrs string, name string) (temp []model.AnimalWithNickNameHit) {
	validSelectedFields := getSelectAttrs(attrs)

	animals := model.CreateAnimalFactory("").ShowByName(name, validSelectedFields...)

	for _, animal := range animals {
		animalWithNameHit := model.AnimalWithNickNameHit{
			Animal:      animal,
			NickNameHit: !gorm_v2.IsLikePatternMatch(animal.Name, name), // 通过对比 name，然后取反；主要是不想让 SQL 过于复杂，、处理起来也麻烦。
		}
		temp = append(temp, animalWithNameHit)
	}

	return
}

func (a *AnimalsCurd) Detail(id int64) *model.Animal {

	return model.CreateAnimalFactory("mysql").ShowByID(id)
}

func (a *AnimalsCurd) MatchAll(query string, num int) (tmp []model.Animal) {
	// STAGE 1. ES 查询
	animalsFromES, err := model_es.CreateAnimalESFactory(nil).QueryDocumentsMatchAll(query, num)
	if err != nil {
		fmt.Println("ES Query error:", err)
		return nil
	}

	var ids []int64
	for _, animal := range animalsFromES {
		ids = append(ids, animal.Id)
	}

	// STAGE 2. MySQL 补充信息
	animalsFromSQL := model.CreateAnimalFactory("").ShowByIDs(ids, "id", "avatar", "status", "department")

	// 3. 合并信息
	for _, animalFromES := range animalsFromES {
		for _, animal := range animalsFromSQL {
			if animal.Id == animalFromES.Id {
				animal.NickNamesList = animalFromES.NickNames
				animal.NickNamesHighlight = animalFromES.NickNamesHighlight
				animal.Description = animalFromES.Description
				animal.Name = animalFromES.Name
				tmp = append(tmp, animal)
			}
		}
	}

	return
}
