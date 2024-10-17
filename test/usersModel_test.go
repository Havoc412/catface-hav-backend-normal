// add_test.go
package test

import (
	"catface/app/model"
	"testing"
)

func TestUsers(t *testing.T) {
	Init()

	user := model.UsersModel{}
	err := DB.AutoMigrate(&user)
	if err != nil {
		t.Error(err)
	}
}
