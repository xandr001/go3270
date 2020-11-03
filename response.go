// This file is part of https://github.com/racingmars/go3270/
// Copyright 2020 by Matthew R. Wilson, licensed under the MIT license. See
// LICENSE in the project root for license information.

package go3270

import (
	"bytes"
	"net"
)

// Response encapsulates data received from a 3270 client in response to the
// previously sent screen.
type Response struct {
	// Which Action ID key did the user press?
	AID AID

	// Row the cursor was on (0-based).
	Row int

	// Column the cursor was on (0-based).
	Col int

	// Field values.
	Values map[string]string
}

// AID is an Action ID character.
type AID byte

const (
	AIDNone  AID = 0x60
	AIDEnter AID = 0x7D
	AIDPF1   AID = 0xF1
	AIDPF2   AID = 0xF2
	AIDPF3   AID = 0xF3
	AIDPF4   AID = 0xF4
	AIDPF5   AID = 0xF5
	AIDPF6   AID = 0xF6
	AIDPF7   AID = 0xF7
	AIDPF8   AID = 0xF8
	AIDPF9   AID = 0xF9
	AIDPF10  AID = 0x7A
	AIDPF11  AID = 0x7B
	AIDPF12  AID = 0x7C
	AIDPF13  AID = 0xC1
	AIDPF14  AID = 0xC2
	AIDPF15  AID = 0xC3
	AIDPF16  AID = 0xC4
	AIDPF17  AID = 0xC5
	AIDPF18  AID = 0xC6
	AIDPF19  AID = 0xC7
	AIDPF20  AID = 0xC8
	AIDPF21  AID = 0xC9
	AIDPF22  AID = 0x4A
	AIDPF23  AID = 0x4B
	AIDPF24  AID = 0x4C
	AIDPA1   AID = 0x6C
	AIDPA2   AID = 0x6E
	AIDPA3   AID = 0x6B
	AIDClear AID = 0x6D
)

func readResponse(c net.Conn) (Response, error) {
	var r Response
	aid, err := readAID(c)
	if err != nil {
		return r, err
	}
	r.AID = aid

	row, col, _, err := readPosition(c)
	if err != nil {
		return r, err
	}
	r.Col = col
	r.Row = row

	if err = readFields(c); err != nil {
		return r, err
	}

	return r, nil
}

func readAID(c net.Conn) (AID, error) {
	buf := make([]byte, 1)
	for {
		_, err := c.Read(buf)
		if err != nil {
			return AIDNone, err
		}
		b := buf[0]
		if (b == 0x60) || (b >= 0x6b && b <= 0x6e) ||
			(b >= 0x7a && b <= 0x7d) || (b >= 0x4a && b <= 0x4c) ||
			(b >= 0xf1 && b <= 0xf9) || (b >= 0xc1 && b <= 0xc9) {
			// We found a valid AID
			debugf("Got AID byte: %x\n", b)
			return AID(b), nil
		}
		// Consume non-AID bytes continuing loop
		debugf("Got non-AID byte: %x\n", b)
	}
}

func readPosition(c net.Conn) (row, col, addr int, err error) {
	buf := make([]byte, 1)
	raw := make([]byte, 2)

	// Read two bytes
	for i := 0; i < 2; i++ {
		if _, err := c.Read(buf); err != nil {
			return 0, 0, 0, err
		}
		raw[i] = buf[0]
	}

	// Decode the raw position
	addr = decodeBufAddr([2]byte{raw[0], raw[1]})
	row = addr % 80
	col = (addr - row) / 80

	debugf("Got position bytes %02x %02x, decoded to %d\n", raw[0], raw[1],
		addr)

	return row, col, addr, nil
}

func readFields(c net.Conn) error {
	buf := make([]byte, 1)
	var infield bool
	var fieldpos int
	var fieldval bytes.Buffer
	var err error

	// consume bytes until we get 0xffef
	for {
		// Read a byte
		if _, err = c.Read(buf); err != nil {
			return err
		}

		// Check for end of data stream (0xffef)
		if buf[0] == 0xff {
			// Finish the current field
			if infield {
				// TODO
				debugf("Field %d: %s\n", fieldpos, e2a(fieldval.Bytes()))
			}

			// consume the next byte, which is probably 0xef
			if _, err = c.Read(buf); err != nil {
				return err
			}
			return nil
		}

		// No? Check for start-of-field
		if buf[0] == 0x11 {
			// Finish the previous field, if necessary
			if infield {
				// TODO
				debugf("Field %d: %s\n", fieldpos, e2a(fieldval.Bytes()))
			}
			// Start a new field
			infield = true
			fieldval = bytes.Buffer{}
			fieldpos = 0

			if _, _, fieldpos, err = readPosition(c); err != nil {
				return err
			}
			continue
		}

		// Consume all other bytes as field contents if we're in a field
		if !infield {
			debugf("Got unexpected byte while processing fields: %02x\n", buf[0])
			continue
		}
		fieldval.WriteByte(buf[0])

	}
}

// decodeBufAddr decodes a raw 2-byte encoded buffer address and returns the
// integer value of the address (i.e. 0-1919)
func decodeBufAddr(raw [2]byte) int {
	hi := decodes[raw[0]] << 6
	lo := decodes[raw[1]]
	return hi | lo
}
