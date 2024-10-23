// add_test.go
package test

import (
	"catface/app/model"
	"fmt"
	"testing"
)

func TestUsers(t *testing.T) {
	Init()

	user := model.UsersModel{}
	err := DB.AutoMigrate(&user)
	if err != nil {
		t.Error(err)
	}
}

func TestEncouner(t *testing.T) {
	Init()

	encounter := model.Encounter{}
	err := DB.AutoMigrate(&encounter)
	if err != nil {
		t.Error(err)
	}
}

func TestEncounterLike(t *testing.T) {
	Init()

	encounterLike := model.EncounterLike{}
	err := DB.AutoMigrate(&encounterLike)
	if err != nil {
		t.Error(err)
	}
}

func TestEncounterLevel_Insert(t *testing.T) {
	Init()

	encounterLevel := model.EncounerLevel{}
	err := DB.AutoMigrate(&encounterLevel)
	if err != nil {
		t.Error(err)
	}

	ZH := []string{"日常", "重大", "标志", "代办", "日程"}
	EN := []string{"daily", "serious", "flag", "todo", "schedule"}
	colorbg := []string{"#F0F0F0", "#FFD700", "#FF69B4", "#87CEFA", "#32CD32"}
	colorfont := []string{"#333333", "#000000", "#FFFFFF", "#000000", "#FFFFFF"}

	for i := 0; i < len(ZH); i++ {
		encounterLevel := model.EncounerLevel{
			BriefModel: model.BriefModel{
				NameZh: ZH[i],
				NameEn: EN[i],
			},
			Color: model.Color{
				ColorBackground: colorbg[i],
				ColorFont:       colorfont[i],
			},
		}
		DB.Create(&encounterLevel)
	}
}

func TestEncounterLike_Create_and_Delete(t *testing.T) {
	Init()

	encounterLike := model.EncounterLike{
		UsersModelId: 1,
		EncounterId:  15,
	}

	// DB.Create(&encounterLike)
	// DB.Delete(&encounterLike)

	// time.Sleep(2 * time.Second)

	// encounterLike = model.EncounterLike{
	// 	UsersModelId: 1,
	// 	EncounterId:  15,
	// }
	DB.Unscoped().Where("encounter_id", encounterLike.EncounterId).First(&encounterLike)
	fmt.Println(encounterLike.Id, encounterLike.IsDel)

	// INFO 这样操作无效
	res := DB.Where("id = ?", encounterLike.Id).Update("is_del", nil)
	fmt.Println(res.RowsAffected)

	// INFO 这样有效。
	encounterLike.IsDel = 0
	res2 := DB.Save(&encounterLike)
	fmt.Println(res2.RowsAffected)
}
