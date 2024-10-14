package main

// INFO 🐱 开发时测试 标迁移；

import (
	model "catface/app/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // 这种写法是方柏包外使用

// 自动迁移表
func autoMigrateTable() {
	err := DB.AutoMigrate(&model.Animal{}, &model.Breed{}, &model.Sterilzation{}, &model.AnmStatus{}, &model.AnmGender{})
	if err != nil {
		fmt.Println("autoMigrateTable error:", err)
	}
}

func testInsertSterilzation() {
	// 定义状态数据
	statusesZH := []string{"不明", "未绝育", "已绝育"}
	statusesEN := []string{"unknown", "unsterilized", "sterilized"}

	for i := 0; i < len(statusesZH); i++ {
		sterilzation := model.Sterilzation{
			NameZh: statusesZH[i],
			NameEn: statusesEN[i],
		}
		// 3.
		tx := DB.Create(&sterilzation)
		if tx.Error != nil {
			fmt.Println("insert sterilzation error:", tx.Error)
		}
	}
}

func testInsertBreed() {
	// INFO 为方便之后扩展，将 unknown 默认为  1
	colorsZH := []string{"不明", "橘白", "奶牛", "白猫", "黑猫", "橘猫", "狸花", "狸白", "简州", "三花", "彩狸"}
	colorsEN := []string{"unknown", "orange", "cow", "white", "black", "orangeCat", "tabby", "tabbyWhite", "jianzhong", "threeColor", "colorCat"}
	for i := 0; i < len(colorsZH); i++ {
		breed := model.Breed{
			BriefModel: model.BriefModel{
				NameZh: colorsZH[i],
				NameEn: colorsEN[i],
			},
		}

		tx := DB.Create(&breed)
		if tx.Error != nil {
			fmt.Println("insert breed error:", tx.Error)
		}
	}
}

func testInsertAnmGender() {
	// 定义性别数据
	gendersZH := []string{"不明", "弟弟", "妹妹"}
	gendersEN := []string{"unknown", "boy", "gril"}

	for i := 0; i < len(gendersZH); i++ {
		anmGender := model.AnmGender{
			BriefModel: model.BriefModel{
				NameZh: gendersZH[i],
				NameEn: gendersEN[i],
			},
		}
		tx := DB.Create(&anmGender)
		if tx.Error != nil {
			fmt.Println("insert gender error:", tx.Error)
		}
	}
}

func testInsertStatus() {
	// 定义状态数据
	statusesZH := []string{"不明", "在校", "毕业", "退学", "喵星"}
	statusesEN := []string{"unknown", "inschool", "graduation", "missing", "catstar"}
	for i := 0; i < len(statusesZH); i++ {
		anmStatus := model.AnmStatus{
			BriefModel: model.BriefModel{
				NameZh: statusesZH[i],
				NameEn: statusesEN[i],
			},
		}

		tx := DB.Create(&anmStatus)
		if tx.Error != nil {
			fmt.Println("insertstatus error:", tx.Error)
		}
	}
}

func insertData() {
	testInsertSterilzation()
	fmt.Println("testInsertSterilzation success.")

	testInsertBreed()
	fmt.Println("testInsertBreed success.")

	testInsertStatus()
	fmt.Println("testInsertStatus success.")

	testInsertAnmGender()
	fmt.Println("testInsertAnmGender success.")
}

func main() {
	// 1.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "Havocantelope412#", "113.44.68.213", "3306", "hav_cats") // danger MySQL
	fmt.Println("dsn:", dsn)
	dbMySQL, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = dbMySQL

	autoMigrateTable()
	fmt.Println("autoMigrateTable over.")

	// insertData() // INFO 记得用完注释掉
}
