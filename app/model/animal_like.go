package model

func CreateAnimalLikeFactory(sqlType string) *AnimalLike {
	return &AnimalLike{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type AnimalLike struct {
	BaseModel
	UsersModelId int `gorm:"column:user_id" json:"user_id"`
	UsersModel   UsersModel
	AnimalId     int `gorm:"column:animal_id" json:"animal_id"`
	Animal       Animal
	DeletedAt
}

func (a *AnimalLike) Create(userId, animalId int) bool {
	a.UsersModelId = userId
	a.AnimalId = animalId

	a.Unscoped().Where("animal_id = ? AND user_id = ?", a.AnimalId, a.UsersModelId).First(a)
	a.IsDel = 0
	return a.Save(a).Error == nil
}

func (a *AnimalLike) SoftDelete(userId, animalId int) bool {
	a.UsersModelId = userId

	a.Unscoped().Where("animal_id = ? AND user_id = ?", a.AnimalId, a.UsersModelId).First(a)
	return a.Delete(a).Error == nil
}
