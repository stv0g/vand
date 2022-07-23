package jbd

func grow(slice []byte, newCap int) []byte {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]byte, len(slice), newCap)
	copy(newSlice, slice)

	return newSlice
}
