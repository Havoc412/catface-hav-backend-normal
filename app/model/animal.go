package model

import (
	"catface/app/global/variable"
	"catface/app/utils/data_bind"
	"catface/app/utils/gorm_v2"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateAnimalFactory(sqlType string) *Animal {
	return &Animal{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type Animal struct {
	// UPDATE 或者这里都应该采取外键连接？
	BaseModel            // 假设 BaseModel 中不需要添加 omitempty 标签
	Name          string `gorm:"type:varchar(20)" json:"name,omitempty"`                            // 名称
	Birthday      string `gorm:"size:10" json:"birthday,omitempty"`                                 // 生日；就简单存string就好
	Gender        uint8  `json:"gender,omitempty"`                                                  // 性别
	Breed         uint8  `json:"breed,omitempty"`                                                   // 品种
	Sterilization uint8  `json:"sterilization,omitempty"`                                           // 1 不明 2 未绝育 3 已绝育
	Vaccination   uint8  `json:"vaccination,omitempty"`                                             // 免疫状态
	Deworming     uint8  `json:"deworming,omitempty"`                                               // 驱虫状态
	NickNames     string `gorm:"type:varchar(31)" json:"nick_names,omitempty"`                      // 别称，辅助查询；存储上采取 , 间隔符的方式; VARCHAR 会比较合适
	Status        uint8  `json:"status,omitempty"`                                                  // 状态
	Description   string `gorm:"column:description;type:varchar(255)" json:"description,omitempty"` // 简明介绍
	Tags          string `json:"tags,omitempty"`
	// TAG imaegs
	Avatar       string `gorm:"type:varchar(50)" json:"avatar,omitempty"`   // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight uint16 `json:"avatar_height,omitempty"`                    // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。  // INFO 暂时没用到
	AvatarWidth  uint16 `json:"avatar_width,omitempty"`                     // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	HeadImg      string `gorm:"type:varchar(50)" json:"head_img,omitempty"` // Head 默认处理为正方形。
	Photos       string `gorm:"type:varchar(255)" json:"photos,omitempty"`  // 图片数组
	// TAG POI
	Department     uint8   `gorm:"column:department" json:"department,omitempty"`
	Latitude       float64 `json:"latitude,omitempty"`        // POI 位置相关
	Longitude      float64 `json:"longitude,omitempty"`       // POI 位置相关
	ActivityRadius uint64  `json:"activity_radius,omitempty"` // 活动半径
	// CatFace
	FaceModelScore float64 `json:"face_model_score,omitempty" gorm:"defalut:0"` // 评估面部模型得分
	FaceBreeds     string  `json:"face_breeds,omitempty" gorm:"size:20"`
	FaceBreedProbs string  `json:"face_breed_probs,omitempty" gorm:"size:20"`
	// 上传者 ID
	UsersModelId int64       `gorm:"column:user_id" json:"user_id,omitempty"` // 上传者 ID
	UsersModel   *UsersModel `json:"users_model,omitempty"`
}

func (a *Animal) TableName() string {
	return "animals"
}

func (a *Animal) Show(attrs []string, gender []uint8, breed []uint8, sterilization []uint8, status []uint8, department []uint8, notInIds []int64, num int, skip int) (temp []Animal) {
	db := a.DB.Table(a.TableName()).Limit(int(num)).Offset(int(skip)).Select(attrs)

	// 创建条件映射
	conditions := map[string][]uint8{
		"gender":        gender,
		"breed":         breed,
		"sterilization": sterilization,
		"status":        status,
		"department":    department,
	}

	db = gorm_v2.BuildWhere(db, conditions) // TIP 这里的 Where 条件连接就很方便了。

	if len(notInIds) > 0 {
		db = db.Where("id not in (?)", notInIds)
	}

	err := db.Find(&temp).Error
	if err != nil {
		variable.ZapLog.Error("Animal Show Error", zap.Error(err))
	}
	return
}

func (a *Animal) ShowByID(id int64) *Animal {
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

func (a *Animal) InsertDate(c *gin.Context) (int64, bool) {
	var tmp Animal
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		if res := a.Create(&tmp); res.Error == nil {
			// 获取插入的 ID
			insertedID := tmp.Id
			return insertedID, true
		} else {
			variable.ZapLog.Error("Animal 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("Animal 数据绑定出错", zap.Error(err))
	}
	return 0, false
}
