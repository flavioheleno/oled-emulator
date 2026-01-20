# AGENTS.md - Development Guide for Coding Agents

This guide explains how to effectively work on the OLED Display Emulator project as a coding agent.

## Quick Start for Agents

### 1. Understanding the Project Structure

```
oled-emulator/
├── device/          # Hardware emulation layer (SSD1322 controller, VRAM management)
├── graphics/        # Drawing API (primitives, fonts, text rendering)
├── animation/       # Animation system (easing, tweens, animator)
├── emulator/        # Desktop window (ebiten integration, rendering)
├── protocol/        # SPI communication bridge and command definitions
├── cmd/emulator/    # Standalone emulator binary entry point
├── examples/        # Usage examples
├── docs/            # API reference and protocol documentation
├── go.mod / go.sum  # Dependency management
├── README.md        # Project overview
├── PLAN.md          # Implementation roadmap
└── tests            # Unit tests throughout packages
```

### 2. How to Get Started

```bash
# Clone and enter directory
cd /home/flavio/Work/fhlabs/dispatch/emulator

# Download dependencies (IMPORTANT: See dependency notes below)
go mod tidy

# Run tests to verify everything works
go test ./...

# Build the emulator
go build -o emulator ./cmd/emulator/

# Try an example
go run examples/hello_world/main.go
```

## Architecture Overview

### Layer Model (Bottom to Top)

```
┌─────────────────────────────────────────────────────┐
│  Emulator (Window Management, Rendering)             │
│  - ebiten game loop integration                      │
│  - VRAM to screen conversion                         │
│  - Palette management                                │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────┐
│  Graphics API (High-level Drawing)                   │
│  - FrameBuffer (drawing primitives)                  │
│  - Text rendering (fonts, alignment)                │
│  - Image support                                     │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────┐
│  Animation System (Frame-based)                      │
│  - Animator (goroutine-based controller)            │
│  - Easing functions (20+ variants)                  │
│  - Tweens (value interpolation)                     │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────┐
│  Protocol Bridge (SPI Communication)                 │
│  - SPIBridge (DC/CS pin control)                    │
│  - CommandBuilder (fluent command construction)     │
│  - Command registry and initialization              │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────────┐
│  Device Layer (Hardware Emulation)                   │
│  - SSD1322 command processor                        │
│  - VRAM management (480x64 internal)                │
│  - Pixel format abstraction (Nibble/Vertical/RGB)   │
│  - Dirty region tracking                            │
└─────────────────────────────────────────────────────┘
```

### Design Principles

1. **Layered Architecture**: Each layer is independent and can be tested separately
2. **Interface-based**: Heavy use of interfaces for flexibility and testability
3. **Efficient Updates**: Dirty region tracking avoids full screen redraws
4. **Authentic Behavior**: SSD1322 command processor mimics real hardware
5. **Extensible**: Easy to add new pixel formats, display controllers, fonts

## Working on Specific Layers

### Device Layer (`device/`)

**Responsibilities:**
- Emulate SSD1322 controller behavior
- Manage VRAM (480×64 internal, 256×64 display)
- Handle pixel format conversions
- Track dirty regions for optimization

**Key Files:**
- `device.go` - Base device interface and configuration
- `ssd1322.go` - SSD1322 command processor (~280 lines)
- `memory.go` - VRAM utilities for nibble/vertical/RGB formats

**When to modify:**
- Adding support for different display controllers (SSD1306, SSD1351, etc.)
- Implementing new commands in SSD1322
- Optimizing VRAM access patterns

**Testing:**
```bash
go test ./device/... -v
```

### Graphics Layer (`graphics/`)

**Responsibilities:**
- Provide high-level drawing API
- Implement geometric primitives (Bresenham algorithms)
- Handle text rendering with multiple font types
- Support image drawing and conversion

**Key Files:**
- `framebuffer.go` - Main drawing API (~150 lines)
- `primitives.go` - Line/circle/ellipse/triangle algorithms (~250 lines)
- `text.go` - Text rendering framework
- `bitmap.go` - Bitmap font implementation
- `image.go` - Image support utilities

**When to modify:**
- Adding new drawing primitives
- Implementing new text rendering features
- Adding image effects or filters

**Testing:**
```bash
go test ./graphics/... -v
```

### Animation Layer (`animation/`)

