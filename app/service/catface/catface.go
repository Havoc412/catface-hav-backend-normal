package catface

import (
	"catface/app/global/variable"
	"catface/app/utils/micro_service"
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

type FaceData struct {
	FaceBreed string `json:"face_breed"`
	Cats      []struct {
		Id   int64   `json:"id"`
		Conf float64 `json:"conf"`
	} `json:"cats"`
}

type CatFaceRes struct {
	Status int      `json:"status"`
	Data   FaceData `json:"data"`
}

func GetCatfaceResult(filePath string) (*FaceData, error) {
	body := map[string]interface{}{
		"file_path": filePath,
	}
	var res CatFaceRes
	err := requests.URL(micro_service.FetchPythonServiceUrl("cnn/detect_cat")).
		BodyJSON(&body).
		ToJSON(&res).
		Fetch(context.Background())

	if err != nil {
		variable.ZapLog.Error("获取cat face结果集失败: " + err.Error())
		return nil, err
	}

	if res.Status != 200 {
		return nil, fmt.Errorf("请求状态码错误: %d", res.Status)
	}

	return &res.Data, nil
}
