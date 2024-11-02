package model

func CreateAnimalLikeFactory(sqlType string) *AnimalLike {
	return &AnimalLike{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type AnimalLike struct {
	BaseModel
	UsersModelId int `gorm:"column:user_id;index" json:"user_id"`
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

func (a *AnimalLike) LikedBatch(userId int, animalIds []int) ([]bool, error) {
	var results []AnimalLike
	db := a.Where("user_id = ? AND animal_id IN (?) AND is_del = 0", userId, animalIds).Find(&results)

	if db.Error != nil {
		return nil, db.Error
	}

	// 创建一个布尔值数组，初始化为 false
	likedMap := make(map[int]bool)
	for _, result := range results {
		likedMap[result.AnimalId] = true
	}

	// 构建返回结果
	var likedResults []bool
	for _, animalId := range animalIds {
		likedResults = append(likedResults, likedMap[animalId])
	}

	return likedResults, nil
}

/**
 * @description: 获取所有关注的猫咪id
 * @param {int} userId
 * @return {*}
 */
func (a *AnimalLike) LikedCats(userId int) (res []int) {
	var results []AnimalLike
	db := a.Select("animal_id").Where("user_id = ?", userId).Find(&results)
	if db.Error != nil {
		return nil
	}

	for _, result := range results {
		res = append(res, result.AnimalId)
	}

	return
}
