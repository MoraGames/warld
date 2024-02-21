package coord

type Coord struct {
	Z int
	X int
}

func NewCoord(z, x int) Coord {
	return Coord{Z: z, X: x}
}

func (cc Coord) GetAlignedNeighbors(maxZ, maxX int) CoordSlice {
	// Create a temporary slice of coordinates that are aligned with the current coordinate
	TempAlignedNeighbors := []Coord{
		NewCoord(cc.Z-1, cc.X), // Up
		NewCoord(cc.Z, cc.X+1), // Right
		NewCoord(cc.Z+1, cc.X), // Down
		NewCoord(cc.Z, cc.X-1), // Left
	}

	// Remove invalid coordinates
	AlignedNeighbors := make([]Coord, 0)
	for i := 0; i < len(TempAlignedNeighbors); i++ {
		if TempAlignedNeighbors[i].Z >= 0 && TempAlignedNeighbors[i].Z <= maxZ && TempAlignedNeighbors[i].X >= 0 && TempAlignedNeighbors[i].X <= maxX {
			AlignedNeighbors = append(AlignedNeighbors, TempAlignedNeighbors[i])
		}
	}

	return AlignedNeighbors
}

func (cc Coord) GetDiagonalNeighbors(maxZ, maxX int) CoordSlice {
	// Create a temporary slice of coordinates that are diagonal from the current coordinate
	TempDiagonalNeighbors := []Coord{
		NewCoord(cc.Z-1, cc.X-1), // Up-Left
		NewCoord(cc.Z-1, cc.X+1), // Up-Right
		NewCoord(cc.Z+1, cc.X+1), // Down-Right
		NewCoord(cc.Z+1, cc.X-1), // Down-Left
	}

	// Remove invalid coordinates
	DiagonalNeighbors := make([]Coord, 0)
	for i := 0; i < len(TempDiagonalNeighbors); i++ {
		if TempDiagonalNeighbors[i].Z >= 0 && TempDiagonalNeighbors[i].Z <= maxZ && TempDiagonalNeighbors[i].X >= 0 && TempDiagonalNeighbors[i].X <= maxX {
			DiagonalNeighbors = append(DiagonalNeighbors, TempDiagonalNeighbors[i])
		}
	}

	return DiagonalNeighbors
}

func (cc Coord) GetLargeRingNeighbors(maxZ, maxX int) CoordSlice {
	// Create a temporary slice of coordinates that have a distance of 2 from the current coordinate
	TempLargeRingNeighbors := []Coord{
		NewCoord(cc.Z-2, cc.X-2), // Up-Left
		NewCoord(cc.Z-2, cc.X-1), // Up-Up-Left
		NewCoord(cc.Z-2, cc.X),   // Up
		NewCoord(cc.Z-2, cc.X+1), // Up-Up-Right
		NewCoord(cc.Z-2, cc.X+2), // Up-Right
		NewCoord(cc.Z-1, cc.X+2), // Up-Right-Right
		NewCoord(cc.Z, cc.X+2),   // Right
		NewCoord(cc.Z+1, cc.X+2), // Down-Right-Right
		NewCoord(cc.Z+2, cc.X+2), // Down-Right
		NewCoord(cc.Z+2, cc.X+1), // Down-Down-Right
		NewCoord(cc.Z+2, cc.X),   // Down
		NewCoord(cc.Z+2, cc.X-1), // Down-Down-Left
		NewCoord(cc.Z+2, cc.X-2), // Down-Left
		NewCoord(cc.Z+1, cc.X-2), // Down-Left-Left
		NewCoord(cc.Z, cc.X-2),   // Left
		NewCoord(cc.Z-1, cc.X-2), // Up-Left-Left
	}

	// Remove invalid coordinates
	LargeRingNeighbors := make([]Coord, 0)
	for i := 0; i < len(TempLargeRingNeighbors); i++ {
		if TempLargeRingNeighbors[i].Z >= 0 && TempLargeRingNeighbors[i].Z <= maxZ && TempLargeRingNeighbors[i].X >= 0 && TempLargeRingNeighbors[i].X <= maxX {
			LargeRingNeighbors = append(LargeRingNeighbors, TempLargeRingNeighbors[i])
		}
	}

	return LargeRingNeighbors
}
