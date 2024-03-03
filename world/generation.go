package world

import (
	"fmt"
	"math"

	"github.com/MoraGames/warld/coord"
	"github.com/MoraGames/warld/noiseMap"
	"github.com/MoraGames/warld/seed"
)

// X Width
// Y Height (Depth if negative)
// Z Length

func InitializeWorld(worldWidth, worldHeight, worldLength int, seeder *seed.Seeder) *World {
	// Create a new world map
	worldMatrix := make(WorldMatrix, worldLength)
	for z := range worldMatrix {
		worldMatrix[z] = make([]Tile, worldWidth)
		for x := range worldMatrix[z] {
			worldMatrix[z][x] = NewTile()
		}
	}

	// Create a new world
	w := World{
		Seeder:          seeder,
		Data:            NewWorldData(worldWidth, worldHeight, worldLength),
		Map:             worldMatrix,
		Temperature:     noiseMap.NewPerlinNoiseMap(worldWidth, worldLength, 0.02, 1.0, 2.0, 0.5, 8, seeder),
		Humidity:        noiseMap.NewPerlinNoiseMap(worldWidth, worldLength, 0.02, 1.0, 2.0, 0.5, 8, seeder),
		Continentalness: noiseMap.NewPerlinNoiseMap(worldWidth, worldLength, 0.01, 1.0, 2.0, 0.5, 8, seeder),
		Altitude:        noiseMap.NewPerlinNoiseMap(worldWidth, worldLength, 0.015, 1.0, 2.0, 0.5, 8, seeder),
		Rivers:          noiseMap.NewPerlinNoiseMap(worldWidth, worldLength, 0.025, 5.0, 1.0, 0.2, 2, seeder),
		Variants:        noiseMap.NewWhiteNoiseMap(worldWidth, worldLength, seeder),
	}

	// Save the original noisemaps as images
	w.Temperature.Image("./noiseMap/runImages/temperature_original%v.png")
	w.Humidity.Image("./noiseMap/runImages/humidity_original%v.png")
	w.Continentalness.Image("./noiseMap/runImages/continentalness_original%v.png")
	w.Altitude.Image("./noiseMap/runImages/altitude_original%v.png")
	w.Rivers.Image("./noiseMap/runImages/rivers_original%v.png")
	w.Variants.Image("./noiseMap/runImages/variants_original%v.png")
	fmt.Println("[DEBUG] > > Noisemap images saved")

	// Generate the base world map (only oceans and continents)
	worldMatrixFirstIteration(&w)

	// Return the world data structure as a pointer
	return &w
}

func worldMatrixFirstIteration(worldPtr *World) {
	// Retrieve the world
	w := *worldPtr

	// Initialize the world map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 30}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// Convert the range value from [-1, 1] to [0, 30]
			w.Continentalness.ScalePixel(coord.Coord{X: x, Z: z}, w.Continentalness.CurrentScale, newScale)

			// Set the typology of the world map tile based on the continentalness noise map
			// 19/30 of the world map is land
			if w.Continentalness.Map[z][x] > (newScale.Max - 19.0) {
				w.Map[z][x].Labels[GroupCategory] = CategoryLand
			} else {
				w.Map[z][x].Labels[GroupCategory] = CategoryOcean
			}
		}
	}

	// Update the continentalness noisemap scale
	w.Continentalness.CurrentScale = newScale

	// Update the world
	*worldPtr = w
}

