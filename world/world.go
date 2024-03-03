package world

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/MoraGames/warld/noiseMap"
	"github.com/MoraGames/warld/seed"
)

const (
	GlobalDebugMode = true
)

var (
	PathFilesNumber = "00"
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
		Rivers          noiseMap.PerlinNoiseMap
		Variants        noiseMap.WhiteNoiseMap
	}
)

func Create(worldWidth, worldHeight, worldLength int, seeder *seed.Seeder) *World {
	fmt.Println("[DEBUG] > Creating world...")

	// Create a new world
	worldPtr := InitializeWorld(worldWidth, worldHeight, worldLength, seeder)
	fmt.Println("[DEBUG] > World initialized")

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("Initialization\n", worldPtr.String(false))

	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println(worldPtr.String(false))

	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)
	//AddIslands(worldPtr)
	//AddIslands(worldPtr)
	//AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("x2 Zoom + x5 AddIslands\n", worldPtr.String(false))

	//RemoveTooMuchOceans(worldPtr)
	//AddIslands(worldPtr)

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("RemoveOceans + AddIslands\n", worldPtr.String(false))

	AddTemperatures(worldPtr)
	fmt.Println("[DEBUG] > AddTemperatures() done")
	AddVariants(worldPtr)
	fmt.Println("[DEBUG] > AddVariants() done")

	// // Update the biomes
	// for z := 0; z < worldPtr.Data.ActualWidth; z++ {
	// 	for x := 0; x < worldPtr.Data.ActualLength; x++ {
	// 		worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
	// 	}
	// }
	// fmt.Println("AddTemperatures + Add Variants\n", worldPtr.String(false))

	//ZoomWorld(worldPtr)
	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)
	AddHumidity(worldPtr)
	fmt.Println("[DEBUG] > AddHumidity() done")
	AddOceanAltitudes(worldPtr)
	fmt.Println("[DEBUG] > AddOceanAltitudes() done")
	AddAltitudes(worldPtr)
	fmt.Println("[DEBUG] > AddAltitudes() done")
	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)
	//ZoomWorld(worldPtr)
	//AddIslands(worldPtr)

	// Update the biomes
	for z := 0; z < worldPtr.Data.Length; z++ {
		for x := 0; x < worldPtr.Data.Width; x++ {
			worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
		}
	}
	fmt.Println("[DEBUG] > Biomes updated")

	AdjustVariantsAltitudes(worldPtr)
	fmt.Println("[DEBUG] > AdjustVariantsAltitudes() done")
	AddRivers(worldPtr)
	fmt.Println("[DEBUG] > AddRivers() done")

	// Update the biomes
	for z := 0; z < worldPtr.Data.Length; z++ {
		for x := 0; x < worldPtr.Data.Width; x++ {
			//fmt.Printf("[DEBUG] > Updating biome [%v %v]: ", z, x)
			worldPtr.Map[z][x].UpdateBiome(worldPtr.Seeder)
			//fmt.Println()
		}
	}
	fmt.Println("[DEBUG] > Biomes updated")

	return worldPtr
}

func (w *World) String(withLegend bool) string {
	str := ""
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if GlobalDebugMode {
				str += w.Map[z][x].StringDebug()
			} else {
				str += w.Map[z][x].String()
			}
		}
		str += "\n"
	}

	if withLegend {
		str += BiomesLegend()
	}

	return str
}

func (w *World) Print(withLegend bool) {
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if GlobalDebugMode {
				fmt.Print(w.Map[z][x].StringDebug())
			} else {
				fmt.Print(w.Map[z][x].String())
			}
		}
		fmt.Println()
	}
	fmt.Println()

	if withLegend {
		fmt.Println(BiomesLegend())
	}
}

