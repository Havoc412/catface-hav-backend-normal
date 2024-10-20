package upload_file

import (
	"catface/app/global/variable"
	"image"
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
