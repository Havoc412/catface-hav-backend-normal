package model

type AnimalWithLikeList struct {
	Animal Animal `json:"animal"`
	Like   bool   `json:"like,omitempty"`
}

type AnimalWithNickNameHit struct {
	Animal      Animal `json:"animal"`
	NickNameHit bool   `json:"nick_name_hit,omitempty"`
}
