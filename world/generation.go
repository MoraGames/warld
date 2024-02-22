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
		Variants:        noiseMap.NewWhiteNoiseMap(worldWidth, worldLength, seeder),
	}

	// Save the original noisemaps as images
	w.Temperature.Image("./noiseMap/runImages/temperature_original%v.png")
	w.Humidity.Image("./noiseMap/runImages/humidity_original%v.png")
	w.Continentalness.Image("./noiseMap/runImages/continentalness_original%v.png")
	w.Altitude.Image("./noiseMap/runImages/altitude_original%v.png")
	w.Variants.Image("./noiseMap/runImages/variants_original%v.png")
	fmt.Println("[DEBUG] > > Noisemap images saved")

	// Generate the base world map (only oceans and continents)
	worldMatrixFirstIteration(&w)
	fmt.Println("[DEBUG] > > First iteration done")

	// Return the world data structure as a pointer
	return &w
}

func worldMatrixFirstIteration(worldPtr *World) {
	// Retrieve the world
	w := *worldPtr

	// Initialize the world map
	oceansNum, landsNum := 0, 0
	minValBS, maxValBS, minValAS, maxValAS := math.MaxFloat64, math.SmallestNonzeroFloat64, math.MaxFloat64, math.SmallestNonzeroFloat64
	minValIndex, maxValIndex := coord.NewCoord(0, 0), coord.NewCoord(0, 0)
	newScale := noiseMap.ScaleRange{Min: 0, Max: 15}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Continentalness.Map[z][x] < minValBS {
				minValBS = w.Continentalness.Map[z][x]
				minValIndex = coord.NewCoord(z, x)
			}
			if w.Continentalness.Map[z][x] > maxValBS {
				maxValBS = w.Continentalness.Map[z][x]
				maxValIndex = coord.NewCoord(z, x)
			}

			// Convert the range value from [-1, 1] to [0, 15]
			w.Continentalness.ScalePixel(coord.Coord{X: x, Z: z}, w.Continentalness.CurrentScale, newScale)

			if w.Continentalness.Map[z][x] < minValAS {
				minValAS = w.Continentalness.Map[z][x]
			}
			if w.Continentalness.Map[z][x] > maxValAS {
				maxValAS = w.Continentalness.Map[z][x]
			}

			// Set the typology of the world map tile based on the continentalness noise map
			// 8/15 of the world map is land
			if w.Continentalness.Map[z][x] > (newScale.Max - 8.0) {
				w.Map[z][x].Labels[GroupCategory] = CategoryLand
				landsNum++
			} else {
				w.Map[z][x].Labels[GroupCategory] = CategoryOcean
				oceansNum++
			}
		}
	}
	w.Continentalness.CurrentScale = newScale
	noiseMap.PNMScalePixelDebugPrintTimes = 0

	fmt.Printf("[DEBUG] > > > Oceans: %v/%v | Lands: %v/%v\n", oceansNum, w.Data.Length*w.Data.Width, landsNum, w.Data.Length*w.Data.Width)
	fmt.Printf("[DEBUG] > > > Continentalness registered range: [%v, %v] | Scaled registered range: [%v, %v]\n", minValBS, maxValBS, minValAS, maxValAS)
	fmt.Printf("[DEBUG] > > > Continentalness scale operation on min/max values:\n   From %v (index %v) ==> To %v\n   From %v (index %v) ==> To %v\n", minValBS, minValIndex, w.Continentalness.Map[minValIndex.Z][minValIndex.X], maxValBS, maxValIndex, w.Continentalness.Map[maxValIndex.Z][maxValIndex.X])

	// Update the world
	*worldPtr = w
}

func ZoomWorld(worldPtr *World) {
	// Retrieve the world
	w := *worldPtr

	// Zoom the world map (each tile becomes 2x2 tiles)
	worldMatrix := make(WorldMatrix, worldPtr.Data.Length*2)
	for z := range worldMatrix {
		worldMatrix[z] = make([]Tile, worldPtr.Data.Width*2)
		for x := range worldMatrix[z] {
			//worldMatrix[z][x] = NewTile()
			// Copy the tile's data from the old world map to the new world map
			worldMatrix[z][x] = w.Map[z/2][x/2]
		}
	}

	w.Data.Width *= 2
	w.Data.Length *= 2
	w.Map = worldMatrix

	// Introduce some randomness variation to the new world map
	Smudge(&w)

	// Update the world
	*worldPtr = w
}

