package model

func CreateEncounterFactory(sqlType string) *Encounter {
	return &Encounter{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type Encounter struct { // Encounter 或者称为 post，指的就是 Human 单次的记录。
	BaseModel
	// TAG 外键关联
	UserId     int
	UsersModel UsersModel
	AnimalsId  string `gorm:"size:20"` // TODO 关联对象存在上限

	Title   string `gorm:"size:20"`
	Content string
	// Time 从 CreatedAt 中解析

	// TAG Avatar 最好是压缩后的备份图像
	Avatar       string `gorm:"type:varchar(1s0)" json:"avatar,omitempty"` // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight uint16 `json:"avatar_height,omitempty"`                   // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	AvatarWidth  uint16 `json:"avatar_width,omitempty"`
	Photos       string `gorm:"type:varchar(100)" json:"photos,omitempty"` // 图片数组
	// POI
	Latitude  float64 `json:"latitude,omitempty"` // POI 位置相关
	Longitude float64 `json:"longitude,omitempty"`
	// TODO Comment Num 然后去统计？
}

func (e *Encounter) TableName() string {
	return "encounters"
}
