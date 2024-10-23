package model

func CreateEncounterLikeFactory(sqlType string) *EncounterLike {
	return &EncounterLike{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type EncounterLike struct {
	BaseModel
	UsersModelId int `gorm:"column:user_id" json:"user_id"`
	UsersModel   UsersModel
	EncounterId  int `gorm:"column:encounter_id" json:"encounter_id"`
	Encounter    Encounter
	DeletedAt
}

func (e *EncounterLike) Create(userId, encounterId int) bool {
	e.UsersModelId = userId
	e.EncounterId = encounterId

	e.Unscoped().Where("encounter_id = ?", e.EncounterId).First(e)
	e.IsDel = 0 //
	if err := e.Save(e).Error; err != nil {
		return false
	}
	return true
}

func (e *EncounterLike) SoftDelete(userId, encounterId int) bool {
	e.UsersModelId = userId
	//
	e.Where("encounter_id = ?", encounterId).First(e)
	if err := e.Delete(e).Error; err != nil {
		return false
	}
	return true
}
