package test

import (
	model "catface/app/model"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "Havocantelope412#", "113.44.68.213", "3306", "hav_cats") // danger MySQL
	fmt.Println("dsn:", dsn)
	dbMySQL, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = dbMySQL
}

const (
	list = `[{
		"Name": "小北",
		"Birthday": "2021-01-01",
		"Gender": 3,
		"Breed": 10,
		"Status": 2,
		"Sterilization": 2,
		"Description": "曾经的小北门女神；喜欢贴贴的可爱大美猫；现已领养；",
		"Latitude": 30.532645,
		"Longitude": 114.367661,
		"Avatar": "1.jpg",
		"AvatarHeight":   178,
		"AvatarWidth":    118,
		"HeadImg":        "1.jpg",
		"Photos":         "",
		"ActivityRadius": 100,
		"Tags": "",
		"NickName": "胖北,北子",
		"face_breeds":     "10,8,2",
		"face_breed_probs": "0.9,0.09,0.01"
	}, {
		"Name": "打三",
		"Birthday": "2019-01-01",
		"Gender": 3,
		"Breed": 10,
		"Status": 1,
		"Sterilization": 1,
		"Description": "常驻梅园区域；高冷的一尊佛；阅历深厚的老学姐；",
		"Latitude": 30.532645,
		"Longitude": 114.367661,
		"Avatar": "2.jpg",
		"AvatarHeight":   160,
		"AvatarWidth":    213,
		"HeadImg":        "2.jpg",
		"Photos":         "",
		"ActivityRadius": 100,
		"NickName": "打人三花,小花",
		"Tags": "",
		"face_breeds":     "10,8,2",
		"face_breed_probs": "0.9,0.09,0.01"
	}, {
		"Name": "猪皮",
		"Birthday": "2017-01-01",
		"Gender": 2,
		"Breed": 2,
		"Status": 1,
		"Sterilization": 1,
		"Description": "信部资深学长；很有个性的司马脸；每晚选择一位大学牲翻牌子。",
		"Latitude": 30.532645,
		"Longitude": 114.367661,
		"Avatar": "4.jpg",
		"AvatarHeight":   191,
		"AvatarWidth":    160,		
		"HeadImg":        "4.jpg",
		"Photos":         "0.png,1.jpg,3.jpg",
		"ActivityRadius": 100,
		"nick_names": "猜皮,猪",
		"Tags": "臭脸,猜皮,玉玉",
		"face_breeds":     "2,6,4",
		"face_breed_probs": "0.9,0.09,0.01"
	}, {
		"Name": "切糕",
		"Birthday": "2022-01-01",
		"Gender": 2,
		"Breed": 8,
		"Tags": "",
		"Status": 1,
		"Sterilization": 1,
		"Description": "脾气老实的男妈妈；活动在湖滨区域；身边还有一只狸花猫叫做小蝶；",
		"Latitude": 30.532645,
		"Longitude": 114.367661,
		"Avatar": "6.jpg",
		"AvatarHeight":   408,
		"AvatarWidth":    306,
		"HeadImg":        "6.jpg",
		"Photos":         "",
		"ActivityRadius": 100,
		"face_breeds":     "8,7,5",
		"face_breed_probs": "0.8,0.19,0.01"
	}]`
)

func TestCreateAnimal(t *testing.T) {
	Init()

	err := DB.AutoMigrate(&model.Animal{})
	if err != nil {
		fmt.Println("autoMigrateTable error:", err)
	}

	// 创建一个新的 Animal 实例
	animal := model.Animal{
		Name:           "斜刘海",
		Birthday:       "2021-01-01",
		Gender:         3,
		Breed:          3, // 示例品种 ID
		Sterilization:  2, // 已绝育
		NickNames:      "cc5,西西务,西西务小姐",
		Status:         1, // 示例状态
		Description:    "信部资深学姐,优雅的代名。",
		Latitude:       30.532645,
		Longitude:      114.367661,
		Avatar:         "0.jpg",
		AvatarHeight:   239,
		AvatarWidth:    160,
		HeadImg:        "0.jpg",
		Photos:         "",
		ActivityRadius: 100,
		Tags:           "优雅",
		FaceBreeds:     "3,4,5",
		FaceBreedProbs: "0.9,0.05,0.05",
	}

	// 插入数据到数据库
	result := DB.Create(&animal)
	// 检查插入是否成功
	assert.Nil(t, result.Error)

	// 可以进一步检查数据是否正确插入,例如通过查询数据库来验证
	var foundAnimal model.Animal
	result = DB.First(&foundAnimal, animal.BaseModel.Id)
	assert.Nil(t, result.Error)
	assert.Equal(t, animal.Name, foundAnimal.Name)
}

func TestCreateAnimalList(t *testing.T) {
	Init()

	var animals []model.Animal
	err := json.Unmarshal([]byte(list), &animals)
	if err != nil {
		fmt.Println("JSON 解析错误:", err)
		return
	}

	for _, animal := range animals {
		fmt.Println("animal:", animal)
		animal := model.Animal{
			Name:           animal.Name,
			Birthday:       animal.Birthday,
			Gender:         animal.Gender,
			Breed:          animal.Breed,
			Sterilization:  animal.Sterilization,
			NickNames:      animal.NickNames,
			Status:         animal.Status,
			Description:    animal.Description,
			Latitude:       animal.Latitude,
			Longitude:      animal.Longitude,
			Avatar:         animal.Avatar,
			AvatarHeight:   animal.AvatarHeight,
			AvatarWidth:    animal.AvatarWidth,
			HeadImg:        animal.HeadImg,
			Photos:         animal.Photos,
			ActivityRadius: animal.ActivityRadius,
			FaceBreeds:     animal.FaceBreeds,
			FaceBreedProbs: animal.FaceBreedProbs,
			Tags:           animal.Tags,
		}

		fmt.Println("TEST: ", animal.FaceBreeds, animal.FaceBreedProbs)

		result := DB.Create(&animal)

		// 检查插入是否成功
		assert.Nil(t, result.Error)
	}
}
