// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package jbd

func grow(slice []byte, newCap int) []byte {
	if cap(slice) >= newCap {
		return slice
	}

	newSlice := make([]byte, len(slice), newCap)
	copy(newSlice, slice)

	return newSlice
}
