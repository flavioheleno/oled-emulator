# OLED Display Emulator

A configurable OLED display emulator written in Go, supporting the SSD1322 hardware command protocol. Perfect for UI design, prototyping, and testing display applications without real hardware.

## Features

- **SSD1322 Emulation**: Full support for SSD1322 controller commands and VRAM management
- **Graphics API**: High-level drawing functions (lines, rectangles, circles, polygons, etc.)
- **Text Rendering**: Bitmap and TrueType font support with alignment options
- **Animation System**: Frame-based animations with 20+ easing functions
- **SPI Protocol Bridge**: Emulate SPI communication for real driver compatibility
- **ebiten Integration**: Native desktop window with pixel scaling and OLED-style rendering
- **Multiple Pixel Formats**: Support for HorizontalNibble (SSD1322), VerticalByte (SSD1306), and RGB888
- **Dirty Region Tracking**: Efficient partial updates

## Installation

```bash
go get github.com/flavioheleno/oled-emulator
```

## Quick Start

### Hello World

```go
package main

import (
    "log"
    "github.com/flavioheleno/oled-emulator/device"
    "github.com/flavioheleno/oled-emulator/emulator"
    "github.com/flavioheleno/oled-emulator/graphics"
)

func main() {
    // Create a 256x64 SSD1322 display
    dev := device.NewSSD1322(256, 64)

    // Create emulator window
    emu := emulator.NewEmulator(dev, 2)
    emu.SetWindowTitle("My Display")

    // Create framebuffer and draw
    fb := graphics.NewFrameBuffer(dev)
    fb.DrawRect(10, 10, 100, 30, 0x0F, true)
    fb.Flush()

    // Show window
    if err := emu.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Architecture

### Core Packages

#### `device/`
- `device.go`: Device interface and configuration
- `memory.go`: VRAM management and pixel format conversions
- `ssd1322.go`: SSD1322 command processor and emulation

#### `graphics/`
- `framebuffer.go`: High-level drawing API
- `primitives.go`: Drawing algorithms (Bresenham, midpoint circle, etc.)
- `text.go`: Font interface and text rendering
- `bitmap.go`: Bitmap font implementation
- `truetype.go`: TrueType font support
- `image.go`: Image and sprite utilities

#### `animation/`
- `animator.go`: Frame-based animation controller
- `easing.go`: 20+ easing functions
- `tween.go`: Value interpolation and composition

#### `emulator/`
- `window.go`: ebiten window manager
- `rendering.go`: VRAM to screen conversion

#### `protocol/`
- `spi.go`: SPI communication bridge
- `commands.go`: Command definitions and builders

## Usage Examples

### Drawing Shapes

```go
fb := graphics.NewFrameBuffer(dev)

// Clear display
fb.Clear(0x00)

// Draw lines
fb.DrawLine(0, 0, 100, 100, 0x0F)

// Draw rectangles (outline and filled)
fb.DrawRect(50, 50, 100, 50, 0x0A, false)
fb.DrawRect(150, 50, 100, 50, 0x0F, true)

// Draw circles
fb.DrawCircle(128, 32, 20, 0x08, false)
fb.DrawCircle(200, 32, 20, 0x0C, true)

// Commit changes
fb.Flush()
```

### Text Rendering

```go
font := graphics.DefaultBitmapFont()
drawer := graphics.NewAlignedTextDrawer(font)

// Draw centered text
drawer.DrawCenteredText(fb, 128, 32, "Hello, World!", 0x0F)

// Draw left-aligned text
drawer.DrawAlignedText(fb, 10, 50, "Left", graphics.AlignLeft, 0x0F)

// Draw right-aligned text
drawer.DrawAlignedText(fb, 246, 50, "Right", graphics.AlignRight, 0x0F)
```

### Animations

```go
animator := animation.NewAnimator(60) // 60 FPS

// Simple value tween
tween := animation.NewTween(0, 100, 2*time.Second, animation.EaseInOutQuad)
tween.SetOnUpdate(func(value float64) {
    // Use value for animation
})

animator.AddAnimation(func(frame int, dt float64) bool {
    return tween.Update(dt)
})

animator.Start()
```

### SPI Protocol Integration

```go
bridge := protocol.NewSPIBridge(dev)

// Send initialization
initSeq := protocol.SSD1322InitSequence()
bridge.SendInitSequence(initSeq)

// Set contrast via command builder
cmd := protocol.ContrastCommand(0x80)
bridge.SetDC(false)
bridge.Write(cmd)

// Set display inversion
cmd = protocol.InversionCommand(true)
bridge.SetDC(false)
bridge.Write(cmd)
```

## Configuration

### Display Configuration

```go
config := device.Config{
    Width:       256,
    Height:      64,
    ColorDepth:  4,
    PixelFormat: device.HorizontalNibble,
    ColumnOffset: 28,
}

dev := device.NewBaseDevice(config)
```

### Emulator Configuration

```go
emu := emulator.NewEmulator(dev, 4) // 4x pixel scale

emu.SetWindowTitle("My Display")
emu.SetFrameRate(60)
emu.ShowDebugInfo(true)

// Custom palette
palette := emulator.NewGrayscalePalette()
emu.SetPalette(palette)
```

## Color Values

Colors are represented as 4-bit values (0x00 to 0x0F):

```
0x00 = Black
0x01 - 0x07 = Dark shades
0x08 - 0x0E = Medium shades
0x0F = White
```

## Examples

The repository includes several examples:

- `examples/hello_world/` - Basic text display
- `examples/primitives/` - Drawing shapes and lines
- `examples/animation/` - Animated transitions
- `examples/real_driver/` - SPI protocol integration

Run examples:

```bash
go run examples/hello_world/main.go
go run examples/primitives/main.go
go run examples/animation/main.go
go run examples/real_driver/main.go
```

## Testing

Run all tests:

```bash
go test ./...
```

Run specific package tests:

```bash
go test ./device/...
go test ./graphics/...
go test ./animation/...
go test ./protocol/...
```

## Performance

The emulator uses:
- **Dirty region tracking** for efficient updates
- **Goroutine-based animation** system
- **Memory-mapped VRAM** for direct pixel manipulation
- **Pixel format abstraction** for hardware compatibility

## Pixel Formats

### HorizontalNibble (SSD1322 Native)
- 2 pixels per byte
- 4-bit grayscale per pixel
- Used by: SSD1322

### VerticalByte (SSD1306 Style)
- 8 pixels per byte
- 1-bit monochrome per pixel
- Used by: SSD1306, SSD1306

### RGB888
- 24-bit color per pixel
- 8-bit per channel

## API Reference

See [API.md](docs/API.md) for detailed API documentation.

See [PROTOCOL.md](docs/PROTOCOL.md) for SSD1322 command reference.

## License

This project is provided for educational and testing purposes.

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## Resources

- [SSD1322 Datasheet](http://datasheets.chipfind.org/pdf/940.pdf)
- [ebiten Documentation](https://ebitenengine.org/)
- [Go Image Libraries](https://golang.org/pkg/image/)

## Acknowledgments

Built with:
- [ebiten](https://ebitenengine.org/) - Go graphics library
- Go standard library
- golang.org/x/image for font support
