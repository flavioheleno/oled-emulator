package graphics

import (
	"testing"

	"github.com/flavioheleno/oled-emulator/device"
)

func TestFrameBufferCreation(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if fb.Width() != 256 {
		t.Errorf("expected width 256, got %d", fb.Width())
	}

	if fb.Height() != 64 {
		t.Errorf("expected height 64, got %d", fb.Height())
	}
}

func TestFrameBufferClear(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.Clear(0x00); err != nil {
		t.Fatalf("clear failed: %v", err)
	}

	// All pixels should be 0
	for y := 0; y < fb.Height(); y++ {
		for x := 0; x < fb.Width(); x++ {
			pixel, err := fb.GetPixel(x, y)
			if err != nil {
				t.Fatalf("failed to get pixel: %v", err)
			}
			if pixel != 0 {
				t.Errorf("expected pixel 0, got %d", pixel)
			}
		}
	}
}

func TestFrameBufferSetPixel(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.SetPixel(10, 20, 0x0F); err != nil {
		t.Fatalf("set pixel failed: %v", err)
	}

	pixel, err := fb.GetPixel(10, 20)
	if err != nil {
		t.Fatalf("get pixel failed: %v", err)
	}

	if pixel != 0x0F {
		t.Errorf("expected pixel 0x0F, got 0x%02X", pixel)
	}

	if !fb.IsDirty() {
		t.Error("framebuffer should be dirty after set pixel")
	}
}

func TestFrameBufferDrawLine(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.DrawLine(0, 0, 10, 10, 0x0F); err != nil {
		t.Fatalf("draw line failed: %v", err)
	}

	// Check that at least some pixels are set
	pixelsSet := 0
	for y := 0; y < fb.Height(); y++ {
		for x := 0; x < fb.Width(); x++ {
			pixel, _ := fb.GetPixel(x, y)
			if pixel != 0 {
				pixelsSet++
			}
		}
	}

	if pixelsSet == 0 {
		t.Error("no pixels were set by draw line")
	}
}

func TestFrameBufferDrawRect(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.DrawRect(10, 10, 20, 20, 0x0F, false); err != nil {
		t.Fatalf("draw rect failed: %v", err)
	}

	// Check corners are set
	corners := []struct {
		x, y int
	}{
		{10, 10}, {29, 10}, {10, 29}, {29, 29},
	}

	for _, c := range corners {
		pixel, err := fb.GetPixel(c.x, c.y)
		if err != nil {
			t.Fatalf("get pixel failed: %v", err)
		}
		if pixel == 0 {
			t.Errorf("corner pixel at (%d, %d) should be set", c.x, c.y)
		}
	}
}

func TestFrameBufferDrawFilledRect(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.DrawRect(10, 10, 10, 10, 0x0F, true); err != nil {
		t.Fatalf("draw filled rect failed: %v", err)
	}

	// Check that interior pixels are set
	pixelsSet := 0
	for y := 10; y < 20; y++ {
		for x := 10; x < 20; x++ {
			pixel, _ := fb.GetPixel(x, y)
			if pixel != 0 {
				pixelsSet++
			}
		}
	}

	if pixelsSet != 100 {
		t.Errorf("expected 100 pixels set, got %d", pixelsSet)
	}
}

func TestFrameBufferDrawCircle(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := fb.DrawCircle(50, 32, 10, 0x0F, false); err != nil {
		t.Fatalf("draw circle failed: %v", err)
	}

	// Check that some pixels are set
	pixelsSet := 0
	for y := 0; y < fb.Height(); y++ {
		for x := 0; x < fb.Width(); x++ {
			pixel, _ := fb.GetPixel(x, y)
			if pixel != 0 {
				pixelsSet++
			}
		}
	}

	if pixelsSet == 0 {
		t.Error("no pixels were set by draw circle")
	}
}

func TestFrameBufferFlush(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	fb.SetPixel(10, 10, 0x0F)

	if !fb.IsDirty() {
		t.Error("framebuffer should be dirty before flush")
	}

	if err := fb.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	if fb.IsDirty() {
		t.Error("framebuffer should not be dirty after flush")
	}
}
