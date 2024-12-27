package hash

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"strconv"
)

func Compute(dataBytes []byte) string {
	var H = [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}
	n := uint64(len(dataBytes)) * 8
	dataBytes = append(dataBytes, byte(128))

	for len(dataBytes)%64 != 56 {
		dataBytes = append(dataBytes, 0)
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)
	dataBytes = append(dataBytes, buf...)
	numBlocks := len(dataBytes) / 64

	for i := 0; i < numBlocks; i++ {
		block := dataBytes[i*64 : (i+1)*64]

		words := make([]uint32, 64)
		for j := 0; j < 16; j++ {
			word := block[4*j : 4*(j+1)]
			words[j] = binary.BigEndian.Uint32(word)
		}
		for j := 16; j < 64; j++ {
			words[j] = sigma0(words[j-15]) + sigma1(words[j-2]) + words[j-7] + words[j-16]
		}

		var (
			a = H[0]
			b = H[1]
			c = H[2]
			d = H[3]
			e = H[4]
			f = H[5]
			g = H[6]
			h = H[7]
		)

		for t := 0; t < 64; t++ {
			T1 := h + sigmaL1(e) + choose(e, f, g) + kArray[t] + words[t]
			T2 := sigmaL0(a) + major(a, b, c)

			h = g
			g = f
			f = e
			e = d + T1
			d = c
			c = b
			b = a
			a = T1 + T2
		}
		H[0] = H[0] + a
		H[1] = H[1] + b
		H[2] = H[2] + c
		H[3] = H[3] + d
		H[4] = H[4] + e
		H[5] = H[5] + f
		H[6] = H[6] + g
		H[7] = H[7] + h
	}
	hash := ""
	for i := 0; i < 8; i++ {
		hash += fmt.Sprintf("%08s", strconv.FormatUint(uint64(H[i]), 16))
	}
	return hash
}

var kArray = [64]uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
}

func choose(x, y, z uint32) uint32 {
	return xor(x&y, (^x)&z)
}

func major(x, y, z uint32) uint32 {
	return xor(xor(x&y, x&z), y&z)
}

func xor(a, b uint32) uint32 {
	return (a & ^b) | (^a & b)
}

func sigma0(x uint32) uint32 {
	a := bits.RotateLeft32(x, -7)
	b := bits.RotateLeft32(x, -18)
	c := x >> 3
	return xor(xor(a, b), c)
}

func sigma1(x uint32) uint32 {
	a := bits.RotateLeft32(x, -17)
	b := bits.RotateLeft32(x, -19)
	c := x >> 10
	return xor(xor(a, b), c)
}

func sigmaL0(x uint32) uint32 {
	a := bits.RotateLeft32(x, -2)
	b := bits.RotateLeft32(x, -13)
	c := bits.RotateLeft32(x, -22)
	return xor(xor(a, b), c)
}

func sigmaL1(x uint32) uint32 {
	a := bits.RotateLeft32(x, -6)
	b := bits.RotateLeft32(x, -11)
	c := bits.RotateLeft32(x, -25)
	return xor(xor(a, b), c)
}
