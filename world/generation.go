package world

import "github.com/MoraGames/warld/seed"

func GenerateWorldMatrix(seed *seed.Seed) *WorldMatrix {
	// Create a new world map
	worldMatrix := make(WorldMatrix, 5)
	for i := range worldMatrix {
		worldMatrix[i] = make([]Tile, 5)
		for j := range worldMatrix[i] {
			worldMatrix[i][j] = NewTile()
		}
	}

	// Generate the world map
	// TODO: Implement the world map generation

	// Return the world map
	return &worldMatrix
}

func ZoomWorld(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Zoom the world map (each tile becomes 2x2 tiles)
	// TODO: Implement ZoomWorld

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func AddIslands(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Add islands to the world map (randomly located, the number of islands is based on the map size)
	// TODO: Implement AddIslands

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func RemoveTooMuchOceans(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Remove too much oceans from the world map, replacing them with land
	// TODO: Implement RemoveTooMuchOceans

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func AddTemperatures(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Add temperatures to the world map
	// TODO: Implement AddTemperatures

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func AddPrecipitations(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Add precipitations to the world map
	// TODO: Implement AddPrecipitations

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func AddAltitudes(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Add altitudes to the world map
	// TODO: Implement AddAltitudes

	// Update the world map
	*worldMatrixPtr = worldMatrix
}

func AddVaraints(worldMatrixPtr *WorldMatrix) {
	// Retrieve the world map
	worldMatrix := *worldMatrixPtr

	// Add variants to the world map
	// TODO: Implement AddVaraints

	// Update the world map
	*worldMatrixPtr = worldMatrix
}
