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

	for i := 0; i < len(ZH); i++ {
		encounterLevel := model.EncounerLevel{
			BriefModel: model.BriefModel{
				NameZh: ZH[i],
				NameEn: EN[i],
			},
		}
		DB.Create(&encounterLevel)
	}
}
