package world

import (
	"fmt"
	"strings"
)

var (
	BiomeWarmOcean      = Biome{"Warm Ocean", "#1EE1F0"}
	BiomeTemperateOcean = Biome{"Temperate Ocean", "#28C8F5"}
	BiomeColdOcean      = Biome{"Cold Ocean", "#32AFFA"}
	BiomeFrozenOcean    = Biome{"Frozen Ocean", "#D6FFFF"}

	BiomeWarmDeepOcean      = Biome{"Warm Deep Ocean", "#00CAD8"}
	BiomeTemperateDeepOcean = Biome{"Temperate Deep Ocean", "#00B3E4"}
	BiomeColdDeepOcean      = Biome{"Cold Deep Ocean", "#0094F1"}
	BiomeFrozenDeepOcean    = Biome{"Frozen Deep Ocean", "#BEFAFF"}

	BiomeDesert       = Biome{"Desert", "#F0E68C"}
	BiomeSavanna      = Biome{"Savanna", "#C2B280"}
	BiomePlains       = Biome{"Plains", "#8DB360"}
	BiomeSwamp        = Biome{"Swamp", "#2F7F6F"}
	BiomeJungleForest = Biome{"Jungle Forest", "#2F7F6F"}
	BiomeBambooForest = Biome{"Bamboo Forest", "#2F7F6F"}
	BiomeOakForest    = Biome{"Oak Forest", "#2F7F6F"}
	BiomeBirchForest  = Biome{"Birch Forest", "#2F7F6F"}
	BiomeSpruceForest = Biome{"Spruce Forest", "#2F7F6F"}
	BiomeSnowyTaiga   = Biome{"Snowy Taiga", "#2F7F6F"}
	BiomeSnowyTundra  = Biome{"Snowy Tundra", "#2F7F6F"}
	BiomeIceSpikes    = Biome{"Ice Spikes", "#2F7F6F"}
)

type (
	Biome struct {
		Name  string
		Color string
	}
)

func NewBiome(name, color string) (Biome, error) {
	if strings.TrimSpace(name) == "" {
		return Biome{}, fmt.Errorf("Invalid name: name is empty")
	}
	if strings.TrimSpace(color) == "" || len(color) != 7 || color[0] != '#' {
		return Biome{}, fmt.Errorf("Invalid color: color is empty or not a valid Hex RGB color")
	}
	biome := Biome{
		Name:  name,
		Color: color,
	}
	return biome, nil
}

func (b Biome) String() string {
	return b.Name + " (" + b.Color + ")"
}

/*
// Ocean [Category, Temperature, Altitude]
// Land [Category, Temperature, Precipitation]
func GenerateBiome(labels map[LabelGroup]Label) (Biome, error) {
	// If the labels are empty, return an empty biome
	if len(labels) == 0 {
		return Biome{}, fmt.Errorf("Invalid labels: labels is empty")
	}

	// Check if are available all the necessary labels to generate a biome
	category, okCategory := labels[GroupCategory]
	if !okCategory {
		return Biome{}, fmt.Errorf("Invalid labels: category is missing")
	}
	temperature, okTemperature := labels[GroupTemperature]
	if !okTemperature {
		return Biome{}, fmt.Errorf("Invalid labels: temperature is missing")
	}
	precipitation, okPrecipitation := labels[GroupPrecipitation]
	if category == CategoryLand && !okPrecipitation {
		return Biome{}, fmt.Errorf("Invalid labels: precipitation is missing")
	}
	altitude, okAltitude := labels[GroupAltitude]
	if category == CategoryOcean && !okAltitude {
		return Biome{}, fmt.Errorf("Invalid labels: altitude is missing")
	}
	variant, okVariant := labels[GroupVariant]

	// Generate the biome
	biome := Biome{}
	switch category {
	case CategoryOcean:
		switch altitude {
		case DepthLow:
			switch temperature {
			case TemperatureWarm:
				biome = BiomeWarmOcean
			case TemperatureTemperate:
				biome = BiomeTemperateOcean
			case TemperatureCold:
				biome = BiomeColdOcean
			case TemperatureFreezing:
				biome = BiomeFrozenOcean
			}
		case DepthHigh:
			switch temperature {
			case TemperatureWarm:
				biome = BiomeWarmDeepOcean
			case TemperatureTemperate:
				biome = BiomeTemperateDeepOcean
			case TemperatureCold:
				biome = BiomeColdDeepOcean
			case TemperatureFreezing:
				biome = BiomeFrozenDeepOcean
			}
		}
	case CategoryLand:
	}

	return biome, nil
}
*/
