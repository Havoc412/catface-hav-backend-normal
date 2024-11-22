package model_res

import "catface/app/model"

type CatfaceCat struct {
	Animal model.Animal `json:"animal"`
	Conf   float64      `json:"conf"`
}
