package upload_file

import (
	"image"
	// "image/jpeg"
	"os"

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
	dstImg := imaging.Thumbnail(srcImg, targetWidth, targetHeight, imaging.Lanczos)
	// image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	// 使用高质量的滤波算法进行缩放
	// draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// Save
	err = imaging.Save(dstImg, dstPath)
	return
}
