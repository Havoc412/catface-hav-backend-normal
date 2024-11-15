package oss

import (
	"log"
	"os"
	"testing"
)

func TestDeleteFromQiNiu(t *testing.T) {
	err := DeleteFromQiNiu("img_1728355351")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func TestUploadToQiNiu(t *testing.T) {
	file, err := os.OpenFile("README.md", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err)
		return
	}

	niu, err := UploadToQiNiu(file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(niu)
	return
}

func TestMultiUploadToQiNiu(t *testing.T) {

	niu, err := MultiUploadToQiNiu(make([]string, 3))
	if err != nil {
		log.Fatalln(err)
		return
	}
	for _, s := range niu {
		log.Println(s)
	}

}
