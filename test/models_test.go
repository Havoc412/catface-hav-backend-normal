// add_test.go
package test

import (
	"catface/app/model"
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

func TestEncounterLevel(t *testing.T) {
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
