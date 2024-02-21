package world

type (
	WorldData struct {
		Width  int
		Height int
		Length int
	}
)

func NewWorldData(width, height, length int) WorldData {
	return WorldData{
		Width:  width,
		Height: height,
		Length: length,
	}
}
