package main

import (
	"log"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/flavioheleno/oled-emulator/emulator"
	"github.com/flavioheleno/oled-emulator/graphics"
)

// Example: Hello World - Basic text display
func main() {
	// Create SSD1322 device
	dev := device.NewSSD1322(256, 64)

	// Create emulator window
	emu := emulator.NewEmulator(dev, 2)
	emu.SetWindowTitle("OLED Emulator - Hello World")
	emu.ShowDebugInfo(false)

	// Create framebuffer for drawing
	fb := graphics.NewFrameBuffer(dev)

	// Create a bitmap font
	font := graphics.DefaultBitmapFont()

	// Create aligned text drawer
	drawer := graphics.NewAlignedTextDrawer(font)

	// Draw title
	drawer.DrawCenteredText(fb, 128, 10, "Hello, OLED!", 0x0F)

	// Draw some decorative lines
	fb.DrawLine(20, 22, 236, 22, 0x08)
	fb.DrawLine(20, 42, 236, 42, 0x08)

	// Draw subtitle
	drawer.DrawCenteredText(fb, 128, 28, "Emulator Example", 0x0A)

	// Draw some info text
	drawer.DrawCenteredText(fb, 128, 50, "Press ESC to exit", 0x07)

	// Flush changes to device
	fb.Flush()

	// Run emulator
	if err := emu.Run(); err != nil {
		log.Fatalf("emulator error: %v", err)
	}
}
