package upload_file

import (
	"image"
	"image/jpeg"
	"os"
)

// ResizeImage 按照指定宽度等比例缩放图片
func ResizeImage(srcPath string, dstPath string, targetWidth int) error {
	// 打开源图片文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 解码源图片
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	// 获取源图片的尺寸
	bounds := srcImg.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// 计算目标高度
	targetHeight := int(float64(srcHeight) * (float64(targetWidth) / float64(srcWidth)))

	// 创建目标图片
	dstImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// 使用高质量的滤波算法进行缩放
	draw.CatmullRom.Scale(dstImg, dstImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, nil)

	// 打开目标图片文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 编码并保存目标图片
	err = jpeg.Encode(dstFile, dstImg, nil)
	if err != nil {
		return err
	}

	return nil
}
