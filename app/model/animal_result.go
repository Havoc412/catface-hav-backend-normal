package model

type AnimalWithLikeList struct {
	Animal Animal `json:"animal"`
	Like   bool   `json:"like,omitempty"`
}
