package device

import "fmt"

// PixelFormat defines how pixels are packed in VRAM
type PixelFormat int

const (
	// HorizontalNibble: 2 pixels per byte, packed horizontally (SSD1322 native)
	HorizontalNibble PixelFormat = iota
	// VerticalByte: 8 pixels per byte, packed vertically (SSD1306 style)
	VerticalByte
	// RGB888: 24-bit RGB color
	RGB888
)

// Config holds device configuration
type Config struct {
	Width        int           // Display width in pixels
	Height       int           // Display height in pixels
	ColorDepth   int           // Bits per pixel: 1, 4, 8, 24
	PixelFormat  PixelFormat   // How pixels are packed in memory
	ColumnOffset int           // Offset for VRAM column (e.g., 28 for SSD1322)
	InitCommands []byte        // Custom initialization sequence
}

// Device defines the interface for display emulation
type Device interface {
	// ProcessCommand processes a device command
	ProcessCommand(cmd byte, data []byte) error

	// GetFrameBuffer returns the current VRAM contents
	GetFrameBuffer() []byte

	// GetDirtyRegion returns the bounding box of changed pixels
	// Returns (x0, y0, x1, y1) or (-1, -1, -1, -1) if no changes
	GetDirtyRegion() (int, int, int, int)

	// ClearDirtyRegion resets the dirty tracking
	ClearDirtyRegion()

	// Configuration getters
	Width() int
	Height() int
	ColorDepth() int
	PixelFormat() PixelFormat

	// Reset performs a hardware reset
	Reset() error

	// SetPixel sets a pixel directly (for testing/high-level API)
	SetPixel(x, y int, color byte) error

	// GetPixel reads a pixel value
	GetPixel(x, y int) (byte, error)
}

// BaseDevice provides common functionality for device implementations
type BaseDevice struct {
	config      Config
	vram        []byte
	dirtyX0     int
	dirtyY0     int
	dirtyX1     int
	dirtyY1     int
	hasDirty    bool
}

// NewBaseDevice creates a new base device
func NewBaseDevice(config Config) *BaseDevice {
	// Validate configuration
	if config.Width <= 0 || config.Height <= 0 {
		panic(fmt.Sprintf("invalid display dimensions: %dx%d", config.Width, config.Height))
	}

	bd := &BaseDevice{
		config:   config,
		dirtyX0:  -1,
		dirtyY0:  -1,
		dirtyX1:  -1,
		dirtyY1:  -1,
		hasDirty: false,
	}

	// Allocate VRAM based on pixel format
	bd.vram = bd.allocateVRAM()

	return bd
}

// allocateVRAM calculates and allocates VRAM
func (bd *BaseDevice) allocateVRAM() []byte {
	var byteCount int

	switch bd.config.PixelFormat {
	case HorizontalNibble:
		// 2 pixels per byte (4 bits each)
		// Include column offset for SSD1322 (480 columns internal)
		columns := 480 // SSD1322 has 480 columns internally
		rows := bd.config.Height
		byteCount = (columns * rows) / 2
	case VerticalByte:
		// 8 pixels per byte, packed vertically
		byteCount = bd.config.Width * ((bd.config.Height + 7) / 8)
	case RGB888:
		// 24-bit color (3 bytes per pixel)
		byteCount = bd.config.Width * bd.config.Height * 3
	default:
		panic("unsupported pixel format")
	}

	return make([]byte, byteCount)
}

// GetFrameBuffer returns the VRAM
func (bd *BaseDevice) GetFrameBuffer() []byte {
	return bd.vram
}

// GetDirtyRegion returns the bounding box of dirty pixels
func (bd *BaseDevice) GetDirtyRegion() (int, int, int, int) {
	if !bd.hasDirty {
		return -1, -1, -1, -1
	}
	return bd.dirtyX0, bd.dirtyY0, bd.dirtyX1, bd.dirtyY1
}

// ClearDirtyRegion resets dirty tracking
func (bd *BaseDevice) ClearDirtyRegion() {
	bd.hasDirty = false
	bd.dirtyX0 = -1
	bd.dirtyY0 = -1
	bd.dirtyX1 = -1
	bd.dirtyY1 = -1
}

// MarkDirty marks a rectangular region as dirty
func (bd *BaseDevice) MarkDirty(x0, y0, x1, y1 int) {
	// Clamp to valid bounds
	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 >= bd.config.Width {
		x1 = bd.config.Width - 1
	}
	if y1 >= bd.config.Height {
		y1 = bd.config.Height - 1
	}

	if !bd.hasDirty {
		bd.dirtyX0 = x0
		bd.dirtyY0 = y0
		bd.dirtyX1 = x1
		bd.dirtyY1 = y1
		bd.hasDirty = true
	} else {
		// Expand dirty region
		if x0 < bd.dirtyX0 {
			bd.dirtyX0 = x0
		}
		if y0 < bd.dirtyY0 {
			bd.dirtyY0 = y0
		}
		if x1 > bd.dirtyX1 {
			bd.dirtyX1 = x1
		}
		if y1 > bd.dirtyY1 {
			bd.dirtyY1 = y1
		}
	}
}

// Width returns display width
func (bd *BaseDevice) Width() int {
	return bd.config.Width
}

// Height returns display height
func (bd *BaseDevice) Height() int {
	return bd.config.Height
}

// ColorDepth returns bits per pixel
func (bd *BaseDevice) ColorDepth() int {
	return bd.config.ColorDepth
}

// PixelFormat returns the pixel format
func (bd *BaseDevice) PixelFormat() PixelFormat {
	return bd.config.PixelFormat
}
