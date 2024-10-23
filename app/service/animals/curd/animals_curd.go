package curd

import (
	"catface/app/model"
	"catface/app/utils/model_handler"
	"catface/app/utils/query_handler"
	"fmt"
	"strconv"
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

func (a *AnimalsCurd) List(attrs string, gender string, breed string, sterilzation string, status string, num int, skip int, userId int) (temp []model.AnimalWithLikeList) {
	validSelectedFields := getSelectAttrs(attrs)
	genderArray := query_handler.StringToUint8Array(gender)
	breedArray := query_handler.StringToUint8Array(breed)
	sterilzationArray := query_handler.StringToUint8Array(sterilzation)
	statusArray := query_handler.StringToUint8Array(status)

	if num == 0 {
		num = 10
	}

	animals := model.CreateAnimalFactory("").Show(validSelectedFields, genderArray, breedArray, sterilzationArray, statusArray, num, skip)

	for i := range animals {
		animalWithLike := model.AnimalWithLikeList{
			Animal: animals[i],
			Like:   model.CreateAnimalLikeFactory("").Liked(userId, int(animals[i].Id)),
		}
		temp = append(temp, animalWithLike)
	}
	return
}

func (a *AnimalsCurd) Detail(id string) *model.Animal {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		// TODO LOG
		fmt.Println("Detail id error:", err)
		return nil
	}

	return model.CreateAnimalFactory("mysql").ShowByID(idInt)
}
