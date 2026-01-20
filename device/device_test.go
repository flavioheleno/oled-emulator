package device

import (
	"testing"
)

func TestBaseDeviceCreation(t *testing.T) {
	config := Config{
		Width:       256,
		Height:      64,
		ColorDepth:  4,
		PixelFormat: HorizontalNibble,
		ColumnOffset: 28,
	}

	bd := NewBaseDevice(config)

	if bd.Width() != 256 {
		t.Errorf("expected width 256, got %d", bd.Width())
	}

	if bd.Height() != 64 {
		t.Errorf("expected height 64, got %d", bd.Height())
	}

	if bd.ColorDepth() != 4 {
		t.Errorf("expected color depth 4, got %d", bd.ColorDepth())
	}
}

func TestNibblePacking(t *testing.T) {
	mh := NewMemoryHelper(256, 64, HorizontalNibble, 28)
	vram := make([]byte, 480*64/2)

	// Test setting pixels
	tests := []struct {
		x, y int
		color byte
	}{
		{0, 0, 0x0F},
		{1, 0, 0x0A},
		{2, 0, 0x05},
	}

	for _, test := range tests {
		if err := mh.SetPixelNibble(vram, test.x, test.y, test.color); err != nil {
			t.Fatalf("failed to set pixel: %v", err)
		}

		// Read back
		pixel, err := mh.GetPixelNibble(vram, test.x, test.y)
		if err != nil {
			t.Fatalf("failed to get pixel: %v", err)
		}

		if pixel != test.color {
			t.Errorf("expected pixel %X, got %X", test.color, pixel)
		}
	}
}

func TestDirtyTracking(t *testing.T) {
	config := Config{
		Width:       256,
		Height:      64,
		ColorDepth:  4,
		PixelFormat: HorizontalNibble,
		ColumnOffset: 28,
	}

	bd := NewBaseDevice(config)

	// Initially no dirty region
	x0, y0, x1, y1 := bd.GetDirtyRegion()
	if x0 != -1 {
		t.Errorf("expected no dirty region initially, got (%d, %d, %d, %d)", x0, y0, x1, y1)
	}

	// Mark a region as dirty
	bd.MarkDirty(10, 20, 30, 40)

	x0, y0, x1, y1 = bd.GetDirtyRegion()
	if x0 != 10 || y0 != 20 || x1 != 30 || y1 != 40 {
		t.Errorf("expected dirty region (10, 20, 30, 40), got (%d, %d, %d, %d)", x0, y0, x1, y1)
	}

	// Expand dirty region
	bd.MarkDirty(5, 15, 50, 50)

	x0, y0, x1, y1 = bd.GetDirtyRegion()
	if x0 != 5 || y0 != 15 || x1 != 50 || y1 != 50 {
		t.Errorf("expected expanded dirty region (5, 15, 50, 50), got (%d, %d, %d, %d)", x0, y0, x1, y1)
	}

	// Clear dirty region
	bd.ClearDirtyRegion()

	x0, y0, x1, y1 = bd.GetDirtyRegion()
	if x0 != -1 {
		t.Errorf("expected no dirty region after clear, got (%d, %d, %d, %d)", x0, y0, x1, y1)
	}
}

func TestSSD1322Creation(t *testing.T) {
	ssd := NewSSD1322(256, 64)

	if ssd.Width() != 256 {
		t.Errorf("expected width 256, got %d", ssd.Width())
	}

	if ssd.Height() != 64 {
		t.Errorf("expected height 64, got %d", ssd.Height())
	}

	if ssd.IsDisplayOn() {
		t.Error("display should be off initially")
	}
}

func TestSSD1322Commands(t *testing.T) {
	ssd := NewSSD1322(256, 64)

	// Unlock commands
	ssd.ProcessCommand(CmdCommandLock, []byte{0xB1})

	if ssd.commandLocked {
		t.Error("commands should be unlocked")
	}

	// Turn display on
	ssd.ProcessCommand(CmdNormalDisplay, nil)

	if !ssd.IsDisplayOn() {
		t.Error("display should be on")
	}

	// Set contrast
	ssd.ProcessCommand(CmdSetContrast, []byte{0x80})

	if ssd.GetContrastLevel() != 0x80 {
		t.Errorf("expected contrast 0x80, got 0x%02X", ssd.GetContrastLevel())
	}

	// Invert display
	ssd.ProcessCommand(CmdInvertDisplay, []byte{0x01})

	if !ssd.IsInverted() {
		t.Error("display should be inverted")
	}
}

func TestSSD1322SetPixel(t *testing.T) {
	ssd := NewSSD1322(256, 64)

	// Set some pixels
	tests := []struct {
		x, y int
		color byte
	}{
		{0, 0, 0x0F},
		{100, 32, 0x08},
		{255, 63, 0x01},
	}

	for _, test := range tests {
		if err := ssd.SetPixel(test.x, test.y, test.color); err != nil {
			t.Fatalf("failed to set pixel: %v", err)
		}

		// Read back
		pixel, err := ssd.GetPixel(test.x, test.y)
		if err != nil {
			t.Fatalf("failed to get pixel: %v", err)
		}

		// Only lower 4 bits should be set
		if pixel != (test.color&0x0F) {
			t.Errorf("expected pixel 0x%02X, got 0x%02X", test.color&0x0F, pixel)
		}
	}
}

func TestSSD1322Reset(t *testing.T) {
	ssd := NewSSD1322(256, 64)

	// Modify some state
	ssd.ProcessCommand(CmdCommandLock, []byte{0xB1})
	ssd.ProcessCommand(CmdNormalDisplay, nil)
	ssd.ProcessCommand(CmdSetContrast, []byte{0x50})

	// Verify state changed
	if ssd.commandLocked {
		t.Error("commands should be unlocked")
	}
	if !ssd.IsDisplayOn() {
		t.Error("display should be on")
	}

	// Reset
	if err := ssd.Reset(); err != nil {
		t.Fatalf("reset failed: %v", err)
	}

	// Verify state restored
	if !ssd.commandLocked {
		t.Error("commands should be locked after reset")
	}
	if ssd.IsDisplayOn() {
		t.Error("display should be off after reset")
	}
	if ssd.GetContrastLevel() != 0x7F {
		t.Errorf("contrast should be 0x7F after reset, got 0x%02X", ssd.GetContrastLevel())
	}
}
