package utils

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

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