func Smudge(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Smudge the world map
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryLand {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				//fmt.Printf("[DEBUG] > z = %v | x = %v | lenght = %v | width = %v | map lenght = %v | map width = %v\n", z, x, w.Data.Length, w.Data.Width, len(w.Map), len(w.Map[0]))

				// Get all 8 neighbors of the current coordinate and choose randomly 3 of them
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)
				choosedNeighbor := neighborsCoords[w.Seeder.Random.IntN(len(neighborsCoords))]

				//fmt.Printf("[DEBUG] -> current = %+v | choosed neighbor = %+v\n", current, neighbor)

				// Set it equal to the current with a 75% chance, if fails, set the current equal to the neighbor with a 75% chance
				if w.Seeder.Random.IntN(4) == 0 {
					w.Map[choosedNeighbor.Z][choosedNeighbor.X] = w.Map[z][x]
				} else {
					if w.Seeder.Random.IntN(4) == 0 {
						w.Map[z][x] = w.Map[choosedNeighbor.Z][choosedNeighbor.X]
					}
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddIslands(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Generate the coordinates of a random amount (based on the map size) of islands
	//proportionalAmount := int(math.Round(float64(w.Data.Length*w.Data.Width) / 8))
	//amountVariation := w.Seeder.Random.IntN(5) - 2 // [-2,+2]
	// amount := proportionalAmount + amountVariation
	amount := int(math.Round(float64(w.Data.Length*w.Data.Width) / 10))
	randomCoords := make(coord.CoordSlice, 0)
	for i := 0; i < amount; i++ {
		randomCoords = append(randomCoords, coord.NewCoord(w.Seeder.Random.IntN(w.Data.Length), w.Seeder.Random.IntN(w.Data.Width)))
	}

	// For each coordinate, set it as land (create an island)
	for _, coord := range randomCoords {
		w.Map[coord.Z][coord.X].Labels[GroupCategory] = CategoryLand
	}

	// Update the world map
	*worldPtr = w
}

func RemoveTooMuchOceans(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Remove too much oceans from the world map, replacing them with land
	// If a tile is ocean and as aligned neighbors all oceans, set it as land with a 50% chance
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryOcean {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get the 4 aligned neighbors of the current coordinate
				alignedNeighborsCoords := current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1)

				// Check if all the aligned neighbors are oceans
				allOceans := true
				for _, neighbor := range alignedNeighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupCategory] != CategoryOcean {
						allOceans = false
						break
					}
				}

				// If all the aligned neighbors are oceans, set the current as land with a 50% chance
				if allOceans {
					if w.Seeder.Random.IntN(2) == 0 {
						w.Map[z][x].Labels[GroupCategory] = CategoryLand
					}
				}
			}
		}
	}

	// Update the world map
	*worldPtr = w
}

func AddTemperatures(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Add temperatures to the world map based on the temperature noise map
	newScale := noiseMap.ScaleRange{Min: 0, Max: 8}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			// if z > w.Data.Length-3 || x > w.Data.Width-3 {
			// 	fmt.Printf("[DEBUG] -> z = %v | x = %v | lenght = %v | width = %v | map lenght = %v | map width = %v\n", z, x, w.Data.Length, w.Data.Width, len(w.Map), len(w.Map[0]))
			// }
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
	w.Temperature.CurrentScale = newScale
	noiseMap.PNMScalePixelDebugPrintTimes = 0

	fmt.Println("[DEBUG] > > Temperatures added")

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
	fmt.Println("[DEBUG] > > Temperatures adjusted")

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
	w.Humidity.CurrentScale = newScale
	noiseMap.PNMScalePixelDebugPrintTimes = 0
	fmt.Println("[DEBUG] > > Humidity added")

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
	fmt.Println("[DEBUG] > > Humidity adjusted")

	// Update the world map
	*worldPtr = w
}