func AddTemperatures(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Add temperatures to the world map based on the temperature noise map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 8}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// Convert the range value from [-1, 1] to [0, 8]
			w.Temperature.ScalePixel(coord.Coord{X: x, Z: z}, w.Temperature.CurrentScale, newScale)

			// Set the temperature of the world map tile based on the temperature noise map
			switch {
			case w.Temperature.Map[z][x] >= 7.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureWarm
				w.Map[z][x].Labels[GroupTemperature] = Temperature7
			case w.Temperature.Map[z][x] >= 6.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureWarm
				w.Map[z][x].Labels[GroupTemperature] = Temperature6
			case w.Temperature.Map[z][x] >= 5.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureTemperate
				w.Map[z][x].Labels[GroupTemperature] = Temperature5
			case w.Temperature.Map[z][x] >= 4.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureTemperate
				w.Map[z][x].Labels[GroupTemperature] = Temperature4
			case w.Temperature.Map[z][x] >= 3.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureCold
				w.Map[z][x].Labels[GroupTemperature] = Temperature3
			case w.Temperature.Map[z][x] >= 2.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureCold
				w.Map[z][x].Labels[GroupTemperature] = Temperature2
			case w.Temperature.Map[z][x] >= 1.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureFreezing
				w.Map[z][x].Labels[GroupTemperature] = Temperature1
			case w.Temperature.Map[z][x] >= 0.0:
				w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureFreezing
				w.Map[z][x].Labels[GroupTemperature] = Temperature0
			}
		}
	}

	// Update the temperature noisemap scale
	w.Temperature.CurrentScale = newScale

	// If a warm tile has at least 1 cold/freezing neighbor, set it as temperate (5)
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupMacroTemperature] == MacroTemperatureWarm {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get all 8 neighbors of the current coordinate
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)

				// Check if at least 1 neighbor is cold/freezing, if so, set the current as temperate (5)
				for _, neighbor := range neighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroTemperature] == MacroTemperatureCold || w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroTemperature] == MacroTemperatureFreezing {
						w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureTemperate
						w.Map[z][x].Labels[GroupTemperature] = Temperature5
						break
					}
				}
			}
		}
	}

	// If a freezing tile has at least 1 warm/temperate neighbor, set it as cold (2)
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupMacroTemperature] == MacroTemperatureFreezing {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get all 8 neighbors of the current coordinate
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)

				// Check if at least 1 neighbor is warm/temperate, if so, set the current as cold
				for _, neighbor := range neighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroTemperature] == MacroTemperatureWarm || w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroTemperature] == MacroTemperatureTemperate {
						w.Map[z][x].Labels[GroupMacroTemperature] = MacroTemperatureCold
						w.Map[z][x].Labels[GroupTemperature] = Temperature2
						break
					}
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddHumidity(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Add humidity to the world map based on the humidity noise map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 8}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// Convert the range value from [-1, 1] to [0, 8]
			w.Humidity.ScalePixel(coord.Coord{X: x, Z: z}, w.Humidity.CurrentScale, newScale)

			// Set the precipitation of the world map tile based on the temperature noise map
			switch {
			case w.Humidity.Map[z][x] >= 7.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityHigh
				w.Map[z][x].Labels[GroupHumidity] = Humidity7
			case w.Humidity.Map[z][x] >= 6.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityHigh
				w.Map[z][x].Labels[GroupHumidity] = Humidity6
			case w.Humidity.Map[z][x] >= 5.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityModerate
				w.Map[z][x].Labels[GroupHumidity] = Humidity5
			case w.Humidity.Map[z][x] >= 4.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityModerate
				w.Map[z][x].Labels[GroupHumidity] = Humidity4
			case w.Humidity.Map[z][x] >= 3.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityLow
				w.Map[z][x].Labels[GroupHumidity] = Humidity3
			case w.Humidity.Map[z][x] >= 2.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityLow
				w.Map[z][x].Labels[GroupHumidity] = Humidity2
			case w.Humidity.Map[z][x] >= 1.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityMinimal
				w.Map[z][x].Labels[GroupHumidity] = Humidity1
			case w.Humidity.Map[z][x] >= 0.0:
				w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityMinimal
				w.Map[z][x].Labels[GroupHumidity] = Humidity0
			}
		}
	}

	// Update the humidity noisemap scale
	w.Humidity.CurrentScale = newScale

	// If a high humidity tile has at least 1 low/minimal humidity neighbor, set it as moderate (5)
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupMacroHumidity] == MacroHumidityHigh {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get all 8 neighbors of the current coordinate
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)

				// Check if at least 1 neighbor is low/minimal humidity, if so, set the current as moderate (5)
				for _, neighbor := range neighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroHumidity] == MacroHumidityLow || w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroHumidity] == MacroHumidityMinimal {
						w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityModerate
						w.Map[z][x].Labels[GroupHumidity] = Humidity5
						break
					}
				}
			}
		}
	}

	// If a minimal humidity tile has at least 1 high/moderate humidity neighbor, set it as low (3)
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupMacroHumidity] == MacroHumidityMinimal {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get all 8 neighbors of the current coordinate
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)

				// Check if at least 1 neighbor is high/moderate humidity, if so, set the current as low (3)
				for _, neighbor := range neighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroHumidity] == MacroHumidityHigh || w.Map[neighbor.Z][neighbor.X].Labels[GroupMacroHumidity] == MacroHumidityModerate {
						w.Map[z][x].Labels[GroupMacroHumidity] = MacroHumidityLow
						w.Map[z][x].Labels[GroupHumidity] = Humidity3
						break
					}
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddOceanAltitudes(worldPtr *World) { // Retrieve the world map
	w := *worldPtr

	// Calculate the altitude levels of the world map based on the continentalness noise map
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryOcean {
				if w.Continentalness.Map[z][x] <= (w.Continentalness.CurrentScale.Max-19.0)-2.0 {
					w.Map[z][x].Labels[GroupAltitude] = Depth1
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroDepthHigh
				} else { // <= (w.Continentalness.CurrentScale.Max - 19.0)
					w.Map[z][x].Labels[GroupAltitude] = Depth0
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroDepthLow
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddAltitudes(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Calculate the altitude levels of the world map based on the altitude noise map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 15}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryLand {
				// Convert the range value from [-1, 1] to [0, 15]
				w.Altitude.ScalePixel(coord.Coord{X: x, Z: z}, w.Altitude.CurrentScale, newScale)

				// Set the altitude of the world map tile based on the altitude noise map
				switch {
				case w.Altitude.Map[z][x] == 15.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height15
				case w.Altitude.Map[z][x] >= 14.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height14
				case w.Altitude.Map[z][x] >= 13.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height13
				case w.Altitude.Map[z][x] >= 12.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height12
				case w.Altitude.Map[z][x] >= 11.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					w.Map[z][x].Labels[GroupAltitude] = Height11
				case w.Altitude.Map[z][x] >= 10.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					w.Map[z][x].Labels[GroupAltitude] = Height10
				case w.Altitude.Map[z][x] >= 9.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height9
				case w.Altitude.Map[z][x] >= 8.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height8
				case w.Altitude.Map[z][x] >= 7.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height7
				case w.Altitude.Map[z][x] >= 6.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height6
				case w.Altitude.Map[z][x] >= 5.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height5
				case w.Altitude.Map[z][x] >= 4.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height4
				case w.Altitude.Map[z][x] >= 3.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height3
				case w.Altitude.Map[z][x] >= 2.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height2
				case w.Altitude.Map[z][x] >= 1.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height1
				case w.Altitude.Map[z][x] >= 0.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height0
				}
			}
		}
	}

	// Update the altitude noisemap scale
	w.Altitude.CurrentScale = newScale

	// Update the world map
	*worldPtr = w
}

