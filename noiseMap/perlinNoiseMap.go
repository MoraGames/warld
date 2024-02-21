package noiseMap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/MoraGames/warld/seed"
)

type (
	PerlinNoiseMap struct {
		Map [][]float64
	}

	randomVector struct {
		x float64
		y float64
	}
)

func NewPerlinNoiseMap(width, length, octaves int, gridSize float64, seeder *seed.Seeder) PerlinNoiseMap {
	values := make([][]float64, length)
	for i := range values {
		values[i] = make([]float64, width)
		for j := range values[i] {
			var (
				frequency float64 = 1.0
				amplitude float64 = 1.0
				value     float64 = 0.0
			)
			for k := 0; k < octaves; k++ {
				// perlinValue() returns a value in the range [-1, 1]
				value += perlinValue(float64(j)*frequency/gridSize, float64(i)*frequency/gridSize, seeder) * amplitude
				amplitude *= 0.5
				frequency *= 2
			}
			values[i][j] = value
		}
	}
	return PerlinNoiseMap{Map: values}
}

func perlinValue(x, y float64, seeder *seed.Seeder) float64 {
	xLeft, yUp := int(math.Floor(x)), int(math.Floor(y))
	xRight, yDown := xLeft+1, yUp+1

	xDiffLeft, yDiffUp := x-float64(xLeft), y-float64(yUp)

	dLU := dot(xLeft, yUp, x, y, seeder)
	dRU := dot(xRight, yUp, x, y, seeder)
	iUp := interpolate(dLU, dRU, xDiffLeft)

	dLD := dot(xLeft, yDown, x, y, seeder)
	dRD := dot(xRight, yDown, x, y, seeder)
	iDown := interpolate(dLD, dRD, xDiffLeft)

	return interpolate(iUp, iDown, yDiffUp)
}

func interpolate(v1, v2, w float64) float64 {
	return ((v2-v1)*(3.0-w*2.0)*w*w + v1)
}

func dot(ix, iy int, fx, fy float64, seeder *seed.Seeder) float64 {
	vec := newRandomVector(ix, iy, seeder)
	xDiff, yDiff := fx-float64(ix), fy-float64(iy)
	return vec.x*xDiff + vec.y*yDiff
}

// TODO: Understand what this function does D:
func newRandomVector(x, y int, seeder *seed.Seeder) randomVector {
	// No precomputed gradients mean this works for any number of grid coordinates
	const bitNum = 64                // bit number
	const rotationWidth = bitNum / 2 // rotation width
	x += math.MaxUint64 / 2
	y += math.MaxUint64 / 2
	uX, uY := uint64(x), uint64(y)
	uX *= 3284157443321321332
	uY ^= (uX << rotationWidth) | (uX >> (bitNum - rotationWidth))
	uY *= 1911520717346556890
	uinqueSeed := seeder.GetSeed().ToUint64()
	uinqueSeed ^= (uY << rotationWidth) | (uY >> (bitNum - rotationWidth))
	uinqueSeed *= 3937510949134324214
	uX ^= (uinqueSeed << rotationWidth) | (uinqueSeed >> (bitNum - rotationWidth))
	uX *= 2048419325134134313
	random := float64(uX) * (math.Pi / float64(^uint64(0)>>1)) // in [0, 2*Pi]
	return randomVector{x: math.Cos(random), y: math.Sin(random)}
}

func (pnm *PerlinNoiseMap) ScalePixel(z, x int, min, max float64) float64 {
	pnm.Map[z][x] = (((pnm.Map[z][x] + 1.0) / 2.0) * (max - min)) + min
	return pnm.Map[z][x]
}

func (pnm *PerlinNoiseMap) String() string {
	str := ""
	for _, row := range pnm.Map {
		for _, value := range row {
			str += "[" + fmt.Sprintf("%.2f", value) + "]"
		}
		str += "\n"
	}
	return str
}

func (pnm *PerlinNoiseMap) Image(name string) {
	width, height := len(pnm.Map[0]), len(pnm.Map)
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//fmt.Printf("[DEBUG] map[%v][%v] = %v", y, x, pnm.Map[y][x])

			// Convert to 0-255 scale
			grayColor := uint8(((1.0 - pnm.Map[y][x]) / 2.0) * 255.0)
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