func (w *World) Image(path string) {
	// Create the base image
	width, height := len(w.Map[0]), len(w.Map)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Convert each pixel to its color
			c, err := w.Map[y][x].Typology.RGBAColor()
			if err != nil {
				panic(err)
			}
			img.Set(x, y, c)
		}
	}

	// Save the image to a file
	file, err := os.Create(fmt.Sprintf(path, PathFilesNumber))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func (w *World) PrintDebugInfos() {

	fmt.Printf("World Seed: %v | World Size: %vl x %vh x %vw | ", w.Seeder.GetSeed().Token, w.Data.Width, w.Data.Height, w.Data.Length)

	totTiles := w.Data.Length * w.Data.Width

	totOceansNum := 0
	notDeepOceansNum, deepOceansNum := 0, 0
	warmOceansNum, temperateOceansNum, coldOceansNum, freezingOceansNum := 0, 0, 0, 0
	warmNotDeepOceansNum, temperateNotDeepOceansNum, coldNotDeepOceansNum, frozenNotDeepOceansNum := 0, 0, 0, 0
	warmDeepOceansNum, temperateDeepOceansNum, coldDeepOceansNum, frozenDeepOceansNum := 0, 0, 0, 0

	totLandsNum := 0
	warmLandsNum, temperateLandsNum, coldLandsNum, frozenLandsNum := 0, 0, 0, 0
	lowHeightLandsNum, mediumHeightLandsNum, highHeightLandsNum, veryHighHeightLandsNum := 0, 0, 0, 0
	minimalHumidityLandsNum, lowHumidityLandsNum, moderateHumidityLandsNum, highHumidityLandsNum := 0, 0, 0, 0
	normalLandsNum, collinarLandsNum, specialLandsNum := 0, 0, 0
	desertBiomesNum, savannaBiomesNum, plainsBiomesNum, swampBiomesNum, jungleForestBiomesNum, bambooForestBiomesNum, oakForestBiomesNum, birchForestBiomesNum, spruceForestBiomesNum, snowyTaigaBiomesNum, snowyTundraBiomesNum, iceSpikesBiomesNum, mountainsBiomesNum, snowyMountainsBiomesNum := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0

	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryOcean {
				totOceansNum++
				if w.Map[z][x].Labels[GroupAltitude] == Depth0 {
					notDeepOceansNum++
					switch w.Map[z][x].Labels[GroupMacroTemperature] {
					case MacroTemperatureWarm:
						warmNotDeepOceansNum++
						warmOceansNum++
					case MacroTemperatureTemperate:
						temperateNotDeepOceansNum++
						temperateOceansNum++
					case MacroTemperatureCold:
						coldNotDeepOceansNum++
						coldOceansNum++
					case MacroTemperatureFreezing:
						frozenNotDeepOceansNum++
						freezingOceansNum++
					}
				} else if w.Map[z][x].Labels[GroupAltitude] == Depth1 {
					deepOceansNum++
					switch w.Map[z][x].Labels[GroupMacroTemperature] {
					case MacroTemperatureWarm:
						warmDeepOceansNum++
						warmOceansNum++
					case MacroTemperatureTemperate:
						temperateDeepOceansNum++
						temperateOceansNum++
					case MacroTemperatureCold:
						coldDeepOceansNum++
						coldOceansNum++
					case MacroTemperatureFreezing:
						frozenDeepOceansNum++
						freezingOceansNum++
					}
				}
			} else if w.Map[z][x].Labels[GroupCategory] == CategoryLand {
				totLandsNum++
				switch w.Map[z][x].Labels[GroupMacroTemperature] {
				case MacroTemperatureWarm:
					warmLandsNum++
				case MacroTemperatureTemperate:
					temperateLandsNum++
				case MacroTemperatureCold:
					coldLandsNum++
				case MacroTemperatureFreezing:
					frozenLandsNum++
				}
				switch w.Map[z][x].Labels[GroupMacroAltitude] {
				case MacroHeightLow:
					lowHeightLandsNum++
				case MacroHeightMedium:
					mediumHeightLandsNum++
				case MacroHeightHigh:
					highHeightLandsNum++
				case MacroHeightVeryHigh:
					veryHighHeightLandsNum++
				}
				switch w.Map[z][x].Labels[GroupMacroHumidity] {
				case MacroHumidityMinimal:
					minimalHumidityLandsNum++
				case MacroHumidityLow:
					lowHumidityLandsNum++
				case MacroHumidityModerate:
					moderateHumidityLandsNum++
				case MacroHumidityHigh:
					highHumidityLandsNum++
				}
				switch w.Map[z][x].Labels[GroupVariant] {
				case VariantCollinar:
					collinarLandsNum++
				case VariantSpecial:
					specialLandsNum++
				default:
					normalLandsNum++
				}
				switch w.Map[z][x].Typology {
				case BiomeDesert:
					desertBiomesNum++
				case BiomeSavanna:
					savannaBiomesNum++
				case BiomePlains:
					plainsBiomesNum++
				case BiomeSwamp:
					swampBiomesNum++
				case BiomeJungleForest:
					jungleForestBiomesNum++
				case BiomeBambooForest:
					bambooForestBiomesNum++
				case BiomeOakForest:
					oakForestBiomesNum++
				case BiomeBirchForest:
					birchForestBiomesNum++
				case BiomeSpruceForest:
					spruceForestBiomesNum++
				case BiomeSnowyTaiga:
					snowyTaigaBiomesNum++
				case BiomeSnowyTundra:
					snowyTundraBiomesNum++
				case BiomeIceSpikes:
					iceSpikesBiomesNum++
				case BiomeMountains:
					mountainsBiomesNum++
				case BiomeSnowyMountains:
					snowyMountainsBiomesNum++
				}
			}
		}
	}

	fmt.Printf("Oceans: %v | Lands: %v | Total: %v\n", totOceansNum, totLandsNum, totTiles)
	fmt.Printf("Oceans:\n")
	fmt.Printf("  - Not Deep: %v | Deep: %v\n  - Warm: %v | Temperate: %v | Cold: %v | Freezing: %v\n", notDeepOceansNum, deepOceansNum, warmOceansNum, temperateOceansNum, coldOceansNum, freezingOceansNum)
	fmt.Printf("  - Warm Not Deep: %v | Temperate Not Deep: %v | Cold Not Deep: %v | Frozen Not Deep: %v\n", warmNotDeepOceansNum, temperateNotDeepOceansNum, coldNotDeepOceansNum, frozenNotDeepOceansNum)
	fmt.Printf("  - Warm Deep: %v | Temperate Deep: %v | Cold Deep: %v | Frozen Deep: %v\n", warmDeepOceansNum, temperateDeepOceansNum, coldDeepOceansNum, frozenDeepOceansNum)
	fmt.Printf("Lands:\n")
	fmt.Printf("  - Warm: %v | Temperate: %v | Cold: %v | Freezing: %v\n", warmLandsNum, temperateLandsNum, coldLandsNum, frozenLandsNum)
	fmt.Printf("  - Low Height: %v | Medium Height: %v | High Height: %v | Very High Height: %v\n", lowHeightLandsNum, mediumHeightLandsNum, highHeightLandsNum, veryHighHeightLandsNum)
	fmt.Printf("  - Minimal Humidity: %v | Low Humidity: %v | Moderate Humidity: %v | High Humidity: %v\n", minimalHumidityLandsNum, lowHumidityLandsNum, moderateHumidityLandsNum, highHumidityLandsNum)
	fmt.Printf("  - Normal: %v | Collinar: %v | Special: %v\n", normalLandsNum, collinarLandsNum, specialLandsNum)
	fmt.Printf("  - Desert: %v | Savanna: %v | Plains: %v | Swamp: %v | Jungle Forest: %v | Bamboo Forest: %v | Oak Forest: %v | Birch Forest: %v | Spruce Forest: %v | Snowy Taiga: %v | Snowy Tundra: %v | Ice Spikes: %v | Mountains: %v | Snowy Mountains: %v\n", desertBiomesNum, savannaBiomesNum, plainsBiomesNum, swampBiomesNum, jungleForestBiomesNum, bambooForestBiomesNum, oakForestBiomesNum, birchForestBiomesNum, spruceForestBiomesNum, snowyTaigaBiomesNum, snowyTundraBiomesNum, iceSpikesBiomesNum, mountainsBiomesNum, snowyMountainsBiomesNum)
}

func StringToFile(text, path string) {
	file, err := os.Create(fmt.Sprintf(path, PathFilesNumber))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
}