**Responsibilities:**
- Manage frame-based animations
- Provide easing functions for smooth transitions
- Support tween composition (sequences and parallel)

**Key Files:**
- `animator.go` - Frame controller with goroutine support (~150 lines)
- `easing.go` - 20+ easing functions (~300 lines)
- `tween.go` - Value/color interpolation and composition (~280 lines)

**When to modify:**
- Adding new easing functions
- Implementing more complex animation patterns
- Optimizing animation updates

**Testing:**
```bash
go test ./animation/... -v
```

### Emulator Layer (`emulator/`)

**Responsibilities:**
- Manage ebiten window and game loop
- Convert VRAM to displayable images
- Handle palette and color mapping

**Key Files:**
- `window.go` - ebiten Game interface (~130 lines)
- `rendering.go` - VRAM to image conversion (~90 lines)

**When to modify:**
- Adding UI features (buttons, keyboard input)
- Changing window behavior
- Adding different rendering modes

**Testing:**
Manual - run examples and verify visual output

### Protocol Layer (`protocol/`)

**Responsibilities:**
- Emulate SPI communication
- Provide command definitions and metadata
- Generate initialization sequences

**Key Files:**
- `spi.go` - SPI bridge with DC/CS pin control (~100 lines)
- `commands.go` - Command registry and builders (~350 lines)

**When to modify:**
- Supporting different communication interfaces (I2C, 8080 parallel)
- Adding new command utilities
- Creating initialization sequences for different configurations

**Testing:**
```bash
go test ./protocol/... -v
```

## Common Development Tasks

### Adding a New Drawing Primitive

1. **Implement algorithm** in `graphics/primitives.go`
   ```go
   func DrawNewShape(fb *FrameBuffer, /* params */, color byte, setPixel func(int, int, byte)) {
       // Algorithm implementation
   }
   ```

2. **Add public API** in `graphics/framebuffer.go`
   ```go
   func (fb *FrameBuffer) DrawNewShape(/* params */, color byte) error {
       // Validation and pixel clamping
       DrawNewShape(fb, /* params */, color, func(x, y int, c byte) {
           if x >= 0 && x < fb.Width() && y >= 0 && y < fb.Height() {
               fb.device.SetPixel(x, y, c)
               fb.dirty = true
           }
       })
       return nil
   }
   ```

3. **Add tests** in `graphics/framebuffer_test.go`

4. **Update examples** if relevant

### Adding a New Font Type

1. **Implement Font interface** in new file `graphics/newfont.go`
   ```go
   type NewFont struct { /* fields */ }

   func (nf *NewFont) DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error) { ... }
   func (nf *NewFont) MeasureString(text string) (width, height int, err error) { ... }
   func (nf *NewFont) Height() int { ... }
   func (nf *NewFont) GetGlyph(ch rune) (GlyphData, error) { ... }
   ```

2. **Add example** in `examples/text_demo/main.go`

3. **Document** in API.md

### Adding Support for a New Display Controller

1. **Create new file** `device/ssd1351.go` (or similar)

2. **Implement Device interface**
   ```go
   type SSD1351 struct {
       *BaseDevice
       // controller-specific fields
   }

   func NewSSD1351(width, height int) *SSD1351 { ... }
   func (s *SSD1351) ProcessCommand(cmd byte, data []byte) error { ... }
   // other interface methods
   ```

3. **Add tests** `device/ssd1351_test.go`

4. **Create example** using new controller

5. **Document** in PROTOCOL.md

### Adding New Easing Function

1. **Add to** `animation/easing.go`
   ```go
   func EaseNewType(t float64) float64 {
       t = clamp(t)
       // Implementation
       return result
   }
   ```

2. **Add test** in `animation/animation_test.go`

3. **Update** docs/API.md

## Testing Guidelines

### Unit Tests

Every package should have comprehensive tests:

```go
// File: package/package_test.go
package package

import "testing"

func TestFeature(t *testing.T) {
    // Arrange
    obj := NewObject()

    // Act
    result := obj.DoSomething()

    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### Running Tests

```bash
# All tests
go test ./...

# Single package
go test ./device/...

# Verbose output
go test ./... -v

