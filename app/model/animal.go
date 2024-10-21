package model

import (
	"catface/app/global/variable"
	"catface/app/utils/gorm_v2"

	"go.uber.org/zap"
)

func CreateAnimalFactory(sqlType string) *Animal {
	return &Animal{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type Animal struct {
	BaseModel            // 假设 BaseModel 中不需要添加 omitempty 标签
	Name          string `gorm:"type:varchar(20)" json:"name,omitempty"`                            // 名称
	Birthday      string `json:"birthday,omitempty"`                                                // 生日
	Gender        uint8  `json:"gender,omitempty"`                                                  // 性别
	Breed         uint8  `json:"breed,omitempty"`                                                   // 品种
	Sterilization uint8  `json:"sterilization,omitempty"`                                           // 1 不明 2 未绝育 3 已绝育
	NickNames     string `gorm:"type:varchar(31)" json:"nick_names,omitempty"`                      // 别称，辅助查询；存储上采取 , 间隔符的方式; VARCHAR 会比较合适
	Status        uint8  `json:"status,omitempty"`                                                  // 状态
	Description   string `gorm:"column:description;type:varchar(255)" json:"description,omitempty"` // 简明介绍
	Tags          string `json:"tags,omitempty"`
	// TAG imaegs
	Avatar       string `gorm:"type:varchar(10)" json:"avatar,omitempty"`   // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight uint16 `json:"avatar_height,omitempty"`                    // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。  // INFO 暂时没用到
	AvatarWidth  uint16 `json:"avatar_width,omitempty"`                     // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	HeadImg      string `gorm:"type:varchar(10)" json:"head_img,omitempty"` // Head 默认处理为正方形。
	Photos       string `gorm:"type:varchar(100)" json:"photos,omitempty"`  // 图片数组
	// TAG POI
	Latitude       float64 `json:"latitude,omitempty"`        // POI 位置相关
	Longitude      float64 `json:"longitude,omitempty"`       // POI 位置相关
	ActivityRadius uint64  `json:"activity_radius,omitempty"` // 活动半径
	// CatFace
	FaceBreeds     string `json:"face_breeds,omitempty" gorm:"size:20"`
	FaceBreedProbs string `json:"face_breed_probs,omitempty" gorm:"size:20"`
}

func (a *Animal) TableName() string {
	return "animals"
}

func (a *Animal) Show(attrs []string, gender []uint8, breed []uint8, sterilzation []uint8, status []uint8, num int, skip int) (temp []Animal) {
	db := a.DB.Table(a.TableName()).Limit(int(num)).Offset(int(skip)).Select(attrs)

	// 创建条件映射
	conditions := map[string][]uint8{
		"gender":        gender,
		"breed":         breed,
		"sterilization": sterilzation,
		"status":        status,
	}

	db = gorm_v2.BuildWhere(db, conditions)

	err := db.Find(&temp).Error
	if err != nil {
		variable.ZapLog.Error("Animal Show Error", zap.Error(err))
	}
	return
}

func (a *Animal) ShowByID(id int) *Animal {
	var temp Animal
	err := a.DB.Table(a.TableName()).Model(&temp).Where("id = ?", id).Scan(&temp).Error
	if err != nil {
		variable.ZapLog.Error("Animal ShowByID Error", zap.Error(err))
	}
	return &temp
}

func (a *Animal) ShowByIDs(ids []int64, attrs ...string) (temp []Animal) {
	db := a.DB.Table(a.TableName())

	if len(attrs) > 0 {
		db = db.Select(attrs)
	}

	err := db.Where("id in (?)", ids).Find(&temp).Error
	if err != nil {
		variable.ZapLog.Error("Animal ShowByIDs Error", zap.Error(err))
	}
	return
}
