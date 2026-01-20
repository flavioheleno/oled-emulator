package main

import (
	"log"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/flavioheleno/oled-emulator/emulator"
	"github.com/flavioheleno/oled-emulator/graphics"
	"github.com/flavioheleno/oled-emulator/protocol"
)

// Example: Real Driver Integration - Shows how to use the SPI bridge
// This demonstrates how actual SSD1322 driver code would work with the emulator
func main() {
	// Create device
	dev := device.NewSSD1322(256, 64)

	// Create SPI bridge (this would normally interface with actual hardware)
	bridge := protocol.NewSPIBridge(dev)

	// Initialize display using the SPI bridge
	initSequence := protocol.SSD1322InitSequence()

	// Send initialization sequence
	log.Println("Sending initialization sequence...")
	if err := bridge.SendInitSequence(initSequence); err != nil {
		log.Fatalf("init error: %v", err)
	}

	// Set contrast using command builder
	log.Println("Setting contrast...")
	bridge.SetDC(false)
	contrastCmd := protocol.ContrastCommand(0x80)
	if err := bridge.Write(contrastCmd); err != nil {
		log.Fatalf("contrast command error: %v", err)
	}

	// Create emulator to visualize
	emu := emulator.NewEmulator(dev, 2)
	emu.SetWindowTitle("OLED Emulator - Real Driver Integration")
	emu.ShowDebugInfo(true)

	// Create framebuffer to draw content
	fb := graphics.NewFrameBuffer(dev)

	// Draw test pattern
	drawTestPattern(fb)

	fb.Flush()

	log.Println("Display initialized. Launching emulator...")

	// Run emulator
	if err := emu.Run(); err != nil {
		log.Fatalf("emulator error: %v", err)
	}
}

// drawTestPattern draws a test pattern on the display
func drawTestPattern(fb *graphics.FrameBuffer) {
	// Clear display
	fb.Clear(0x00)

	// Draw gradient bars
	for row := 0; row < 16; row++ {
		shade := byte(row)
		for col := 0; col < 16; col++ {
			fb.SetPixel(col*16, row, shade)
		}
	}

	// Draw checkerboard pattern
	for y := 32; y < 48; y++ {
		for x := 0; x < 128; x++ {
			if ((x + y) % 4) == 0 {
				fb.SetPixel(x, y, 0x0F)
			}
		}
	}

	// Draw vertical lines
	for x := 130; x < 256; x += 10 {
		for y := 0; y < 64; y++ {
			if (y % 2) == 0 {
				fb.SetPixel(x, y, 0x08)
			}
		}
	}

	// Draw text info
	font := graphics.DefaultBitmapFont()
	font.DrawString(fb, 10, 50, "SPI Bridge Test", 0x0F)
}

// This example shows how the emulator can replace actual hardware
// In a real scenario, this code could work with both:
// 1. The emulator (for testing/development)
// 2. Actual hardware (by replacing the device initialization)
//
// The SPI bridge provides a protocol-level interface that matches
// the actual SSD1322 communication protocol, making it easy to
// switch between hardware and emulation.