func AddOceanAltitudes(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Each tile of water is set with label DepthHigh if all 8 neighbors are water, otherwise DepthLow
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryOcean {
				// Get the current coordinate
				current := coord.NewCoord(z, x)

				// Get all 8 neighbors of the current coordinate
				neighborsCoords := coord.ConcatenateCoordSlices(
					current.GetAlignedNeighbors(w.Data.Length-1, w.Data.Width-1),
					current.GetDiagonalNeighbors(w.Data.Length-1, w.Data.Width-1),
				)

				// Check if all the neighbors are water, if so, set the current as DepthHigh
				allWater := true
				for _, neighbor := range neighborsCoords {
					if w.Map[neighbor.Z][neighbor.X].Labels[GroupCategory] != CategoryOcean {
						allWater = false
						break
					}
				}
				if allWater {
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroDepthHigh
					w.Map[z][x].Labels[GroupAltitude] = Depth1
				} else {
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroDepthLow
					w.Map[z][x].Labels[GroupAltitude] = Depth0
				}
			}
		}
	}
}

func AddAltitudes(worldPtr *World) {
	// Retrieve the world map
	w := *worldPtr

	// Calculate the altitude levels of the world map based on the altitude noise map
	oceansNum, landsNum := 0, 0
	minValBS, maxValBS, minValAS, maxValAS := math.MaxFloat64, math.SmallestNonzeroFloat64, math.MaxFloat64, math.SmallestNonzeroFloat64
	newScale := noiseMap.ScaleRange{Min: 0, Max: 10}
	for z := 0; z < w.Data.Length; z++ {
		for x := 0; x < w.Data.Width; x++ {
			if w.Map[z][x].Labels[GroupCategory] == CategoryLand {
				landsNum++
				if w.Altitude.Map[z][x] < minValBS {
					minValBS = w.Altitude.Map[z][x]
				}
				if w.Altitude.Map[z][x] > maxValBS {
					maxValBS = w.Altitude.Map[z][x]
				}

				// Convert the range value from [-1, 1] to [0, 10]
				w.Altitude.ScalePixel(coord.Coord{X: x, Z: z}, w.Altitude.CurrentScale, newScale)

				if w.Altitude.Map[z][x] < minValAS {
					minValAS = w.Altitude.Map[z][x]
				}
				if w.Altitude.Map[z][x] > maxValAS {
					maxValAS = w.Altitude.Map[z][x]
				}

				// Set the altitude of the world map tile based on the altitude noise map
				switch {
				case w.Altitude.Map[z][x] >= 9.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height9
				case w.Altitude.Map[z][x] >= 8.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightVeryHigh
					w.Map[z][x].Labels[GroupAltitude] = Height8
				case w.Altitude.Map[z][x] >= 7.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					w.Map[z][x].Labels[GroupAltitude] = Height7
				case w.Altitude.Map[z][x] >= 6.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					w.Map[z][x].Labels[GroupAltitude] = Height6
				case w.Altitude.Map[z][x] >= 5.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightHigh
					w.Map[z][x].Labels[GroupAltitude] = Height5
				case w.Altitude.Map[z][x] >= 4.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height4
				case w.Altitude.Map[z][x] >= 3.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height3
				case w.Altitude.Map[z][x] >= 2.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					w.Map[z][x].Labels[GroupAltitude] = Height2
				case w.Altitude.Map[z][x] >= 1.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height1
				case w.Altitude.Map[z][x] >= 0.0:
					w.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightLow
					w.Map[z][x].Labels[GroupAltitude] = Height0
				}
			} else {
				oceansNum++

				// Burn the value (convert the range value from [-1, 1] to [0, 0])
				w.Altitude.ScalePixel(coord.Coord{X: x, Z: z}, w.Altitude.CurrentScale, noiseMap.ScaleRange{Min: 0, Max: 0})
			}
		}
	}
	w.Altitude.CurrentScale = newScale
	noiseMap.PNMScalePixelDebugPrintTimes = 0

	fmt.Printf("[DEBUG] > > > Oceans: %v/%v | Lands: %v/%v\n", oceansNum, w.Data.Length*w.Data.Width, landsNum, w.Data.Length*w.Data.Width)
	fmt.Printf("[DEBUG] > > > Altitude registered range: [%v, %v] | Scaled registered range: [%v, %v]\n", minValBS, maxValBS, minValAS, maxValAS)

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
						worldPtr.Map[z][x].Labels[GroupMacroAltitude] = MacroHeightMedium
					case Height2:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height3
					case Height3:
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height4
					case Height4:
						// In this case becomes a Mountains biome
						worldPtr.Map[z][x].Labels[GroupAltitude] = Height5
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
	w.Variants.CurrentScale = newScale
	noiseMap.WNMScalePixelDebugPrintTimes = 0

	// Update the world map
	*worldPtr = w
}
