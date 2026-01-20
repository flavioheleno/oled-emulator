package graphics

import (
	"fmt"

	"github.com/flavioheleno/oled-emulator/device"
)

// FrameBuffer provides a high-level drawing API on top of a device
type FrameBuffer struct {
	device device.Device
	buffer []byte
	dirty  bool
}

// NewFrameBuffer creates a new framebuffer for a device
func NewFrameBuffer(dev device.Device) *FrameBuffer {
	fb := &FrameBuffer{
		device: dev,
		buffer: make([]byte, len(dev.GetFrameBuffer())),
		dirty:  false,
	}

	// Copy initial buffer
	copy(fb.buffer, dev.GetFrameBuffer())

	return fb
}

// Clear fills the entire framebuffer with a color
func (fb *FrameBuffer) Clear(color byte) error {
	width := fb.device.Width()
	height := fb.device.Height()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if err := fb.SetPixel(x, y, color); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetPixel sets a pixel at the given coordinates
func (fb *FrameBuffer) SetPixel(x, y int, color byte) error {
	if err := fb.device.SetPixel(x, y, color); err != nil {
		return err
	}

	fb.dirty = true
	return nil
}

// GetPixel reads a pixel at the given coordinates
func (fb *FrameBuffer) GetPixel(x, y int) (byte, error) {
	return fb.device.GetPixel(x, y)
}

// DrawLine draws a line from (x0, y0) to (x1, y1)
func (fb *FrameBuffer) DrawLine(x0, y0, x1, y1 int, color byte) error {
	color = color & 0x0F // Ensure 4-bit color for SSD1322

	DrawLineBresenham(fb, x0, y0, x1, y1, color, func(x, y int, c byte) {
		// Clamp coordinates
		if x >= 0 && x < fb.device.Width() && y >= 0 && y < fb.device.Height() {
			fb.device.SetPixel(x, y, c)
			fb.dirty = true
		}
	})

	return nil
}

// DrawRect draws a rectangle outline or filled rectangle
func (fb *FrameBuffer) DrawRect(x, y, w, h int, color byte, filled bool) error {
	if w < 0 || h < 0 {
		return fmt.Errorf("invalid rectangle dimensions: %dx%d", w, h)
	}

	color = color & 0x0F

	DrawRect(fb, x, y, w, h, color, filled, func(px, py int, c byte) {
		if px >= 0 && px < fb.device.Width() && py >= 0 && py < fb.device.Height() {
			fb.device.SetPixel(px, py, c)
			fb.dirty = true
		}
	})

	return nil
}

// DrawCircle draws a circle outline or filled circle
func (fb *FrameBuffer) DrawCircle(x, y, r int, color byte, filled bool) error {
	if r < 0 {
		return fmt.Errorf("invalid circle radius: %d", r)
	}

	color = color & 0x0F

	DrawCircle(fb, x, y, r, color, filled, func(px, py int, c byte) {
		if px >= 0 && px < fb.device.Width() && py >= 0 && py < fb.device.Height() {
			fb.device.SetPixel(px, py, c)
			fb.dirty = true
		}
	})

	return nil
}

// DrawEllipse draws an ellipse outline or filled ellipse
func (fb *FrameBuffer) DrawEllipse(x, y, rx, ry int, color byte, filled bool) error {
	if rx < 0 || ry < 0 {
		return fmt.Errorf("invalid ellipse radii: %dx%d", rx, ry)
	}

	color = color & 0x0F

	DrawEllipse(fb, x, y, rx, ry, color, filled, func(px, py int, c byte) {
		if px >= 0 && px < fb.device.Width() && py >= 0 && py < fb.device.Height() {
			fb.device.SetPixel(px, py, c)
			fb.dirty = true
		}
	})

	return nil
}

// DrawTriangle draws a triangle outline or filled triangle
func (fb *FrameBuffer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, color byte, filled bool) error {
	color = color & 0x0F

	DrawTriangle(fb, x1, y1, x2, y2, x3, y3, color, filled, func(px, py int, c byte) {
		if px >= 0 && px < fb.device.Width() && py >= 0 && py < fb.device.Height() {
			fb.device.SetPixel(px, py, c)
			fb.dirty = true
		}
	})

	return nil
}

// FillRegion fills a rectangular region with a solid color
func (fb *FrameBuffer) FillRegion(x, y, w, h int, color byte) error {
	if w < 0 || h < 0 {
		return fmt.Errorf("invalid fill region dimensions: %dx%d", w, h)
	}

	color = color & 0x0F

	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if px >= 0 && px < fb.device.Width() && py >= 0 && py < fb.device.Height() {
				fb.device.SetPixel(px, py, color)
				fb.dirty = true
			}
		}
	}

	return nil
}

// Flush commits any changes to the device's VRAM
func (fb *FrameBuffer) Flush() error {
	if !fb.dirty {
		return nil
	}

	// Update internal buffer from device
	copy(fb.buffer, fb.device.GetFrameBuffer())
	fb.dirty = false

	return nil
}

// IsDirty returns whether the framebuffer has been modified since last flush
func (fb *FrameBuffer) IsDirty() bool {
	return fb.dirty
}

// GetBuffer returns a copy of the current framebuffer
func (fb *FrameBuffer) GetBuffer() []byte {
	result := make([]byte, len(fb.buffer))
	copy(result, fb.buffer)
	return result
}

// GetDevice returns the underlying device
func (fb *FrameBuffer) GetDevice() device.Device {
	return fb.device
}

// Width returns the framebuffer width
func (fb *FrameBuffer) Width() int {
	return fb.device.Width()
}

// Height returns the framebuffer height
func (fb *FrameBuffer) Height() int {
	return fb.device.Height()
}
