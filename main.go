package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MoraGames/terminal"
	"github.com/MoraGames/warld/noiseMap"
	"github.com/MoraGames/warld/seed"
	"github.com/MoraGames/warld/world"
)

const (
	EnableTest = true
)

var (
	WorldSeeder     *seed.Seeder
	WorldFilePath   = "./saves/world%v.png"
	WorldDimensions = NewDimensions("(2048,2048)")
)

func init() {
	// Obtain the filesNumber
	var pathFilesNumberStr string
	fmt.Print("Enter the files number: ")
	Scan(&pathFilesNumberStr)
	noiseMap.PathFilesNumber, world.PathFilesNumber = pathFilesNumberStr, pathFilesNumberStr

	// Obtain the world seed
	var seedStr string
	fmt.Print("Enter the world seed: ")
	Scan(&seedStr)

	// Obtain the world dimensions
	var dimensionsStr string
	fmt.Print("Enter the world (length,width): ")
	Scan(&dimensionsStr)
	WorldDimensions = NewDimensions(dimensionsStr)

	// Create the random number generator based on the seed
	seedPtr, err := seed.GetSeed(seedStr)
	if err != nil {
		log.Panicln(err)
	}
	WorldSeeder = seed.NewSeeder(seedPtr)

	// Clear the console
	time.Sleep(2 * time.Second)
	terminal.Clear()
}

func main() {
	var t0_wc time.Time
	if EnableTest {
		t0_wc = time.Now()
		testingNoiseMaps(WorldDimensions.Length, WorldDimensions.Width)
	}

	t1_wc := time.Now()
	w := world.Create(WorldDimensions.Length, 17, WorldDimensions.Width, WorldSeeder)

	t2_wc := time.Now()
	w.Image(WorldFilePath)

	t3_wc := time.Now()

	fmt.Println()
	w.PrintDebugInfos()
	fmt.Println(world.BiomesLegend())
	t4_wc := time.Now()

	fmt.Println()
	if EnableTest {
		fmt.Printf("Test NoiseMaps (%vx%v) created and saved at %v in %v\n", 2048, 2048, "./noiseMap/testImages/", t1_wc.Sub(t0_wc))
	}
	fmt.Printf("World %v (%vx%v) created in %v and saved at %v in %v (printed in %v)\n", WorldSeeder.GetSeed().Token, w.Data.Length, w.Data.Width, t2_wc.Sub(t1_wc), fmt.Sprintf(WorldFilePath, world.PathFilesNumber), t3_wc.Sub(t2_wc), t4_wc.Sub(t3_wc))
}

func testingNoiseMaps(width, height int) {
	// Create a new white noise map and generate an image file
	t0_wnm := time.Now()
	whiteNoise := noiseMap.NewWhiteNoiseMap(width, height, WorldSeeder)
	t1_wnm := time.Now()
	whiteNoise.Image("./noiseMap/testImages/whiteNoise_test%v.png")
	fmt.Printf("White Noise elapsed time: %v\n", t1_wnm.Sub(t0_wnm))

	// Create a new perlin noise map and generate an image file
	frequency, amplitude, lacunarity, persistence := 0.005, 1.0, 2.0, 0.5
	octaves := 4
	t0_pnm := time.Now()
	perlinNoise := noiseMap.NewPerlinNoiseMap(width, height, frequency, amplitude, lacunarity, persistence, octaves, WorldSeeder)
	t1_pnm := time.Now()
	perlinNoise.Image("./noiseMap/testImages/perlinNoise_test%v.png")
	fmt.Printf("Perlin Noise elapsed time: %v\n", t1_pnm.Sub(t0_pnm))
}

/* Noise Map Tested Parameter */
// White Noise Map
// 		seed = 0123456789
//		width = 2048, height = 2048
//		avg. elapsed time: 0.080s
// Perlin Noise Map
//		seed = 0123456789
// 		width = 2048, height = 2048, octaves = 4, gridSize = 256.0
// 		avg. elapsed time: 4.530s
