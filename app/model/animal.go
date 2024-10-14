package model

type Animal struct {
	BaseModel              // 假设 BaseModel 中不需要添加 omitempty 标签
	Name           string  `gorm:"type:varchar(20)" json:"name,omitempty"`                            // 名称
	Birthday       string  `json:"birthday,omitempty"`                                                // 生日
	Gender         uint8   `json:"gender,omitempty"`                                                  // 性别
	Breed          uint8   `json:"breed,omitempty"`                                                   // 品种
	Sterilization  uint8   `json:"sterilization,omitempty"`                                           // 1 不明 2 未绝育 3 已绝育
	NickName       string  `gorm:"type:varchar(31)" json:"nick_name,omitempty"`                       // 别称，辅助查询；存储上采取 , 间隔符的方式; VARCHAR 会比较合适
	Status         uint8   `json:"status,omitempty"`                                                  // 状态
	Description    string  `gorm:"column:description;type:varchar(255)" json:"description,omitempty"` // 简明介绍
	Avatar         string  `gorm:"type:varchar(10)" json:"avatar,omitempty"`                          // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight   uint16  `json:"avatar_height,omitempty"`                                           // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	AvatarWidth    uint16  `json:"avatar_width,omitempty"`                                            // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	HeadImg        string  `gorm:"type:varchar(10)" json:"head_img,omitempty"`                        // Head 默认处理为正方形。
	Photos         string  `gorm:"type:varchar(100)" json:"photos,omitempty"`                         // 图片数组
	Latitude       float64 `json:"latitude,omitempty"`                                                // POI 位置相关
	Longitude      float64 `json:"longitude,omitempty"`                                               // POI 位置相关
	ActivityRadius uint64  `json:"activity_radius,omitempty"`
	Tags           string  `json:"tags,omitempty"` // 活动半径
}

type Breed struct {
	BriefModel
}

type Sterilzation struct { // TEST How to use BriefModel, the dif between Common
	Id     int64  `json:"id"`
	NameZh string `json:"name_zh"`
	NameEn string `json:"name_en"`
}

type AnmStatus struct {
	BriefModel
}

type AnmGender struct {
	BriefModel
}

/**
 * @description: 保留 Top 3, 辅助 catface - breed 子模型判断； 单独建表，因为只会被 CatFace 模块使用。
 * @return {*}
 */
type AnmFaceBreed struct { // TODO 迁移 python 的时候再考虑一下细节
	BriefModel
	Top1  uint8
	Prob1 float64
	Top2  uint8
	Prob2 float64
	Top3  uint8
	Prob3 float64

	AnimalId int64 // INFO 外键设定?
	Animal   Animal
}
