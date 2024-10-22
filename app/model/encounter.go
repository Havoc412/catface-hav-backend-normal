package model

import (
	"catface/app/global/variable"
	"catface/app/utils/data_bind"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateEncounterFactory(sqlType string) *Encounter {
	return &Encounter{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type Encounter struct { // Encounter 或者称为 post，指的就是 Human 单次的记录。
	BaseModel
	// TAG 外键关联
	UsersModelId int64       `gorm:"column:user_id" json:"user_id"`
	UsersModel   *UsersModel `json:"users_model,omitempty"`     // INFO 由于 Detail 返回空子段有些麻烦，先尝试采用指针。
	AnimalsId    string      `gorm:"size:20" json:"animals_id"` // TODO 关联对象存在上限

	Title   string `gorm:"size:20;column:title" json:"title"`
	Content string `json:"content"`
	// Time 从 CreatedAt 中解析

	// TAG Avatar 最好是压缩后的备份图像
	Avatar       string   `gorm:"type:varchar(50)" json:"avatar,omitempty"` // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight uint16   `json:"avatar_height,omitempty"`                  // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	AvatarWidth  uint16   `json:"avatar_width,omitempty"`
	Photos       string   `gorm:"type:varchar(100)" json:"photos,omitempty"` // 图片数组
	PhotosSlice  []string `gorm:"-" json:"photos_list,omitempty"`            // TIP GORM 忽略
	// POI
	Latitude  float64 `json:"latitude,omitempty"` // POI 位置相关
	Longitude float64 `json:"longitude,omitempty"`
	// TODO Comment Num 然后去统计？
}

func (e *Encounter) TableName() string {
	return "encounters"
}

func (e *Encounter) InsertDate(c *gin.Context) bool {
	var tmp Encounter
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		if res := e.Create(&tmp); res.Error == nil {
			return true
		} else {
			variable.ZapLog.Error("Encounter 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("Encounter 数据绑定出错", zap.Error(err))
	}
	return false
}

func (e *Encounter) Show(num, skip, user_id int) (temp []EncounterList) {
	sql := `
		SELECT e.id, e.user_id, title, avatar, avatar_height, avatar_width, e.updated_at, user_name, user_avatar,
			EXISTS (
				SELECT 1
				FROM encounter_likes l
				WHERE l.user_id = ? AND l.encounter_id = e.id
			) AS ue_like
		FROM encounters e 
		JOIN tb_users u ON e.user_id = u.id
		LIMIT ? OFFSET ?
	`
	// err := e.Raw(sql, user_id, num, skip).Scan(&temp).Error
	// fmt.Println(err)

	var rows *gorm.DB
	if rows = e.Raw(sql, user_id, num, skip); rows.Error != nil {
		log.Println("查询失败:", rows.Error)
		return nil
	}
	// 获取底层的 sql.Rows 对象
	sqlRows, err := rows.Rows()
	if err != nil {
		log.Println("获取 sql.Rows 失败:", err)
		return nil
	}
	defer sqlRows.Close()

	for sqlRows.Next() {
		var item EncounterList
		var ueLikeInt int
		dest := []interface{}{
			&item.Id, &item.UserId, &item.Title, &item.Avatar, &item.AvatarHeight, &item.AvatarWidth, &item.UpdatedAt, &item.UserName, &item.UserAvatar, &ueLikeInt,
		}
		if err := sqlRows.Scan(dest...); err != nil {
			log.Println("扫描失败:", err)
			continue
		}
		item.Like = ueLikeInt == 1
		temp = append(temp, item)
	}

	return
}

func (e *Encounter) ShowByID(id int64) (temp *Encounter, err error) {
	// 1. search encounter
	if err = e.Where("id = ?", id).First(&temp).Error; err != nil {
		return
	}
	return
	// // 2. search user data
	// user := UsersModel{BaseModel: BaseModel{Id: encounter.UsersModelId}}
	// if err := user.Select("user_name", "user_avatar").First(&user).Error; err != nil {
	// 	return
	// }

	// // 3. search animals data
	// animals_id := query_handler.StringToint64Array(encounter.AnimalsId)
	// var animals []Animal
	// if err := e.Model(&animals).Select("id", "avatar", "name").Where("id in (?)", animals_id).Find(&animals).Error; err != nil {
	// 	return
	// }

	// // TODO 4. 然后整合
	// return
}
