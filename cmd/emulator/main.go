package main

import (
	"log"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/flavioheleno/oled-emulator/emulator"
	"github.com/flavioheleno/oled-emulator/graphics"
)

func main() {
	// Create SSD1322 device (256x64 4-bit grayscale)
	dev := device.NewSSD1322(256, 64)

	// Create emulator window with 2x pixel scale
	emu := emulator.NewEmulator(dev, 2)
	emu.SetWindowTitle("OLED Emulator - SSD1322 (256x64)")
	emu.ShowDebugInfo(true)
	emu.SetFrameRate(60)

	// Create a framebuffer for drawing
	fb := graphics.NewFrameBuffer(dev)

	// Draw a test pattern
	drawTestPattern(fb, dev)

	// Run the emulator
	if err := emu.Run(); err != nil {
		log.Fatalf("emulator error: %v", err)
	}
}

// drawTestPattern draws a test pattern on the display
func drawTestPattern(fb *graphics.FrameBuffer, dev device.Device) {
	// Clear display to black
	fb.Clear(0x00)

	// Draw some rectangles in different shades
	for i := 0; i < 4; i++ {
		shade := byte((i + 1) * 3)
		x := i * 64
		fb.DrawRect(x, 0, 64, 32, shade, true)
	}

	// Draw a circle
	fb.DrawCircle(128, 32, 15, 0x0F, false)
	fb.DrawCircle(128, 32, 12, 0x08, false)

	// Draw some lines
	fb.DrawLine(0, 32, 256, 32, 0x07)
	fb.DrawLine(128, 0, 128, 64, 0x07)

	// Draw a filled rectangle with gradient effect
	for i := 0; i < 32; i++ {
		shade := byte((i * 15) / 32)
		fb.DrawRect(200+i, 48, 1, 16, shade, true)
	}

	fb.Flush()
}
