package upload_file

import (
	"catface/app/global/my_errors"
	"catface/app/global/variable"
	"catface/app/utils/md5_encrypt"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(context *gin.Context, savePath string) (r bool, finnalSavePath interface{}) {
	newSavePath, newReturnPath := generateYearMonthPath(savePath)

	// time.Sleep(2 * time.Second) // TEST 模拟服务器的访问延迟。

	//  1.获取上传的文件名(参数验证器已经验证完成了第一步错误，这里简化)
	file, _ := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField")) //  file 是一个文件结构体（文件对象）

	//  保存文件，原始文件名进行全局唯一编码加密、md5 加密，保证在后台存储不重复
	var saveErr error
	if sequence := variable.SnowFlake.GetId(); sequence > 0 {
		saveFileName := fmt.Sprintf("%d%s", sequence, file.Filename)
		saveFileName = md5_encrypt.MD5(saveFileName) + path.Ext(saveFileName)

		if saveErr = context.SaveUploadedFile(file, filepath.Join(newSavePath, saveFileName)); saveErr == nil {
			//  上传成功,返回资源的相对路径，这里请根据实际返回绝对路径或者相对路径
			finnalSavePath = gin.H{
				"path": strings.ReplaceAll(filepath.Join(newReturnPath, saveFileName), variable.BasePath, ""),
			}
			return true, finnalSavePath
		}
	} else {
		saveErr = errors.New(my_errors.ErrorsSnowflakeGetIdFail)
		variable.ZapLog.Error("文件保存出错：" + saveErr.Error())
	}
	return false, nil
}

// 文件上传可以设置按照 xxx年-xx月 格式存储
// INFO 但这个 returnPath 我还基本没有用到。
func generateYearMonthPath(savePathPre string) (string, string) {
	returnPath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileReturnPath") // UPDATE 因为没用到，所以就先不调整了。
	curYearMonth := time.Now().In(time.Local).Format("2006_01")
	newSavePathPre := filepath.Join(savePathPre, curYearMonth)
	newReturnPathPre := filepath.Join(returnPath, curYearMonth)
	// 相关路径不存在，创建目录
	if _, err := os.Stat(newSavePathPre); err != nil {
		if err = os.MkdirAll(newSavePathPre, os.ModePerm); err != nil {
			variable.ZapLog.Error("文件上传创建目录出错" + err.Error())
			return "", ""
		}
	}
	return newSavePathPre + "/", newReturnPathPre + "/"
}
