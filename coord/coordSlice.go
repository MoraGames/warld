package coord

type CoordSlice []Coord

func NewCoordSlice(coords ...Coord) CoordSlice {
	return coords
}

func ConcatenateCoordSlices(slices ...CoordSlice) CoordSlice {
	var result CoordSlice
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func (cs CoordSlice) Contains(c Coord) bool {
	for _, coord := range cs {
		//fmt.Printf("[DEBUG] coord = %v | c = %v | equal? = %v\n", coord, c, coord.Z == c.Z && coord.X == c.X)
		if coord.Z == c.Z && coord.X == c.X {
			return true
		}
	}
	return false
}
