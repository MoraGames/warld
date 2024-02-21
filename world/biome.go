package world

import (
	"fmt"
	"strings"

	"github.com/MoraGames/warld/seed"
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

	BiomeDesert       = Biome{"Desert", "#FFF08C"}
	BiomeSavanna      = Biome{"Savanna", "#CCD043"}
	BiomePlains       = Biome{"Plains", "#8DD060"}
	BiomeSwamp        = Biome{"Swamp", "#8DB35A"}
	BiomeJungleForest = Biome{"Jungle Forest", "#5DA500"}
	BiomeBambooForest = Biome{"Bamboo Forest", "#78B019"}
	BiomeOakForest    = Biome{"Oak Forest", "#2F912F"}
	BiomeBirchForest  = Biome{"Birch Forest", "#61AE61"}
	BiomeSpruceForest = Biome{"Spruce Forest", "#459045"}
	BiomeSnowyTaiga   = Biome{"Snowy Taiga", "#7BB38C"}
	BiomeSnowyTundra  = Biome{"Snowy Tundra", "#D0F1E4"}
	BiomeIceSpikes    = Biome{"Ice Spikes", "#E3ECE9"}

	BiomeMountains      = Biome{"Mountains", "#BABABA"}
	BiomeSnowyMountains = Biome{"Snowy Mountains", "#F8F8F8"}

	BiomeError   = Biome{"Error", "#000000"}
	BiomeUnknown = Biome{"Unknown", "#000000"}
)

type (
	Biome struct {
		Name  string
		Color string
	}
)

func BiomesList() []Biome {
	return []Biome{
		BiomeWarmOcean,
		BiomeTemperateOcean,
		BiomeColdOcean,
		BiomeFrozenOcean,
		BiomeWarmDeepOcean,
		BiomeTemperateDeepOcean,
		BiomeColdDeepOcean,
		BiomeFrozenDeepOcean,
		BiomeDesert,
		BiomeSavanna,
		BiomePlains,
		BiomeSwamp,
		BiomeJungleForest,
		BiomeBambooForest,
		BiomeOakForest,
		BiomeBirchForest,
		BiomeSpruceForest,
		BiomeSnowyTaiga,
		BiomeSnowyTundra,
		BiomeIceSpikes,
		BiomeMountains,
		BiomeSnowyMountains,

		BiomeError,
		BiomeUnknown,
	}
}

