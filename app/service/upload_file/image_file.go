package upload_file

import (
	"catface/app/global/my_errors"
	"catface/app/global/variable"
	"catface/app/utils/md5_encrypt"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

/**
 * @description:  ResizeImage 按照指定宽度等比例缩放图片
 * @param {string} srcPath 需要完整路径
 * @param {string} dstPath
 * @param {int} targetWidth
 * @return {*}
 */
func ResizeImage(srcPath string, dstPath string, targetWidth int) (targetHeight int, err error) {
	// 打开源图片文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// 解码源图片
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return
	}

	// 获取源图片的尺寸
	bounds := srcImg.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// 计算目标高度
	targetHeight = int(float64(srcHeight) * (float64(targetWidth) / float64(srcWidth)))

	// 创建目标图片
	dstImg := imaging.Resize(srcImg, targetWidth, targetHeight, imaging.Lanczos)
	// image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	// 使用高质量的滤波算法进行缩放
	// draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// Save
	// 相关路径不存在，创建目录
	dstFolderPath := filepath.Dir(dstPath)
	if _, err = os.Stat(dstFolderPath); err != nil {
		if err = os.MkdirAll(dstFolderPath, os.ModePerm); err != nil {
			variable.ZapLog.Error("文件上传创建目录出错" + err.Error())
			return
		}
	}
	err = imaging.Save(dstImg, dstPath)
	return
}

func DownloadImage(imageUrl, dstPath string) error {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("图片下载失败，状态码: %d", resp.StatusCode)
	}

	if sequence := variable.SnowFlake.GetId(); sequence > 0 {
		saveFileName := fmt.Sprintf("%d%s", sequence, filepath.Base(imageUrl))
		saveFileName = md5_encrypt.MD5(saveFileName) + ".jpg"

		fullSavePath := filepath.Join(dstPath, saveFileName)
		file, err := os.Create(fullSavePath)
		if err != nil {
			variable.ZapLog.Error("文件保存出错：" + err.Error())
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			variable.ZapLog.Error("文件写入出错：" + err.Error())
			return err
		}
	} else {
		err := errors.New(my_errors.ErrorsSnowflakeGetIdFail)
		variable.ZapLog.Error("文件保存出错：" + err.Error())
		return err
	}
	return nil
}
