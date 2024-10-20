package model

import "time"

type EncounterList struct {
	UserId       int `form:"user_id" json:"user_id"`
	Title        string
	Avatar       string
	AvatarHeight int        `form:"avatar_height" json:"avatar_height"`
	AvatarWidth  int        `form:"avatar_width" json:"avatar_width"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"` // TIP 设为 *time.Time，omitempty 和 autoUpdated 就都可以生效

	UserName string `form:"user_name" json:"user_name"`
}
