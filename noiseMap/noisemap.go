package noiseMap

import "github.com/MoraGames/warld/coord"

var (
	PathFilesNumber string
)

type (
	ScaleRange struct {
		Min float64
		Max float64
	}

	NoiseMap interface {
		ScalePixel(pixel coord.Coord, from, to ScaleRange)
		ScaleCopyPixel(pixel coord.Coord, from, to ScaleRange) float64
		String() string
		Image(path string)
	}
)
