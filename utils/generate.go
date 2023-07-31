package utils

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"

	"github.com/mowshon/moviego"
	"github.com/nullrocks/identicon"
	"github.com/szeliga/goray/engine"
)

// 根据用户id生成随机的背景图
func GenerateBackground(id uint) {
	width := 2000
	height := 800
	random_num1 := rand.Intn(400)
	random_num2 := rand.Intn(7)
	n := math.Sqrt(float64(random_num1))*10 + float64(random_num2*random_num2)
	scene := engine.NewScene(width, height)
	scene.EachPixel(func(x, y int) color.RGBA {
		return color.RGBA{
			uint8(x * 255 / width),
			uint8(y * 255 / height),
			uint8(n),
			255,
		}
	})
	scene.Save(fmt.Sprintf("./public/background/%d.png", id))
}

// 根据用户名生成头像（类似于github默认随机头像的那种）
func GenerateAvatar(username string, id uint) {
	gen, _ := identicon.New(
		"github",
		5,
		3,
	)
	img, _ := gen.Draw(username)
	//新建文件
	file, _ := os.Create(fmt.Sprintf("./public/avatar/%d.png", id))
	defer file.Close()
	//写入图片
	img.Png(300, file)
}

// 为视频生成封面
func GenerateCover(id uint, p string) {
	video, _ := moviego.Load("./" + p)
	video.Screenshot(0, fmt.Sprintf("./public/screenshot/%d.png", id))
}