func AdjustVariantsAltitudes(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Transform special plains into collinar plains
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			//If a tile is a special plains, set it as collinar plains
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
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height7
					case Height7:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height8
					case Height8:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height9
					case Height9:
						// In this case becomes a Mountains biome
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height10
						worldPtr.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					}
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddVariants(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Add variants to the world map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 13}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// Convert the range value from [0, 1] to [0, 13]
			w.Variants.ScalePixel(coord.Coord{X: x, Z: z}, w.Variants.CurrentScale, newScale)

			// Set the variant of the world map tile based on the variant noise map
			if w.Variants.Map[z][x] >= 12.0 {
				w.Map[z][x].Labels[GroupVariant] = VariantSpecial
			}
		}
	}

	// Update the variant noisemap scale
	w.Variants.CurrentScale = newScale

	// Update the world map
	*worldPtr = w
}

func AddRivers(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Add rivers to the world map based on the river noise map
	newScale := noiseMap.ScaleRange{Min: -15, Max: 15}
	// fmt.Printf("[DEBUG] > > > Current range and new range: %+v => %+v\n", w.Rivers.CurrentScale, newScale)
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// Convert the range value from [-1, 1] to [-15, 15]
			w.Rivers.ScalePixel(coord.Coord{X: x, Z: z}, w.Rivers.CurrentScale, newScale)

			// Make the river noise map equal to the absolute value of it's self
			w.Rivers.Map[z][x] = math.Abs(w.Rivers.Map[z][x])

			// Filter the river noise map
			if w.Rivers.Map[z][x] <= 0.5 {
				if w.Map[z][x].Labels[GroupCategory] == CategoryLand {
					if w.Map[z][x].Labels[GroupMacroAltitude] == MacroHeightLow || w.Map[z][x].Labels[GroupMacroAltitude] == MacroHeightMedium || w.Map[z][x].Labels[GroupAltitude] == Height10 {
						w.Map[z][x].Labels[GroupCategory] = CategoryRiver
					}
				}
			}
		}
	}

	// Update the river noisemap scale
	w.Rivers.CurrentScale = newScale

	// Save the edited river noise map as an image
	w.Rivers.Image("./noiseMap/runImages/rivers_edited%v.png")

	// Update the world map
	*worldPtr = w
}
