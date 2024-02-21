package world

import (
	"github.com/MoraGames/warld/noiseMap"
	"github.com/MoraGames/warld/seed"
	"github.com/gookit/color"
)

type (
	World struct {
		Seeder *seed.Seeder
		Data   WorldData

		Map WorldMatrix

		Temperature     noiseMap.PerlinNoiseMap
		Humidity        noiseMap.PerlinNoiseMap
		Continentalness noiseMap.PerlinNoiseMap
		Altitude        noiseMap.PerlinNoiseMap
		Variants        noiseMap.WhiteNoiseMap
	}
)

func Create(worldWidth, worldHeight, worldLength int, seeder *seed.Seeder) *World {
	// Create a new world
	worldPtr := InitializeWorld(worldWidth, worldHeight, worldLength, seeder)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("Initialization\n", worldPtr.String(false))

	ZoomWorld(worldPtr)
	AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println(worldPtr.String(false))

	ZoomWorld(worldPtr)
	AddIslands(worldPtr)
	AddIslands(worldPtr)
	AddIslands(worldPtr)
	AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("x2 Zoom + x5 AddIslands\n", worldPtr.String(false))

	RemoveTooMuchOceans(worldPtr)
	AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("RemoveOceans + AddIslands\n", worldPtr.String(false))

	AddTemperatures(worldPtr)
	AddVariants(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("AddTemperatures + Add Variants\n", worldPtr.String(false))

	ZoomWorld(worldPtr)
	ZoomWorld(worldPtr)
	AddIslands(worldPtr)
	AddHumidity(worldPtr)
	AddOceanAltitudes(worldPtr)
	AddAltitudes(worldPtr)
	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)
	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)

	// Update the biomes
	for z := 0; z < worldPtr.Data.Length; z++ {
		for x := 0; x < worldPtr.Data.Width; x++ {
			worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)

			// Transform in "Collinar" the "Special" plains
			if worldPtr.Map[z][x].Typology == BiomePlains {
				if variant, okVariant := worldPtr.Map[z][x].Labels[GroupVariant]; okVariant && variant == VariantSpecial {
					// Update the labels
					worldPtr.Map[z][x].Labels[GroupVariant] = VariantCollinar

					//Adjust the altitude
					switch worldPtr.Map[z][x].Labels[GroupAltitude] {
					case Height0:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height1
					case Height1:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height2
					case Height2:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height3
					case Height3:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height4
						worldPtr.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					case Height4:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height5
					case Height5:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height6
					case Height6:
						// In this case becomes a Mountains biome
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height7
						worldPtr.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
						worldPtr.Map[z][x].Typology = BiomeMountains

						// Re-update the biome
						worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
					}
				}
			}
		}
	}

	return worldPtr
}

func (w *World) String(withLegend bool) string {
	str := ""
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			str += w.Map[z][x].String()
		}
		str += "\n"
	}

	if withLegend {
		str += "\n Biomes Legend:\n"
		for _, biome := range BiomesList() {
			str += "   " + color.NewRGBStyle(color.HEX("222222"), color.HEX(biome.Color)).Sprint(" ") + " " + biome.Name + "\n"
		}
	}

	return str
}
