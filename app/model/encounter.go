package model

import (
	"catface/app/global/variable"
	"catface/app/utils/data_bind"
	"log"
	"strings"

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
	UsersModel   *UsersModel `json:"users_model,omitempty"` // INFO 由于 Detail 返回空子段有些麻烦，先尝试采用指针。
	// AnimalsId    string      `gorm:"size:20" json:"animals_id"` // 关联对象存在上限  // INFO 还是采取分表，方便查询。

	Title    string   `gorm:"size:20;column:title" json:"title"`
	Content  string   `json:"content"`
	Level    uint8    `json:"level" gorm:"column:level;default:1"`
	Tags     string   `json:"tags,omitempty" gorm:"column:tags;size:50"`
	TagsList []string `gorm:"-" json:"tags_list,omitempty"`

	// TAG Avatar 最好是压缩后的备份图像
	Avatar       string   `gorm:"type:varchar(50)" json:"avatar,omitempty"` // 缩略图 url，为 Go 获取 Photo 之后压缩处理后的图像，单独存储。
	AvatarHeight uint16   `json:"avatar_height,omitempty"`                  // 为了方便前端在加载图像前的骨架图 & 瀑布流展示。
	AvatarWidth  uint16   `json:"avatar_width,omitempty"`
	Photos       string   `gorm:"type:varchar(100)" json:"photos,omitempty"` // 图片数组
	PhotosList   []string `gorm:"-" json:"photos_list,omitempty"`            // TIP GORM 忽略
	// POI
	Latitude  float64 `json:"latitude,omitempty"` // POI 位置相关
	Longitude float64 `json:"longitude,omitempty"`
	// TODO Comment Num 然后去统计？
}

func (e *Encounter) TableName() string {
	return "encounters"
}

/**
 * @description:
 * @param {*gin.Context} c
 * @return {*} 返回创建的绑定对象，之后 model_es 的利用。
 */
func (e *Encounter) InsertDate(c *gin.Context) (tmp Encounter, ok bool) {
	if err := data_bind.ShouldBindFormDataToModel(c, &tmp); err == nil {
		if res := e.Create(&tmp); res.Error == nil {
			return tmp, true
		} else {
			variable.ZapLog.Error("Encounter 数据新增出错", zap.Error(res.Error))
		}
	} else {
		variable.ZapLog.Error("Encounter 数据绑定出错", zap.Error(err))
	}
	return tmp, false
}

func formatEncounterList(rows *gorm.DB) (temp []EncounterList, err error) {
	// 获取底层的 sql.Rows 对象
	sqlRows, err := rows.Rows()
	if err != nil {
		log.Println("获取 sql.Rows 失败:", err)
		return nil, err
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
		// STAGE: Set link Status
		item.Like = ueLikeInt == 1
		// STAGE: Check Url exist
		if item.Avatar == "" {
			if animal_id, ok := CreateEncounterAnimalLinkFactory("").ShowByEncounterIdFirst(int64(item.Id)); ok {
				if animal := CreateAnimalFactory("").ShowByID(animal_id); animal != nil {
					item.Avatar = animal.Avatar
					item.AvatarHeight = int(animal.AvatarHeight)
					item.AvatarWidth = int(animal.AvatarWidth)
					item.UseAnimalAvatar = true
				}
			}
		}
		// append
		temp = append(temp, item)
	}

	return
}

func (e *Encounter) Show(num, skip, user_id int, animals_id []int) (temp []EncounterList) {
	// STAGE - 1：build SQL
	var sqlBuilder strings.Builder

	// 构建基础查询
	sqlBuilder.WriteString(`
SELECT e.id, e.user_id, title, avatar, avatar_height, avatar_width, e.updated_at, user_name, user_avatar,
	EXISTS (
		SELECT 1
		FROM encounter_likes l
		WHERE l.user_id = ? AND l.encounter_id = e.id AND l.is_del = 0
	) AS ue_like
FROM encounters e 
JOIN tb_users u ON e.user_id = u.id`)

	// 动态插入 animals_id 条件
	if len(animals_id) > 0 {
		sqlBuilder.WriteString(`
JOIN encounter_animal_links eal ON e.id = eal.encounter_id
WHERE eal.animal_id IN (?)`)
	}

	// 添加排序和分页
	sqlBuilder.WriteString(`
ORDER BY e.updated_at DESC
LIMIT ? OFFSET ?
	`)

	sql := sqlBuilder.String() // 获取到 SQL；

	// STAGE - 2: Exe SQL
	var rows *gorm.DB
	if len(animals_id) > 0 {
		rows = e.Raw(sql, user_id, animals_id, num, skip)
	} else {
		rows = e.Raw(sql, user_id, num, skip)
	}
	if rows.Error != nil {
		log.Println("查询失败:", rows.Error)
		return nil
	}
	var err error
	if temp, err = formatEncounterList(rows); err != nil {
		log.Println("格式化数据失败:", err)
		return nil
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

func (e *Encounter) ShowByIDs(ids []int64, attrs ...string) (temp []Encounter) {
	db := e.DB.Table(e.TableName())

	if len(attrs) > 0 {
		db = db.Select(attrs)
	}

	err := db.Where("id in (?)", ids).Find(&temp).Error
	if err != nil {
		variable.ZapLog.Error("Encounter ShowByIDs Error", zap.Error(err))
	}
	return
}

/**
 * @description: 过去 1 个月，发送过路遇表的 ids，同时去重。
 * @param {*} user_id
 * @param {int} num
 * @return {*}
 */
func (e *Encounter) EncounteredCats(user_id, num int) ([]int64, error) {
	sql := `SELECT eal.animal_id 
            FROM encounter_animal_links eal
            JOIN encounters e 
                ON e.id = eal.encounter_id AND e.user_id = ?
			WHERE e.updated_at >= DATE_SUB(CURDATE(), INTERVAL 1 MONTH)
			ORDER BY e.updated_at DESC
			LIMIT ?`

	rows, err := e.Raw(sql, user_id, num).Rows()
	if err != nil {
		log.Println("查询失败:", err)
		return nil, err
	}
	defer rows.Close()

	// Scan 同时去重。
	var temp []int64
	seen := make(map[int64]bool)

	for rows.Next() {
		var animal_id int64
		if err := rows.Scan(&animal_id); err != nil {
			log.Println("扫描失败:", err)
			return nil, err
		}
		if !seen[animal_id] {
			seen[animal_id] = true
			temp = append(temp, animal_id)
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("遍历失败:", err)
		return nil, err
	}

	return temp, nil
}
