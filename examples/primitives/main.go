package main

import (
	"log"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/flavioheleno/oled-emulator/emulator"
	"github.com/flavioheleno/oled-emulator/graphics"
)

// Example: Drawing Primitives - Demonstrates lines, rectangles, circles
func main() {
	// Create device
	dev := device.NewSSD1322(256, 64)

	// Create emulator
	emu := emulator.NewEmulator(dev, 2)
	emu.SetWindowTitle("OLED Emulator - Primitives")
	emu.ShowDebugInfo(false)

	// Create framebuffer
	fb := graphics.NewFrameBuffer(dev)

	// Clear to black
	fb.Clear(0x00)

	// Draw title
	font := graphics.DefaultBitmapFont()
	font.DrawString(fb, 10, 2, "Shapes", 0x0F)

	// Draw lines in different shades
	fb.DrawLine(10, 12, 50, 12, 0x0F)
	fb.DrawLine(10, 18, 80, 18, 0x0A)
	fb.DrawLine(10, 24, 110, 24, 0x05)

	// Draw rectangles
	fb.DrawRect(140, 5, 30, 20, 0x08, false)   // Outline
	fb.DrawRect(180, 5, 30, 20, 0x0C, true)    // Filled

	// Draw circles
	fb.DrawCircle(30, 45, 10, 0x0F, false)     // Outline
	fb.DrawCircle(70, 45, 10, 0x0A, true)      // Filled

	// Draw ellipses
	fb.DrawEllipse(120, 45, 15, 8, 0x08, false) // Outline
	fb.DrawEllipse(170, 45, 15, 8, 0x0C, true)  // Filled

	// Draw a triangle
	fb.DrawTriangle(220, 35, 240, 35, 230, 55, 0x0F, false)

	// Draw gradient effect
	for i := 0; i < 16; i++ {
		shade := byte(i)
		fb.SetPixel(10+i, 55, shade)
	}

	fb.Flush()

	// Run emulator
	if err := emu.Run(); err != nil {
		log.Fatalf("emulator error: %v", err)
	}
}
