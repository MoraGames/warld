package noiseMap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/MoraGames/warld/seed"
)

type WhiteNoiseMap struct {
	Map [][]float64
}

func NewWhiteNoiseMap(width, height int, seeder *seed.Seeder) WhiteNoiseMap {
	values := make([][]float64, height)
	for i := range values {
		values[i] = make([]float64, width)
		for j := range values[i] {
			values[i][j] = seeder.Random.Float64()
		}
	}
	return WhiteNoiseMap{Map: values}
}

func (wnm *WhiteNoiseMap) String() string {
	str := ""
	for _, row := range wnm.Map {
		for _, value := range row {
			str += "[" + fmt.Sprintf("%.2f", value) + "]"
		}
		str += "\n"
	}
	return str
}

func (wnm *WhiteNoiseMap) Image(name string) {
	width, height := len(wnm.Map[0]), len(wnm.Map)
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//fmt.Printf("[DEBUG] map[%v][%v] = %v", y, x, wnm.Map[y][x])

			// Convert to 0-255 scale
			grayColor := uint8(wnm.Map[y][x] * 255.0)
			img.SetGray(x, y, color.Gray{Y: grayColor})
		}
	}

	// Save the image to a file
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}
