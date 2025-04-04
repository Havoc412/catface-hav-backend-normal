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

	ZH := []string{"日常", "重要", "标志", "代办", "日程"}
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

func TestAnimalLike(t *testing.T) {
	Init()
	animalLike := model.AnimalLike{}
	err := DB.AutoMigrate(&animalLike)
	if err != nil {
		t.Error(err)
	}
}

func TestEaLink(t *testing.T) {
	Init()
	eaLink := model.EncounterAnimalLink{}
	err := DB.AutoMigrate(&eaLink)
	if err != nil {
		t.Error(err)
	}
}

// // 测试函数
// func TestInsertEncounterAnimalLinks(t *testing.T) {
// 	Init()

// 	// 定义要插入的数据
// 	data := []struct {
// 		EncounterId int
// 		AnimalIds   string
// 	}{
// 		{10, "4"},
// 		{11, "2,3"},
// 		{13, "4"},
// 		{14, "4"},
// 		{15, "4"},
// 		{16, "4"},
// 		{17, "4"},
// 		{18, "4"},
// 		{19, "4"},
// 		{20, "4"},
// 	}

// 	// 插入数据
// 	for _, item := range data {
// 		animalIds := strings.Split(item.AnimalIds, ",")
// 		for _, animalIdStr := range animalIds {
// 			animalId, err := strconv.Atoi(animalIdStr)
// 			if err != nil {
// 				t.Errorf("Failed to convert animal Id: %v", err)
// 				continue
// 			}
// 			link := model.EncounterAnimalLink{
// 				EncounterId: item.EncounterId,
// 				AnimalId:    animalId,
// 			}
// 			if err := DB.Create(&link).Error; err != nil {
// 				t.Errorf("Failed to insert link: %v", err)
// 			}
// 		}
// 	}

// 	// 验证数据是否正确插入
// 	var links []model.EncounterAnimalLink
// 	if err := DB.Find(&links).Error; err != nil {
// 		t.Errorf("Failed to fetch links: %v", err)
// 	}

// 	expectedLinks := []model.EncounterAnimalLink{
// 		{EncounterId: 10, AnimalId: 4},
// 		{EncounterId: 11, AnimalId: 2},
// 		{EncounterId: 11, AnimalId: 3},
// 		{EncounterId: 13, AnimalId: 4},
// 		{EncounterId: 14, AnimalId: 4},
// 		{EncounterId: 15, AnimalId: 4},
// 		{EncounterId: 16, AnimalId: 4},
// 		{EncounterId: 17, AnimalId: 4},
// 		{EncounterId: 18, AnimalId: 4},
// 		{EncounterId: 19, AnimalId: 4},
// 		{EncounterId: 20, AnimalId: 4},
// 	}

// 	if !reflect.DeepEqual(links, expectedLinks) {
// 		t.Errorf("Expected links: %v, but got: %v", expectedLinks, links)
// 	}
// }

// TestDepartment_Insert 测试 Department 模型的自动迁移和数据插入
func TestDepartment_Insert(t *testing.T) {
	Init()

	// 自动迁移 Department 表
	if err := DB.AutoMigrate(&model.Department{}); err != nil {
		t.Errorf("Failed to auto migrate Department table: %v", err)
	}

	// 定义要插入的部门名称
	departmentsZh := []string{"其他", "信部", "医学部", "工部", "文理", "湖滨", "枫园", "网安"}
	departmentsEn := []string{"other", "it", "medical", "engine", "art", "lake", "maple", "net"}

	// 插入数据
	for i := 0; i < len(departmentsZh); i++ {
		department := model.Department{
			BriefModel: model.BriefModel{
				NameZh: departmentsZh[i],
				NameEn: departmentsEn[i],
			},
		}
		if err := DB.Create(&department).Error; err != nil {
			t.Errorf("Failed to insert data into Department table: %v", err)
		}
	}

	// 查询数据并检查
	var result []model.Department
	if err := DB.Find(&result).Error; err != nil {
		t.Errorf("Failed to query data from Department table: %v", err)
	}

	if len(result) != len(departmentsZh) {
		t.Errorf("Expected %d departments, got %d", len(departmentsZh), len(result))
	}

	for i, dept := range result {
		if dept.NameZh != departmentsZh[i] || dept.NameEn != departmentsEn[i] {
			t.Errorf("Expected department name '%s' (zh) and '%s' (en), got '%s' (zh) and '%s' (en)", departmentsZh[i], departmentsEn[i], dept.NameZh, dept.NameEn)
		}
	}
}
