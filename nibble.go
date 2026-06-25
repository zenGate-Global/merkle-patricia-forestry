// Copyright 2025 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mpf

import (
	"encoding/hex"
	"strings"
)

type Nibble byte

func (n Nibble) String() string {
	return hex.EncodeToString([]byte{byte(n)})[1:]
}

// bytesToNibbles splits a series of bytes into a series of nibbles
func bytesToNibbles(data []byte) []Nibble {
	ret := make([]Nibble, 0, len(data)*2)
	for _, dataByte := range data {
		tmpNibbles := byteToNibbles(dataByte)
		ret = append(ret, tmpNibbles...)
	}
	return ret
}

// byteToNibbles splits a byte into two bytes representing the upper and lower 4 bits of the original byte.
// The value 0xab would be returned as [0x0a, 0x0b]
func byteToNibbles(data byte) []Nibble {
	// Split byte into two bytes representing the upper and lower 4 bits
	return []Nibble{
		Nibble(data >> 4),
		Nibble(data & 0xf),
	}
}

// nibblesToBytes packs a series of Nibbles into a byte slice, two nibbles per
// byte (high nibble first). This is the inverse of bytesToNibbles and is valid
// for an even number of nibbles (e.g. a full hash path). A trailing odd nibble
// is packed into the low half of a final byte (zero high half) instead of
// indexing out of bounds and panicking.
func nibblesToBytes(data []Nibble) []byte {
	ret := make([]byte, 0, (len(data)+1)/2)
	i := 0
	for ; i+1 < len(data); i += 2 {
		tmpByte := byte(data[i]<<4) + byte(data[i+1])
		ret = append(ret, tmpByte)
	}
	if i < len(data) {
		// Odd trailing nibble: place it in the low half of a final byte.
		ret = append(ret, byte(data[i]))
	}
	return ret
}

// nibblesToExpandedBytes encodes a series of Nibbles as a byte slice with one
// nibble per byte (the low half of each byte). This matches the on-chain Aiken
// merkle-patricia-forestry encoding of a Fork neighbour's prefix, where the
// prefix is consumed nibble-by-nibble (see helpers.nibbles and do_fork's
// combine(neighbor.prefix, neighbor.root)) and the reference TS library's
// nibbles(prefix) buffer. Works for any length, including odd.
func nibblesToExpandedBytes(data []Nibble) []byte {
	ret := make([]byte, len(data))
	for i, n := range data {
		ret[i] = byte(n)
	}
	return ret
}

// nibblesToHexString converts a series of Nibbles into a hex string representing those nibbles.
func nibblesToHexString(data []Nibble) string {
	var sb strings.Builder
	for _, nibble := range data {
		sb.WriteString(nibble.String())
	}
	return sb.String()
}

// keyToPath converts an arbitrary key to the sequence of Nibbles representing the path to the value
func keyToPath(key []byte) []Nibble {
	keyHash := HashValue(key)
	keyHashNibbles := bytesToNibbles(keyHash.Bytes())
	return keyHashNibbles
}
