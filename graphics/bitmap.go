package graphics

import (
	"fmt"
)

// BitmapFont provides a simple bitmap-based font for monospace text
type BitmapFont struct {
	glyphs  map[rune]GlyphData
	width   int
	height  int
	advance int
}

// NewBitmapFont creates a new bitmap font
func NewBitmapFont(width, height, advance int) *BitmapFont {
	return &BitmapFont{
		glyphs:  make(map[rune]GlyphData),
		width:   width,
		height:  height,
		advance: advance,
	}
}

// Height returns the font height
func (bf *BitmapFont) Height() int {
	return bf.height
}

// AddGlyph adds a glyph to the font
func (bf *BitmapFont) AddGlyph(ch rune, data GlyphData) {
	bf.glyphs[ch] = data
}

// DrawString draws text at the specified position
func (bf *BitmapFont) DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error) {
	currentX := x
	color = color & 0x0F

	for _, ch := range text {
		glyph, ok := bf.glyphs[ch]
		if !ok {
			// Use space character as fallback
			if ch == ' ' {
				currentX += bf.advance
				continue
			}
			// Try to find a replacement glyph
			glyph, ok = bf.glyphs[' ']
			if !ok {
				currentX += bf.advance
				continue
			}
		}

		// Draw the glyph
		if err := bf.drawGlyph(fb, currentX, y, glyph, color); err != nil {
			return 0, err
		}

		currentX += bf.advance
	}

	return currentX - x, nil
}

// MeasureString returns the width and height of text
func (bf *BitmapFont) MeasureString(text string) (width, height int, err error) {
	return len(text) * bf.advance, bf.height, nil
}

// GetGlyph returns glyph data for a character
func (bf *BitmapFont) GetGlyph(ch rune) (GlyphData, error) {
	glyph, ok := bf.glyphs[ch]
	if !ok {
		return GlyphData{}, fmt.Errorf("glyph not found: %c", ch)
	}
	return glyph, nil
}

// drawGlyph draws a single glyph to the framebuffer
func (bf *BitmapFont) drawGlyph(fb *FrameBuffer, x, y int, glyph GlyphData, color byte) error {
	if glyph.Width <= 0 || glyph.Height <= 0 || len(glyph.Data) == 0 {
		return nil // Empty glyph
	}

	byteIndex := 0

	for glyphY := 0; glyphY < glyph.Height; glyphY++ {
		bitIndex := 0

		for glyphX := 0; glyphX < glyph.Width; glyphX++ {
			// Make sure we don't go out of bounds
			if byteIndex >= len(glyph.Data) {
				return nil
			}

			// Check if current bit is set
			bitMask := (1 << (7 - bitIndex))
			isSet := (glyph.Data[byteIndex] & byte(bitMask)) != 0

			if isSet {
				// Draw pixel to framebuffer
				screenX := x + glyphX + glyph.BearingX
				screenY := y + glyphY + glyph.BearingY

				if screenX >= 0 && screenY >= 0 {
					fb.SetPixel(screenX, screenY, color)
				}
			}

			bitIndex++
			if bitIndex == 8 {
				bitIndex = 0
				byteIndex++
			}
		}

		// Move to next row, aligned to byte boundary
		if bitIndex != 0 {
			byteIndex++
		}
	}

	return nil
}

// DefaultBitmapFont creates a default monospace bitmap font with ASCII characters
func DefaultBitmapFont() *BitmapFont {
	bf := NewBitmapFont(5, 7, 6)

	// Simple ASCII bitmap glyphs (5x7 pixel font)
	// These are basic 5x7 character representations
	for ch := rune(32); ch <= rune(126); ch++ {
		bf.AddGlyph(ch, createASCIIGlyph(ch))
	}

	return bf
}

// createASCIIGlyph creates a simple ASCII glyph
func createASCIIGlyph(ch rune) GlyphData {
	// This is a simplified implementation
	// In a real system, you would have pre-rendered glyphs
	width := 5
	height := 7
	bytesPerRow := (width + 7) / 8

	// Create basic glyphs for common characters
	var data []byte

	switch ch {
	case ' ':
		// Space - empty
		data = make([]byte, bytesPerRow*height)

	case 'A':
		// Letter A (5 bits wide, 7 bits tall)
		data = []byte{
			0b01110000,
			0b10001000,
			0b10001000,
			0b11111000,
			0b10001000,
			0b10001000,
			0b10001000,
		}

	case 'B':
		// Letter B
		data = []byte{
			0b11110000,
			0b10001000,
			0b10001000,
			0b11100000,
			0b10001000,
			0b10001000,
			0b11110000,
		}

	case 'H':
		// Letter H
		data = []byte{
			0b10001000,
			0b10001000,
			0b10001000,
			0b11111000,
			0b10001000,
			0b10001000,
			0b10001000,
		}

	case 'O':
		// Letter O
		data = []byte{
			0b01110000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b01110000,
		}

	case '0':
		// Digit 0
		data = []byte{
			0b01110000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b10001000,
			0b01110000,
		}

	case '1':
		// Digit 1
		data = []byte{
			0b00100000,
			0b01100000,
			0b00100000,
			0b00100000,
			0b00100000,
			0b00100000,
			0b01110000,
		}

	default:
		// Default character - simple block
		data = make([]byte, bytesPerRow*height)
		for i := 0; i < len(data); i++ {
			data[i] = 0x78 // 0b01111000 (5 bits set)
		}
	}

	return GlyphData{
		Width:    width,
		Height:   height,
		AdvanceX: 6,
		BearingX: 0,
		BearingY: 0,
		Data:     data,
	}
}
