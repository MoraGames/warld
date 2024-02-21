package world

import (
	"github.com/MoraGames/warld/seed"
	"github.com/gookit/color"
)

type (
	WorldMatrix [][]Tile
	Tile        struct {
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

func CopyTile(source, destination *Tile) {
	destination.Labels = make(map[LabelGroup]Label)
	for k, v := range source.Labels {
		destination.Labels[k] = v
	}
	destination.Typology = source.Typology
}

func (t *Tile) UpdateBiome(seeder *seed.Seeder) {
	biome, err := GenerateBiome(t.Labels, seeder)
	if err != nil {
		panic(err)
	}
	t.Typology = biome
}

func (t *Tile) String() string {
	return color.NewRGBStyle(color.HEX("222222"), color.HEX(t.Typology.Color)).Sprint(" ")
}
