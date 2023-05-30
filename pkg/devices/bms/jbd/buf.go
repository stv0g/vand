// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package jbd

import "encoding/binary"

type Buf []byte

func (buf Buf) Uint16(off int) uint16 {
	return binary.BigEndian.Uint16(buf[off:])
}

func (buf Buf) Int16(off int) int16 {
	return int16(binary.BigEndian.Uint16(buf[off:]))
}

func (buf Buf) Uint32(off int) uint32 {
	return binary.BigEndian.Uint32(buf[off:])
}

func (buf Buf) Bit(off, bit int) bool {
	return buf[off]&(1<<bit) != 0
}
