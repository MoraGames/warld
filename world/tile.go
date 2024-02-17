package world

type (
	Tile struct {
		Labels   map[LabelGroup]Label
		Typology Biome
	}
)

func NewTile() Tile {
	return Tile{
		Labels:   make(map[LabelGroup]Label),
		Typology: Biome{},
	}
}
