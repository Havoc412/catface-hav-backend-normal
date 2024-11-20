package catface

import (
	"catface/app/global/variable"
	"catface/app/utils/micro_service"
	"context"
	"github.com/carlmjohnson/requests"
)

type FaceRes struct {
	FaceBreed int `json:"face_breed"`
	Cats      []struct {
		Id   int64   `json:"id"`
		Prob float64 `json:"prob"`
	} `json:"cats"`
}

func GetCatfaceResult(filePath string) FaceRes {
	body := map[string]interface{}{
		"file_path": filePath,
	}
	var res FaceRes
	err := requests.URL(micro_service.FetchPythonServiceUrl("cnn/detect_cat")).
		BodyJSON(&body).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		variable.ZapLog.Error("获取cat face结果集失败: " + err.Error())
	}

	return res
}
