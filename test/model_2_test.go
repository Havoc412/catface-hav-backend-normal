package test

import (
	"catface/app/model"
	_ "catface/bootstrap"
	"testing"
)

func TestDocModel(t *testing.T) {
	Init()

	doc := model.Doc{}
	err := DB.AutoMigrate(&doc)
	if err != nil {
		t.Error(err)
	}
}

func TestAnmBreed(t *testing.T) {
	res := model.CreateAnmBreedFactory("").GetNameZhByEn("li")
	t.Log(res)
}
