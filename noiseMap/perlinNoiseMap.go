package noiseMap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/MoraGames/warld/coord"
	"github.com/MoraGames/warld/seed"
)

type (
	PerlinNoiseMap struct {
		CurrentScale ScaleRange
		Map          [][]float64
	}

	randomVector struct {
		x float64
		y float64
	}
)

func NewPerlinNoiseMap(width, length int, frequency, amplitude, lacunarity, persistence float64, octaves int, seeder *seed.Seeder) PerlinNoiseMap {
	noisemapSeed := seeder.Random.Uint64()
	values := make([][]float64, length)
	min, max := math.MaxFloat64, math.SmallestNonzeroFloat64
	for i := range values {
		values[i] = make([]float64, width)
		for j := range values[i] {
			var (
				freq          = frequency
				ampl          = amplitude
				value float64 = 0.0
			)
			for k := 0; k < octaves; k++ {
				// perlinValue() returns a value in the range [-1, 1]
				value += perlinValue(float64(j)*freq, float64(i)*freq, noisemapSeed) * float64(ampl)
				freq *= lacunarity
				ampl *= persistence
			}
			values[i][j] = value

			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}
	return PerlinNoiseMap{CurrentScale: ScaleRange{min, max}, Map: values}
}

func perlinValue(x, y float64, noisemapSeed uint64) float64 {
	xLeft, yUp := int(math.Floor(x)), int(math.Floor(y))
	xRight, yDown := xLeft+1, yUp+1

	xDiffLeft, yDiffUp := x-float64(xLeft), y-float64(yUp)

	dLU := dot(xLeft, yUp, x, y, noisemapSeed)
	dRU := dot(xRight, yUp, x, y, noisemapSeed)
	iUp := interpolate(dLU, dRU, xDiffLeft)

	dLD := dot(xLeft, yDown, x, y, noisemapSeed)
	dRD := dot(xRight, yDown, x, y, noisemapSeed)
	iDown := interpolate(dLD, dRD, xDiffLeft)

	return interpolate(iUp, iDown, yDiffUp)
}

func interpolate(v1, v2, w float64) float64 {
	return ((v2-v1)*(3.0-w*2.0)*w*w + v1)
}

func dot(ix, iy int, fx, fy float64, noisemapSeed uint64) float64 {
	vec := newRandomVector(ix, iy, noisemapSeed)
	xDiff, yDiff := fx-float64(ix), fy-float64(iy)
	return vec.x*xDiff + vec.y*yDiff
}

// TODO: Understand what this function does D:
func newRandomVector(x, y int, noisemapSeed uint64) randomVector {
	// No precomputed gradients mean this works for any number of grid coordinates
	const bitNum = 64                // bit number
	const rotationWidth = bitNum / 2 // rotation width
	x += math.MaxUint64 / 2
	y += math.MaxUint64 / 2
	uX, uY := uint64(x), uint64(y)
	uX *= 3284157443321321332
	uY ^= (uX << rotationWidth) | (uX >> (bitNum - rotationWidth))
	uY *= 1911520717346556890
	uinqueSeed := noisemapSeed
	uinqueSeed ^= (uY << rotationWidth) | (uY >> (bitNum - rotationWidth))
	uinqueSeed *= 3937510949134324214
	uX ^= (uinqueSeed << rotationWidth) | (uinqueSeed >> (bitNum - rotationWidth))
	uX *= 2048419325134134313
	random := float64(uX) * (math.Pi / float64(^uint64(0)>>1)) // in [0, 2*Pi]
	return randomVector{x: math.Cos(random), y: math.Sin(random)}
}

func (pnm *PerlinNoiseMap) ScalePixel(pixel coord.Coord, from, to ScaleRange) {
	if PNMScalePixelDebugPrintTimes < 3 {
		// fmt.Printf("[DEBUG] > > > Scaling pixel [%v][%v] from %+v to %+v\n", pixel.Z, pixel.X, from, to)
		// fmt.Printf("[DEBUG] > > > (((pixel + %v) / %v) * %v) + %v\n", from.Min, from.Max-from.Min, to.Max-to.Min, to.Min)
		PNMScalePixelDebugPrintTimes++
	}
	pnm.Map[pixel.Z][pixel.X] = pnm.ScaleCopyPixel(pixel, from, to)
}

var PNMScalePixelDebugPrintTimes = 0

func (pnm *PerlinNoiseMap) ScaleCopyPixel(pixel coord.Coord, from, to ScaleRange) float64 {
	return (((pnm.Map[pixel.Z][pixel.X] - from.Min) / (from.Max - from.Min)) * (to.Max - to.Min)) + to.Min
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

func (pnm *PerlinNoiseMap) Image(path string) {
	width, height := len(pnm.Map[0]), len(pnm.Map)
	img := image.NewGray(image.Rect(0, 0, width, height))
	// fmt.Println("[DEBUG] > > > > Image base created")

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			//fmt.Printf("[DEBUG] map[%v][%v] = %v", y, x, pnm.Map[y][x])
			//fmt.Printf("[DEBUG] > Image pixel [%v][%v] = ", y, x)

			// Convert to 0-255 scale
			grayColor := uint8(pnm.ScaleCopyPixel(coord.Coord{X: x, Z: y}, pnm.CurrentScale, ScaleRange{0, 255}))

			//fmt.Printf("%+v\n", grayColor)

			img.SetGray(x, y, color.Gray{Y: grayColor})

			//fmt.Printf("[DEBUG] > Image pixel [%v][%v] set\n", y, x)
		}
	}
	// fmt.Println("[DEBUG] > > > > Image pixels setted")

	// Save the image to a file
	file, err := os.Create(fmt.Sprintf(path, PathFilesNumber))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	// fmt.Println("[DEBUG] > > > > Image saved")
}
