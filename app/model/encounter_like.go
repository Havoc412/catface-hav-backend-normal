package model

type EncounterLike struct {
	BaseModel
	UsersModelId int `gorm:"column:user_id" json:"user_id"`
	UsersModel   UsersModel
	EncounterId  int `gorm:"column:encounter_id" json:"encounter_id"`
	Encounter    Encounter
}
