package test

import (
	"catface/app/model"
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