func NewBiome(name, color string) (Biome, error) {
	if strings.TrimSpace(name) == "" {
		return BiomeError, fmt.Errorf("Invalid name: name is empty")
	}
	if strings.TrimSpace(color) == "" || len(color) != 7 || color[0] != '#' {
		return BiomeError, fmt.Errorf("Invalid color: color is empty or not a valid Hex RGB color")
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

// Ocean [Category, Temperature, Altitude]
// Land [Category, Temperature, Humidity, Altitude]
func GenerateBiome(labels map[LabelGroup]Label, seeder *seed.Seeder) (Biome, error) {
	// If the labels are empty, return an empty biome
	if len(labels) == 0 {
		return Biome{}, fmt.Errorf("Invalid labels: labels is empty")
	}

	// Check if are available all the necessary labels to generate a biome
	category, okCategory := labels[GroupCategory]
	if !okCategory {
		return Biome{}, fmt.Errorf("Invalid labels: category is missing")
	}
	macroTemperature /*, okMacroTemperature*/ := labels[GroupMacroTemperature]
	temperature /*, okTemperature*/ := labels[GroupTemperature]
	// if !okTemperature {
	// 	return Biome{}, fmt.Errorf("Invalid labels: temperature is missing")
	// }
	humidity /*, okHumidity*/ := labels[GroupHumidity]
	// if category == CategoryLand && !okHumidity {
	// 	return Biome{}, fmt.Errorf("Invalid labels: humidity is missing")
	// }
	macroAltitude /*, okMacroAltitude*/ := labels[GroupMacroAltitude]
	altitude /*, okAltitude*/ := labels[GroupAltitude]
	//if category == CategoryOcean && !okAltitude {
	//	return Biome{}, fmt.Errorf("Invalid labels: altitude is missing")
	//}
	variant, okVariant := labels[GroupVariant]

	// Generate the biome
	biome := BiomePlains
	switch category {
	case CategoryOcean:
		switch altitude {
		case Depth0:
			switch macroTemperature {
			case MacroTemperatureWarm:
				biome = BiomeWarmOcean
			case MacroTemperatureTemperate:
				biome = BiomeTemperateOcean
			case MacroTemperatureCold:
				biome = BiomeColdOcean
			case MacroTemperatureFreezing:
				biome = BiomeFrozenOcean
			default:
				biome = BiomeTemperateOcean
			}
		case Depth1:
			switch macroTemperature {
			case MacroTemperatureWarm:
				biome = BiomeWarmOcean
			case MacroTemperatureTemperate:
				biome = BiomeTemperateOcean
			case MacroTemperatureCold:
				biome = BiomeColdOcean
			case MacroTemperatureFreezing:
				biome = BiomeFrozenOcean
			default:
				biome = BiomeTemperateOcean
			}
		default:
			switch macroTemperature {
			case MacroTemperatureWarm:
				biome = BiomeWarmOcean
			case MacroTemperatureTemperate:
				biome = BiomeTemperateOcean
			case MacroTemperatureCold:
				biome = BiomeColdOcean
			case MacroTemperatureFreezing:
				biome = BiomeFrozenOcean
			default:
				biome = BiomeTemperateOcean
			}
		}
	case CategoryLand:
		switch macroAltitude {
		case MacroHeightLow, MacroHeightMedium:
			switch temperature {
			case Temperature7:
				switch humidity {
				case Humidity7:
					if okVariant && variant == VariantSpecial {
						biome = BiomeBambooForest
					} else {
						biome = BiomeJungleForest
					}
				case Humidity6:
					biome = BiomeJungleForest
				case Humidity5:
					biome = BiomePlains
				case Humidity4:
					biome = BiomePlains
				case Humidity3:
					biome = BiomeSavanna
				case Humidity2:
					biome = BiomeDesert
				case Humidity1:
					biome = BiomeDesert
				case Humidity0:
					biome = BiomeDesert
				}
			case Temperature6:
				switch humidity {
				case Humidity7:
					biome = BiomeJungleForest
				case Humidity6:
					biome = BiomeJungleForest
				case Humidity5:
					biome = BiomeJungleForest
				case Humidity4:
					biome = BiomePlains
				case Humidity3:
					biome = BiomePlains
				case Humidity2:
					biome = BiomeSavanna
				case Humidity1:
					biome = BiomeDesert
				case Humidity0:
					biome = BiomeDesert
				}
			case Temperature5:
				switch humidity {
				case Humidity7:
					biome = BiomeJungleForest
				case Humidity6:
					biome = BiomeSwamp
				case Humidity5:
					biome = BiomePlains
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomePlains
				case Humidity2:
					biome = BiomeSavanna
				case Humidity1:
					biome = BiomeSavanna
				case Humidity0:
					biome = BiomeSavanna
				}
			case Temperature4:
				switch humidity {
				case Humidity7:
					biome = BiomeBirchForest
				case Humidity6:
					biome = BiomeBirchForest
				case Humidity5:
					biome = BiomeOakForest
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomeSavanna
				case Humidity0:
					biome = BiomeSavanna
				}
			case Temperature3:
				switch humidity {
				case Humidity7:
					biome = BiomeSpruceForest
				case Humidity6:
					biome = BiomeSpruceForest
				case Humidity5:
					biome = BiomeSpruceForest
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomePlains
				case Humidity0:
					biome = BiomePlains
				}
			case Temperature2:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSpruceForest
				case Humidity5:
					biome = BiomeSpruceForest
				case Humidity4:
					biome = BiomeSpruceForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomePlains
				case Humidity0:
					biome = BiomePlains
				}
			case Temperature1:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSnowyTaiga
				case Humidity5:
					biome = BiomeSnowyTaiga
				case Humidity4:
					biome = BiomeSnowyTaiga
				case Humidity3:
					biome = BiomeSpruceForest
				case Humidity2:
					biome = BiomeSnowyTundra
				case Humidity1:
					biome = BiomeSnowyTundra
				case Humidity0:
					biome = BiomeSnowyTundra
				}
			case Temperature0:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSnowyTaiga
				case Humidity5:
					biome = BiomeSnowyTaiga
				case Humidity4:
					biome = BiomeSnowyTaiga
				case Humidity3:
					biome = BiomeSnowyTundra
				case Humidity2:
					biome = BiomeSnowyTundra
				case Humidity1:
					biome = BiomeSnowyTundra
				case Humidity0:
					if okVariant && variant == VariantSpecial {
						biome = BiomeIceSpikes
					} else {
						biome = BiomeSnowyTundra
					}
				}
			default:
				biome = BiomePlains
			}
		case MacroHeightHigh:
			biome = BiomeMountains
		case MacroHeightVeryHigh:
			biome = BiomeSnowyMountains
		default:
			switch temperature {
			case Temperature7:
				switch humidity {
				case Humidity7:
					if okVariant && variant == VariantSpecial {
						biome = BiomeBambooForest
					} else {
						biome = BiomeJungleForest
					}
				case Humidity6:
					biome = BiomeJungleForest
				case Humidity5:
					biome = BiomePlains
				case Humidity4:
					biome = BiomePlains
				case Humidity3:
					biome = BiomeSavanna
				case Humidity2:
					biome = BiomeDesert
				case Humidity1:
					biome = BiomeDesert
				case Humidity0:
					biome = BiomeDesert
				}
			case Temperature6:
				switch humidity {
				case Humidity7:
					biome = BiomeJungleForest
				case Humidity6:
					biome = BiomeJungleForest
				case Humidity5:
					biome = BiomeJungleForest
				case Humidity4:
					biome = BiomePlains
				case Humidity3:
					biome = BiomePlains
				case Humidity2:
					biome = BiomeSavanna
				case Humidity1:
					biome = BiomeDesert
				case Humidity0:
					biome = BiomeDesert
				}
			case Temperature5:
				switch humidity {
				case Humidity7:
					biome = BiomeJungleForest
				case Humidity6:
					biome = BiomeSwamp
				case Humidity5:
					biome = BiomePlains
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomePlains
				case Humidity2:
					biome = BiomeSavanna
				case Humidity1:
					biome = BiomeSavanna
				case Humidity0:
					biome = BiomeSavanna
				}
			case Temperature4:
				switch humidity {
				case Humidity7:
					biome = BiomeBirchForest
				case Humidity6:
					biome = BiomeBirchForest
				case Humidity5:
					biome = BiomeOakForest
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomeSavanna
				case Humidity0:
					biome = BiomeSavanna
				}
			case Temperature3:
				switch humidity {
				case Humidity7:
					biome = BiomeSpruceForest
				case Humidity6:
					biome = BiomeSpruceForest
				case Humidity5:
					biome = BiomeSpruceForest
				case Humidity4:
					biome = BiomeOakForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomePlains
				case Humidity0:
					biome = BiomePlains
				}
			case Temperature2:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSpruceForest
				case Humidity5:
					biome = BiomeSpruceForest
				case Humidity4:
					biome = BiomeSpruceForest
				case Humidity3:
					biome = BiomeOakForest
				case Humidity2:
					biome = BiomePlains
				case Humidity1:
					biome = BiomePlains
				case Humidity0:
					biome = BiomePlains
				}
			case Temperature1:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSnowyTaiga
				case Humidity5:
					biome = BiomeSnowyTaiga
				case Humidity4:
					biome = BiomeSnowyTaiga
				case Humidity3:
					biome = BiomeSpruceForest
				case Humidity2:
					biome = BiomeSnowyTundra
				case Humidity1:
					biome = BiomeSnowyTundra
				case Humidity0:
					biome = BiomeSnowyTundra
				}
			case Temperature0:
				switch humidity {
				case Humidity7:
					biome = BiomeSnowyTaiga
				case Humidity6:
					biome = BiomeSnowyTaiga
				case Humidity5:
					biome = BiomeSnowyTaiga
				case Humidity4:
					biome = BiomeSnowyTaiga
				case Humidity3:
					biome = BiomeSnowyTundra
				case Humidity2:
					biome = BiomeSnowyTundra
				case Humidity1:
					biome = BiomeSnowyTundra
				case Humidity0:
					if okVariant && variant == VariantSpecial {
						biome = BiomeIceSpikes
					} else {
						biome = BiomeSnowyTundra
					}
				}
			default:
				biome = BiomePlains
			}
		}
		/*
			switch temperature {
			case TemperatureWarm:
				switch humidity {
				case HumidityMinimal:
					biome = BiomeDesert
				case HumidityLow:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0:
						biome = BiomeDesert
					case 1, 2:
						biome = BiomeSavanna
					case 3:
						biome = BiomePlains
					}
				case HumidityModerate:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1, 2:
						biome = BiomePlains
					case 3:
						biome = BiomeJungleForest
					}
				case HumidityHigh:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1, 3:
						biome = BiomeJungleForest
					case 2:
						if okVariant && variant == VariantSpecial {
							biome = BiomeBambooForest
						} else {
							biome = BiomeJungleForest
						}
					}
				default:
					biome = BiomePlains
				}
			case TemperatureTemperate:
				switch humidity {
				case HumidityMinimal:
					biome = BiomeSavanna
				case HumidityLow:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0:
						biome = BiomeSavanna
					case 1, 2:
						biome = BiomePlains
					case 3:
						biome = BiomeOakForest
					}
				case HumidityModerate:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1, 3:
						biome = BiomeOakForest
					case 2:
						biome = BiomeSwamp
					}
				case HumidityHigh:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0:
						biome = BiomeSwamp
					case 1, 3:
						biome = BiomeBirchForest
					case 2:
						biome = BiomeJungleForest
					}
				default:
					biome = BiomePlains
				}
			case TemperatureCold:
				switch humidity {
				case HumidityMinimal:
					biome = BiomePlains
				case HumidityLow:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1:
						biome = BiomePlains
					case 2, 3:
						biome = BiomeOakForest
					}
				case HumidityModerate:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0:
						biome = BiomeOakForest
					case 1, 2, 3:
						biome = BiomeSpruceForest
					}
				case HumidityHigh:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1, 2:
						biome = BiomeSpruceForest
					case 3:
						biome = BiomeSnowyTaiga
					}
				default:
					biome = BiomePlains
				}
			case TemperatureFreezing:
				switch humidity {
				case HumidityMinimal:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 2, 3:
						biome = BiomeSnowyTundra
					case 1:
						if okVariant && variant == VariantSpecial {
							biome = BiomeIceSpikes
						} else {
							biome = BiomeSnowyTundra
						}
					}
				case HumidityLow:
					randVal := seeder.Random.IntN(4)
					switch randVal {
					case 0, 1, 3:
						biome = BiomeSnowyTundra
					case 2:
						biome = BiomeSpruceForest
					}
				case HumidityModerate:
					biome = BiomeSnowyTaiga
				case HumidityHigh:
					biome = BiomeSnowyTaiga
				default:
					biome = BiomePlains
				}
			default:
				biome = BiomePlains
			}
		*/
	}

	return biome, nil
}
