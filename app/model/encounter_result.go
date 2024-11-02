package model

import "time"

// TIP 用这种方式绑定 GORM 的 Raw 就会比较有效了。
type EncounterList struct {
	Id           int        `form:"id" json:"id"`
	UserId       int        `form:"user_id" json:"user_id"`
	Title        string     `form:"title" json:"title"`
	Avatar       string     `form:"avatar" json:"url"`
	AvatarHeight int        `form:"avatar_height" json:"height"`
	AvatarWidth  int        `form:"avatar_width" json:"width"`
	UpdatedAt    *time.Time `form:"updated_at" json:"time"` // TIP 设为 *time.Time，omitempty 和 autoUpdated 就都可以生效

	Like            bool `form:"ue_like" json:"like"`
	UseAnimalAvatar bool `form:"use_animal_avatar" json:"use_animal_avatar"`

	UserName   string `form:"user_name" json:"user_name"`
	UserAvatar string `form:"user_avatar" json:"user_avatar"`
}

type EncounterDetail struct {
	Encounter  Encounter  `json:"encounter"`
	UsersModel UsersModel `json:"user"`
	Animals    []Animal   `json:"animals"`
}