# With coverage
go test ./... -cover
```

### Current Test Status

- ✅ device: 7 tests passing
- ✅ graphics: 9 tests passing
- ✅ animation: 10 tests passing
- ✅ protocol: 10 tests passing
- **Total: 36 tests, all passing**

## Dependency Management

### Current Dependencies

```
✓ Only github.com/hajimehoshi/ebiten/v2 v2.9.7 (NOT v1)
✓ Only golang.org/x/image v0.35.0
✓ Standard library for everything else
```

### Important: Do NOT Add v1 Imports

**Wrong:**
```go
import "github.com/hajimehoshi/ebiten"  // This is v1!
```

**Correct:**
```go
import "github.com/hajimehoshi/ebiten/v2"
```

### Adding Dependencies

Before adding external dependencies:

1. **Ask yourself**: Can this be done with standard library?
2. **Check existing**: Are we already using something similar?
3. **Evaluate**: Performance, maintenance, license compatibility
4. **Discuss**: This should be coordinated with team

Safe dependencies to consider:
- `golang.org/x/*` - Official Go extensions
- `github.com/go-*` - Well-maintained Go packages
- Other game dev libraries compatible with ebiten

### Updating Dependencies

```bash
# Check for updates
go list -u -m all

# Update everything carefully
go get -u ./...

# Verify integrity
go mod verify

# Run tests after updates
go test ./...
```

## Code Organization

### File Naming Conventions

```
feature.go           # Main implementation
feature_test.go      # Unit tests
feature_internal.go  # Internal helpers (if needed)
```

### Function Organization

```go
// 1. Types and constants
type MyType struct { }
const MyConstant = 42

// 2. Constructor
func NewMyType() *MyType { }

// 3. Public methods (alphabetical)
func (m *MyType) DoSomething() { }
func (m *MyType) GetValue() int { }

// 4. Private helpers (alphabetical)
func (m *MyType) helper() { }
```

### Commenting Guidelines

```go
// Package device provides hardware emulation for OLED displays.
package device

// Device defines the interface for display emulation.
type Device interface {
    // ProcessCommand processes a device command.
    ProcessCommand(cmd byte, data []byte) error
}

// NewDevice creates a new device with the given configuration.
func NewDevice(config Config) Device {
    // Implementation
}
```

**Rules:**
- Public functions/types: Always comment starting with name
- Private functions: Comment if behavior is non-obvious
- Complex logic: Add inline comments explaining why, not what
- No comments for obvious code: `x := 5  // Set x to 5` ❌

## Git Workflow

### Commit Messages

**Format:**
```
[Package] Brief description

Optional longer explanation of changes and rationale.

- Bullet points for multiple changes
- Reference related packages or issues
```

**Examples:**
```
[device] Add SSD1351 controller support

- Implement SSD1351 command processor
- Add 256-color support via RGB888 format
- Include initialization sequence
```

```
[graphics] Implement text rotation and scaling

Support rotated and scaled text rendering through
matrix transformation in the drawing pipeline.
```

### Branch Naming

```
feature/implement-feature-name
bugfix/fix-bug-description
docs/update-documentation
refactor/reorganize-code
```

### Before Committing

```bash
# 1. Format code
go fmt ./...

# 2. Run tests
go test ./...

# 3. Check for unused imports
go vet ./...

# 4. Build
go build ./...
```

## Documentation

### When to Update Docs

- **API changes**: Update `docs/API.md`
- **Protocol changes**: Update `docs/PROTOCOL.md`
- **New features**: Add examples and README section
- **Architecture changes**: Update this file (AGENTS.md)

### How to Document Code

1. **Function documentation**:
   ```go
   // DrawCircle draws a circle at (x, y) with radius r.
   // color is a 4-bit grayscale value (0-15).
   // If filled is true, the circle is filled; otherwise only outline.
   func (fb *FrameBuffer) DrawCircle(x, y, r int, color byte, filled bool) error
   ```

2. **Complex algorithms**: Include reference or explanation
   ```go
   // Bresenham's line algorithm for efficient line drawing.
   // See: https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
   func DrawLineBresenham(...)
   ```

3. **Design decisions**: Comment non-obvious choices
   ```go
   // We use dirty region tracking instead of full screen updates
   // for efficiency. This avoids re-rendering unchanged areas.
   type BaseDevice struct {
       dirtyX0, dirtyY0, dirtyX1, dirtyY1 int
   }
   ```

## Performance Considerations

### Optimization Strategies

1. **Dirty Region Tracking**: Only update changed pixels
   ```go
   x0, y0, x1, y1 := device.GetDirtyRegion()
   if x0 >= 0 {  // Has dirty region
       // Render only dirty area
   }
   ```

2. **VRAM Caching**: Minimize memory allocations
   ```go
   // Good: Reuse buffer
   buffer := make([]byte, width*height)

   // Bad: Allocate in loop
   for i := 0; i < count; i++ {
       b := make([]byte, width*height)  // ❌ Avoid
   }
   ```

3. **Bit Packing**: Use efficient pixel formats
   - HorizontalNibble: 2 pixels/byte (most efficient)
   - VerticalByte: 8 pixels/byte
   - RGB888: 3 bytes/pixel (for color displays)

### Profiling

```bash
# Generate CPU profile
go test ./device/... -cpuprofile=cpu.prof

# View profile
go tool pprof cpu.prof
```

## Extending the Project

### Adding a New Example

1. **Create directory**: `examples/my_feature/main.go`
2. **Implement**: Demonstrate feature usage
3. **Document**: Add comments explaining what it shows
4. **Commit**: With clear message

### Adding New Color Format

1. **Update** `device.PixelFormat` enum
2. **Implement** pixel format conversions in `memory.go`
3. **Add tests** for conversions
4. **Update** documentation

### Adding Command-line Interface

1. **Create** `cmd/cli/main.go`
2. **Add** flag parsing
3. **Integrate** with existing packages
4. **Document** usage in README

## Troubleshooting

### Common Issues

**Issue**: `package ebiten not found`
- **Solution**: Use `github.com/hajimehoshi/ebiten/v2` (with `/v2`)

**Issue**: Tests failing after dependency update
- **Solution**: Run `go mod tidy` and `go mod verify`

**Issue**: Dirty region not updating correctly
- **Solution**: Ensure `MarkDirty()` is called after pixel modifications

**Issue**: Display showing garbage
- **Solution**: Check VRAM addressing - verify column/row ranges match display size

### Debug Techniques

```go
// Enable debug output
fmt.Printf("Dirty region: (%d, %d) to (%d, %d)\n", x0, y0, x1, y1)

// Verify pixel values
pixel, _ := device.GetPixel(100, 32)
fmt.Printf("Pixel at (100, 32): 0x%02X\n", pixel)

// Check VRAM integrity
if len(device.GetFrameBuffer()) != expectedSize {
    fmt.Println("VRAM size mismatch!")
}
```

## Resources

### Key References

- **SSD1322 Datasheet**: `docs/PROTOCOL.md` has command reference
- **Go Module Docs**: `go help mod`
- **ebiten Documentation**: https://ebitenengine.org/
- **Project Documentation**: See `README.md`, `PLAN.md`, `docs/API.md`

### Related Documentation

- `README.md` - Project overview and features
- `PLAN.md` - Implementation phases and architecture
- `docs/API.md` - Complete API reference
- `docs/PROTOCOL.md` - SSD1322 command reference

## Quick Reference

### Essential Commands

```bash
# Setup
go mod tidy                    # Download dependencies
go test ./...                  # Run all tests

# Development
go fmt ./...                   # Format code
go vet ./...                   # Check for errors
go build ./cmd/emulator/       # Build binary
go run examples/hello_world/   # Run example

# Maintenance
go mod verify                  # Verify integrity
go list -u -m all             # Check for updates
git log --oneline             # View commit history
```

### Common Packages to Import

```go
// Device layer
"github.com/flavioheleno/oled-emulator/device"

// Graphics layer
"github.com/flavioheleno/oled-emulator/graphics"

// Animation layer
"github.com/flavioheleno/oled-emulator/animation"

// Emulator window
"github.com/flavioheleno/oled-emulator/emulator"

// Protocol layer
"github.com/flavioheleno/oled-emulator/protocol"

// Standard library (always available)
"fmt"
"image"
"time"
```

## Next Steps for Agents

1. **Read** README.md for project overview
2. **Read** PLAN.md for architecture details
3. **Run** tests to verify setup: `go test ./...`
4. **Build** emulator: `go build -o emulator ./cmd/emulator/`
5. **Try** examples: `go run examples/hello_world/main.go`
6. **Choose task** from issues or feature ideas
7. **Follow** this guide for implementation
8. **Test** thoroughly before committing
9. **Document** changes in relevant files
10. **Create** commit with clear message

Good luck! 🚀
