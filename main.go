package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MoraGames/warld/noiseMap"
	"github.com/MoraGames/warld/seed"
	"github.com/MoraGames/warld/world"
)

var (
	WorldSeeder *seed.Seeder
)

func init() {
	// Obtain the world seed
	var seedStr string
	fmt.Print("Enter the world seed: ")
	Scan(&seedStr)

	// Create the random number generator based on the seed
	seedPtr, err := seed.GetSeed(seedStr)
	if err != nil {
		log.Panicln(err)
	}
	WorldSeeder = seed.NewSeeder(seedPtr)
}

func main() {
	fmt.Printf("World seed: %v | World rand: %v\n", WorldSeeder.GetSeed().Token, WorldSeeder.Random)
	// TODO: Implement main()

	//testingNoiseMaps()
	t0_wc := time.Now()
	w := world.Create(512, 12, 512, WorldSeeder)
	t1_wc := time.Now()
	fmt.Println(w.String(true))
	t2_wc := time.Now()
	fmt.Printf("\nworld %v (%vx%v) created in %v (and printed in %v)\n", WorldSeeder.GetSeed().Token, w.Data.Length, w.Data.Width, t1_wc.Sub(t0_wc), t2_wc.Sub(t1_wc))
}

func testingNoiseMaps() {
	// Create a new white noise map and generate an image file
	width, height := 2048, 2048
	t0_wnm := time.Now()
	whiteNoise := noiseMap.NewWhiteNoiseMap(width, height, WorldSeeder)
	t1_wnm := time.Now()
	whiteNoise.Image("./noiseMap/testImages/whiteNoise_test01.png")
	fmt.Printf("White Noise elapsed time: %v\n", t1_wnm.Sub(t0_wnm))

	// Create a new perlin noise map and generate an image file
	octaves := 4
	gridSize := 256.0
	t0_pnm := time.Now()
	perlinNoise := noiseMap.NewPerlinNoiseMap(width, height, octaves, gridSize, WorldSeeder)
	t1_pnm := time.Now()
	perlinNoise.Image("./noiseMap/testImages/perlinNoise_test01.png")
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
