// Package imagehash computes a perceptual hash (pHash-style) for a
// listing image. The intent is to detect listing-image reuse across
// scam networks, without depending on a paid reverse-image-search
// vendor.
package imagehash

import (
	"image"
	"image/color"
	"image/draw"
	// Side-effect imports for image format decoding.
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/bits"
)

// Size is the side of the working thumbnail used for the hash.
const Size = 32

// Hash returns a 64-bit perceptual hash of the supplied image. Two
// images whose hashes differ by < 6 bits of Hamming distance are
// likely duplicates.
func Hash(r io.Reader) (uint64, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}
	thumb := image.NewGray(image.Rect(0, 0, Size, Size))
	draw.Draw(thumb, thumb.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(thumb, thumb.Bounds(), src, src.Bounds().Min, draw.Over)
	avg := uint64(0)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			avg += uint64(thumb.GrayAt(x*Size/8, y*Size/8).Y)
		}
	}
	avg /= 64
	var h uint64
	bit := uint64(1)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			v := uint64(thumb.GrayAt(x*Size/8, y*Size/8).Y)
			if v >= avg {
				h |= bit
			}
			bit <<= 1
		}
	}
	return h, nil
}

// Distance is the Hamming distance between two hashes.
func Distance(a, b uint64) int { return bits.OnesCount64(a ^ b) }
