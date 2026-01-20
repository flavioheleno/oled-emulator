package device

import "fmt"

// MemoryHelper provides utilities for memory operations
type MemoryHelper struct {
	width       int
	height      int
	pixelFormat PixelFormat
	colOffset   int
}

// NewMemoryHelper creates a new memory helper
func NewMemoryHelper(width, height int, pixelFormat PixelFormat, colOffset int) *MemoryHelper {
	return &MemoryHelper{
		width:       width,
		height:      height,
		pixelFormat: pixelFormat,
		colOffset:   colOffset,
	}
}

// PixelToByteOffset converts pixel coordinates to VRAM byte offset for HorizontalNibble format
func (mh *MemoryHelper) PixelToByteOffsetNibble(x, y int) (int, int, error) {
	if x < 0 || x >= mh.width || y < 0 || y >= mh.height {
		return 0, 0, fmt.Errorf("pixel out of bounds: (%d, %d)", x, y)
	}

	// For SSD1322 with HorizontalNibble format (2 pixels per byte)
	// Each row has 480 columns internally (even if display is 256 wide)
	columns := 480
	byteOffset := (y * columns + x + mh.colOffset) / 2
	nibbleIndex := (x + mh.colOffset) % 2

	return byteOffset, nibbleIndex, nil
}

// SetPixelNibble sets a pixel in HorizontalNibble format (4-bit gray)
func (mh *MemoryHelper) SetPixelNibble(vram []byte, x, y int, color byte) error {
	byteOffset, nibbleIndex, err := mh.PixelToByteOffsetNibble(x, y)
	if err != nil {
		return err
	}

	if byteOffset >= len(vram) {
		return fmt.Errorf("VRAM offset out of bounds: %d", byteOffset)
	}

	// Ensure color is 4-bit
	color = color & 0x0F

	if nibbleIndex == 0 {
		// Lower nibble
		vram[byteOffset] = (vram[byteOffset] & 0xF0) | color
	} else {
		// Upper nibble
		vram[byteOffset] = (vram[byteOffset] & 0x0F) | (color << 4)
	}

	return nil
}

// GetPixelNibble reads a pixel in HorizontalNibble format
func (mh *MemoryHelper) GetPixelNibble(vram []byte, x, y int) (byte, error) {
	byteOffset, nibbleIndex, err := mh.PixelToByteOffsetNibble(x, y)
	if err != nil {
		return 0, err
	}

	if byteOffset >= len(vram) {
		return 0, fmt.Errorf("VRAM offset out of bounds: %d", byteOffset)
	}

	if nibbleIndex == 0 {
		return vram[byteOffset] & 0x0F, nil
	}
	return (vram[byteOffset] >> 4) & 0x0F, nil
}

// PixelToByteOffsetVertical converts pixel coordinates to VRAM byte offset for VerticalByte format
func (mh *MemoryHelper) PixelToByteOffsetVertical(x, y int) (int, int, error) {
	if x < 0 || x >= mh.width || y < 0 || y >= mh.height {
		return 0, 0, fmt.Errorf("pixel out of bounds: (%d, %d)", x, y)
	}

	// Vertical packing: 8 pixels per byte, stacked vertically (SSD1306 style)
	byteOffset := x*((mh.height+7)/8) + y/8
	bitOffset := y % 8

	return byteOffset, bitOffset, nil
}

// SetPixelVertical sets a pixel in VerticalByte format (1-bit mono)
func (mh *MemoryHelper) SetPixelVertical(vram []byte, x, y int, color byte) error {
	byteOffset, bitOffset, err := mh.PixelToByteOffsetVertical(x, y)
	if err != nil {
		return err
	}

	if byteOffset >= len(vram) {
		return fmt.Errorf("VRAM offset out of bounds: %d", byteOffset)
	}

	if color > 0 {
		vram[byteOffset] |= (1 << bitOffset)
	} else {
		vram[byteOffset] &= ^(1 << bitOffset)
	}

	return nil
}

// GetPixelVertical reads a pixel in VerticalByte format
func (mh *MemoryHelper) GetPixelVertical(vram []byte, x, y int) (byte, error) {
	byteOffset, bitOffset, err := mh.PixelToByteOffsetVertical(x, y)
	if err != nil {
		return 0, err
	}

	if byteOffset >= len(vram) {
		return 0, fmt.Errorf("VRAM offset out of bounds: %d", byteOffset)
	}

	if (vram[byteOffset] & (1 << bitOffset)) != 0 {
		return 1, nil
	}
	return 0, nil
}

// SetPixelRGB888 sets a pixel in RGB888 format (24-bit color)
func (mh *MemoryHelper) SetPixelRGB888(vram []byte, x, y int, r, g, b byte) error {
	if x < 0 || x >= mh.width || y < 0 || y >= mh.height {
		return fmt.Errorf("pixel out of bounds: (%d, %d)", x, y)
	}

	offset := (y*mh.width + x) * 3
	if offset+2 >= len(vram) {
		return fmt.Errorf("VRAM offset out of bounds")
	}

	vram[offset] = r
	vram[offset+1] = g
	vram[offset+2] = b

	return nil
}

// GetPixelRGB888 reads a pixel in RGB888 format
func (mh *MemoryHelper) GetPixelRGB888(vram []byte, x, y int) (byte, byte, byte, error) {
	if x < 0 || x >= mh.width || y < 0 || y >= mh.height {
		return 0, 0, 0, fmt.Errorf("pixel out of bounds: (%d, %d)", x, y)
	}

	offset := (y*mh.width + x) * 3
	if offset+2 >= len(vram) {
		return 0, 0, 0, fmt.Errorf("VRAM offset out of bounds")
	}

	return vram[offset], vram[offset+1], vram[offset+2], nil
}

// FillRegionNibble fills a rectangular region with a color in HorizontalNibble format
func (mh *MemoryHelper) FillRegionNibble(vram []byte, x0, y0, x1, y1 int, color byte) error {
	color = color & 0x0F

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if err := mh.SetPixelNibble(vram, x, y, color); err != nil {
				return err
			}
		}
	}

	return nil
}

// FillRegionVertical fills a rectangular region with a color in VerticalByte format
func (mh *MemoryHelper) FillRegionVertical(vram []byte, x0, y0, x1, y1 int, color byte) error {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if err := mh.SetPixelVertical(vram, x, y, color); err != nil {
				return err
			}
		}
	}

	return nil
}

// ExtractRegionNibble extracts a rectangular region as a new buffer
func (mh *MemoryHelper) ExtractRegionNibble(vram []byte, x0, y0, x1, y1 int) ([]byte, error) {
	width := x1 - x0 + 1
	height := y1 - y0 + 1

	// Create new buffer for extracted region
	extracted := make([]byte, (width*height)/2+1)

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			pixel, err := mh.GetPixelNibble(vram, x, y)
			if err != nil {
				return nil, err
			}

			// Write to extracted buffer
			relX := x - x0
			relY := y - y0
			offset := (relY * width + relX) / 2
			nibbleIndex := (relX) % 2

			if nibbleIndex == 0 {
				extracted[offset] = (extracted[offset] & 0xF0) | pixel
			} else {
				extracted[offset] = (extracted[offset] & 0x0F) | (pixel << 4)
			}
		}
	}

	return extracted, nil
}
