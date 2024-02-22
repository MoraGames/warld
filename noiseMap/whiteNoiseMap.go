package noiseMap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/MoraGames/warld/coord"
	"github.com/MoraGames/warld/seed"
)

type WhiteNoiseMap struct {
	CurrentScale ScaleRange
	Map          [][]float64
}

func NewWhiteNoiseMap(width, height int, seeder *seed.Seeder) WhiteNoiseMap {
	values := make([][]float64, height)
	for i := range values {
		values[i] = make([]float64, width)
		for j := range values[i] {
			values[i][j] = seeder.Random.Float64()
		}
	}
	return WhiteNoiseMap{CurrentScale: ScaleRange{0, 1}, Map: values}
}

var WNMScalePixelDebugPrintTimes = 0

func (wnm *WhiteNoiseMap) ScalePixel(pixel coord.Coord, from, to ScaleRange) {
	if WNMScalePixelDebugPrintTimes < 3 {
		fmt.Printf("[DEBUG] > > > Scaling pixel [%v][%v] from %+v to %+v\n", pixel.Z, pixel.X, from, to)
		fmt.Printf("[DEBUG] > > > (((pixel + %v) / %v) * %v) + %v\n", from.Min, from.Max-from.Min, to.Max-to.Min, to.Min)
		WNMScalePixelDebugPrintTimes++
	}
	wnm.Map[pixel.Z][pixel.X] = wnm.ScaleCopyPixel(pixel, from, to)
}

func (wnm *WhiteNoiseMap) ScaleCopyPixel(pixel coord.Coord, from, to ScaleRange) float64 {
	return (((wnm.Map[pixel.Z][pixel.X] - from.Min) / (from.Max - from.Min)) * (to.Max - to.Min)) + to.Min
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

func (wnm *WhiteNoiseMap) Image(path string) {
	width, height := len(wnm.Map[0]), len(wnm.Map)
	img := image.NewGray(image.Rect(0, 0, width, height))
	fmt.Println("[DEBUG] > > > > Image base created")
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//fmt.Printf("[DEBUG] map[%v][%v] = %v", y, x, wnm.Map[y][x])

			// Convert to 0-255 scale
			grayColor := uint8(wnm.ScaleCopyPixel(coord.Coord{X: x, Z: y}, wnm.CurrentScale, ScaleRange{0, 255}))
			img.SetGray(x, y, color.Gray{Y: grayColor})
		}
	}
	fmt.Println("[DEBUG] > > > > Image pixels setted")

	// Save the image to a file
	file, err := os.Create(fmt.Sprintf(path, PathFilesNumber))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
	fmt.Println("[DEBUG] > > > > Image saved")
}
