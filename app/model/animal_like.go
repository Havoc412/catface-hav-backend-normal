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

	a.Unscoped().Where("animal_id = ? AND user_id = ?", animalId, a.UsersModelId).First(a)
	return a.Delete(a).Error == nil
}

/**
 * @description: 查询是否存在关注记录
 * @param {*} userId
 * @param {int} animalId
 * @return {*}
 */
func (a *AnimalLike) Liked(userId, animalId int) bool {
	// 需要考虑 IsDel = 0;
	return a.Where("animal_id = ? AND user_id = ?", animalId, userId).First(a).RowsAffected > 0
}
