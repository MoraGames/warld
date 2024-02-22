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

func (t *Tile) StringDebug() string {
	return color.NewRGBStyle(color.HEX("222222"), color.HEX(t.Typology.Color)).Sprint(t.PlainStringDebug())
}

func (t *Tile) PlainStringDebug() string {
	str := " "
	switch t.Labels[GroupCategory] {
	case CategoryLand:
		switch t.Labels[GroupVariant] {
		case VariantCollinar:
			str = "c"
		case VariantSpecial:
			str = "s"
		}
		switch t.Typology {
		case BiomeOakForest:
			str = "o"
		case BiomeBirchForest:
			str = "b"
		case BiomeSpruceForest:
			str = "s"
		case BiomeJungleForest:
			str = "j"
		case BiomeSwamp:
			str = "@"
		}
	case CategoryOcean:
		switch t.Labels[GroupAltitude] {
		case Depth1:
			str = "-"
		case Depth0:
			str = "+"
		}
	}
	return str
}
