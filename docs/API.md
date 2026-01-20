# API Reference

## Device Package

### Device Interface

```go
type Device interface {
    ProcessCommand(cmd byte, data []byte) error
    GetFrameBuffer() []byte
    GetDirtyRegion() (x0, y0, x1, y1 int)
    ClearDirtyRegion()
    Width() int
    Height() int
    ColorDepth() int
    PixelFormat() PixelFormat
    Reset() error
    SetPixel(x, y int, color byte) error
    GetPixel(x, y int) (byte, error)
}
```

### SSD1322

```go
func NewSSD1322(width, height int) *SSD1322
func (ssd *SSD1322) WriteData(data []byte) error
func (ssd *SSD1322) IsDisplayOn() bool
func (ssd *SSD1322) GetContrastLevel() byte
func (ssd *SSD1322) IsInverted() bool
```

### Memory Helper

```go
type MemoryHelper struct {}

func NewMemoryHelper(width, height int, pixelFormat PixelFormat, colOffset int) *MemoryHelper
func (mh *MemoryHelper) SetPixelNibble(vram []byte, x, y int, color byte) error
func (mh *MemoryHelper) GetPixelNibble(vram []byte, x, y int) (byte, error)
func (mh *MemoryHelper) SetPixelVertical(vram []byte, x, y int, color byte) error
func (mh *MemoryHelper) GetPixelVertical(vram []byte, x, y int) (byte, error)
func (mh *MemoryHelper) FillRegionNibble(vram []byte, x0, y0, x1, y1 int, color byte) error
func (mh *MemoryHelper) FillRegionVertical(vram []byte, x0, y0, x1, y1 int, color byte) error
```

## Graphics Package

### FrameBuffer

```go
type FrameBuffer struct {}

func NewFrameBuffer(dev device.Device) *FrameBuffer
func (fb *FrameBuffer) Clear(color byte) error
func (fb *FrameBuffer) SetPixel(x, y int, color byte) error
func (fb *FrameBuffer) GetPixel(x, y int) (byte, error)
func (fb *FrameBuffer) DrawLine(x0, y0, x1, y1 int, color byte) error
func (fb *FrameBuffer) DrawRect(x, y, w, h int, color byte, filled bool) error
func (fb *FrameBuffer) DrawCircle(x, y, r int, color byte, filled bool) error
func (fb *FrameBuffer) DrawEllipse(x, y, rx, ry int, color byte, filled bool) error
func (fb *FrameBuffer) DrawTriangle(x1, y1, x2, y2, x3, y3 int, color byte, filled bool) error
func (fb *FrameBuffer) FillRegion(x, y, w, h int, color byte) error
func (fb *FrameBuffer) Flush() error
func (fb *FrameBuffer) IsDirty() bool
func (fb *FrameBuffer) Width() int
func (fb *FrameBuffer) Height() int
```

### Drawing Primitives

```go
// Lines
func DrawLineBresenham(fb *FrameBuffer, x0, y0, x1, y1 int, color byte, setPixel func(int, int, byte))

// Circles
func DrawCircle(fb *FrameBuffer, cx, cy, r int, color byte, filled bool, setPixel func(int, int, byte))
func DrawCircleOutline(fb *FrameBuffer, cx, cy, r int, color byte, setPixel func(int, int, byte))
func DrawFilledCircle(fb *FrameBuffer, cx, cy, r int, color byte, setPixel func(int, int, byte))

// Ellipses
func DrawEllipse(fb *FrameBuffer, cx, cy, rx, ry int, color byte, filled bool, setPixel func(int, int, byte))
func DrawEllipseOutline(fb *FrameBuffer, cx, cy, rx, ry int, color byte, setPixel func(int, int, byte))
func DrawFilledEllipse(fb *FrameBuffer, cx, cy, rx, ry int, color byte, setPixel func(int, int, byte))

// Rectangles and Triangles
func DrawRect(fb *FrameBuffer, x, y, w, h int, color byte, filled bool, setPixel func(int, int, byte))
func DrawTriangle(fb *FrameBuffer, x1, y1, x2, y2, x3, y3 int, color byte, filled bool, setPixel func(int, int, byte))
func DrawFilledTriangle(fb *FrameBuffer, x1, y1, x2, y2, x3, y3 int, color byte, setPixel func(int, int, byte))

// Utility functions
func Clamp(value, minVal, maxVal int) int
func Lerp(a, b float64, t float64) float64
func Map(value, inMin, inMax, outMin, outMax float64) float64
func Distance(x1, y1, x2, y2 float64) float64
```

### Font Interface

```go
type Font interface {
    DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error)
    MeasureString(text string) (width, height int, err error)
    Height() int
    GetGlyph(ch rune) (GlyphData, error)
}

type GlyphData struct {
    Width    int
    Height   int
    AdvanceX int
    BearingX int
    BearingY int
    Data     []byte
}
```

### Text Rendering

```go
type TextRenderer struct {}
func NewTextRenderer(font Font) *TextRenderer
func (tr *TextRenderer) SetOptions(opts TextOptions)
func (tr *TextRenderer) DrawText(fb *FrameBuffer, x, y int, text string) (int, error)
func (tr *TextRenderer) DrawMultilineText(fb *FrameBuffer, x, y int, text string) error
func (tr *TextRenderer) MeasureMultilineText(text string) (width, height int, err error)

type AlignedTextDrawer struct {}
func NewAlignedTextDrawer(font Font) *AlignedTextDrawer
func (atd *AlignedTextDrawer) DrawAlignedText(fb *FrameBuffer, x, y int, text string, alignment TextAlignment, color byte) error
func (atd *AlignedTextDrawer) DrawCenteredText(fb *FrameBuffer, x, y int, text string, color byte) error
func (atd *AlignedTextDrawer) DrawRightAlignedText(fb *FrameBuffer, x, y int, text string, color byte) error
```

### Bitmap Font

```go
type BitmapFont struct {}
func NewBitmapFont(width, height, advance int) *BitmapFont
func DefaultBitmapFont() *BitmapFont
func (bf *BitmapFont) AddGlyph(ch rune, data GlyphData)
func (bf *BitmapFont) DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error)
func (bf *BitmapFont) MeasureString(text string) (width, height int, err error)
func (bf *BitmapFont) GetGlyph(ch rune) (GlyphData, error)
```

### Image Support

```go
func DrawImage(fb *FrameBuffer, x, y int, img image.Image) error
func DrawImageScaled(fb *FrameBuffer, x, y, w, h int, img image.Image) error

type ImageTiler struct {}
func NewImageTiler(img image.Image) *ImageTiler
func (it *ImageTiler) DrawTiled(fb *FrameBuffer, x, y, w, h int) error

func ConvertToGrayscale(src image.Image) image.Image
func ConvertToBitmap(src image.Image, threshold uint8) image.Image
```

## Animation Package

### Animator

```go
type Animator struct {}
type AnimationFunc func(frame int, dt float64) bool

func NewAnimator(fps int) *Animator
func (a *Animator) SetFrameRate(fps int)
func (a *Animator) AddAnimation(fn AnimationFunc)
func (a *Animator) SetOnFrame(fn func(frame int, dt float64))
func (a *Animator) Start()
func (a *Animator) Stop()
func (a *Animator) IsRunning() bool
func (a *Animator) GetFrameCount() int
func (a *Animator) GetAnimationCount() int
func (a *Animator) Clear()
func (a *Animator) WaitForCompletion(timeout time.Duration) bool
```

### Easing Functions

```go
type EasingFunc func(t float64) float64

// Basic easing
func Linear(t float64) float64
func EaseInQuad(t float64) float64
func EaseOutQuad(t float64) float64
func EaseInOutQuad(t float64) float64

// Cubic easing
func EaseInCubic(t float64) float64
func EaseOutCubic(t float64) float64
func EaseInOutCubic(t float64) float64

// Quartic easing
func EaseInQuart(t float64) float64
func EaseOutQuart(t float64) float64
func EaseInOutQuart(t float64) float64

// Quintic easing
func EaseInQuint(t float64) float64
func EaseOutQuint(t float64) float64
func EaseInOutQuint(t float64) float64

// Sinusoidal easing
func EaseInSine(t float64) float64
func EaseOutSine(t float64) float64
func EaseInOutSine(t float64) float64

// Exponential easing
func EaseInExpo(t float64) float64
func EaseOutExpo(t float64) float64
func EaseInOutExpo(t float64) float64

// Circular easing
func EaseInCirc(t float64) float64
func EaseOutCirc(t float64) float64
func EaseInOutCirc(t float64) float64

// Back easing
func EaseInBack(t float64) float64
func EaseOutBack(t float64) float64
func EaseInOutBack(t float64) float64

// Elastic easing
func EaseInElastic(t float64) float64
func EaseOutElastic(t float64) float64
func EaseInOutElastic(t float64) float64

// Bounce easing
func EaseInBounce(t float64) float64
func EaseOutBounce(t float64) float64
func EaseInOutBounce(t float64) float64
```

### Tween

```go
type Tween struct {}
func NewTween(from, to float64, duration time.Duration, easing EasingFunc) *Tween
func (t *Tween) SetOnComplete(fn func()) *Tween
func (t *Tween) SetOnUpdate(fn func(value float64)) *Tween
func (t *Tween) GetValue() float64
func (t *Tween) IsComplete() bool
func (t *Tween) GetProgress() float64
func (t *Tween) Update(dt float64) bool

type ColorTween struct {}
func NewColorTween(fromR, fromG, fromB, toR, toG, toB byte, duration time.Duration, easing EasingFunc) *ColorTween
func (ct *ColorTween) SetOnComplete(fn func()) *ColorTween
func (ct *ColorTween) SetOnUpdate(fn func(r, g, b byte)) *ColorTween
func (ct *ColorTween) GetColor() (byte, byte, byte)
func (ct *ColorTween) IsComplete() bool
func (ct *ColorTween) Update(dt float64) bool

type SequenceTween struct {}
func NewSequenceTween(tweens ...*Tween) *SequenceTween
func (st *SequenceTween) SetOnComplete(fn func()) *SequenceTween
func (st *SequenceTween) Update(dt float64) bool
func (st *SequenceTween) IsComplete() bool

type ParallelTween struct {}
func NewParallelTween(tweens ...*Tween) *ParallelTween
func (pt *ParallelTween) SetOnComplete(fn func()) *ParallelTween
func (pt *ParallelTween) Update(dt float64) bool
func (pt *ParallelTween) IsComplete() bool
```

## Emulator Package

### Emulator

```go
type Emulator struct {}
type Palette struct {
    Colors [16]color.Color
}

func NewEmulator(dev device.Device, scale int) *Emulator
func (e *Emulator) SetWindowTitle(title string)
func (e *Emulator) SetFrameRate(fps int)
func (e *Emulator) ShowDebugInfo(show bool)
func (e *Emulator) SetBackgroundColor(c color.Color)
func (e *Emulator) SetPalette(p *Palette)
func (e *Emulator) Run() error
func (e *Emulator) GetDevice() device.Device
func (e *Emulator) GetFrameCount() int
func (e *Emulator) GetFPS() float64

func NewGrayscalePalette() *Palette
```

## Protocol Package

### SPI Bridge

```go
type SPIBridge struct {}
type Status struct {
    DCPin       bool
    CSPin       bool
    CommandMode bool
    LastCommand byte
}

func NewSPIBridge(dev device.Device) *SPIBridge
func (sb *SPIBridge) SetDC(state bool)
func (sb *SPIBridge) SetCS(state bool)
func (sb *SPIBridge) Write(data []byte) error
func (sb *SPIBridge) Reset() error
func (sb *SPIBridge) ReadData(length int) ([]byte, error)
func (sb *SPIBridge) SendInitSequence(sequence []byte) error
func (sb *SPIBridge) GetDevice() device.Device
func (sb *SPIBridge) GetStatus() Status
```

### Command Utilities

```go
type CommandInfo struct {
    Code        byte
    Name        string
    Description string
    DataBytes   int
}

type CommandBuilder struct {}
func NewCommandBuilder() *CommandBuilder
func (cb *CommandBuilder) AddCommand(code byte) *CommandBuilder
func (cb *CommandBuilder) AddData(data byte) *CommandBuilder
func (cb *CommandBuilder) AddBytes(data ...byte) *CommandBuilder
func (cb *CommandBuilder) Build() []byte
func (cb *CommandBuilder) Reset() *CommandBuilder

func GetCommandInfo(code byte) (CommandInfo, error)
func SSD1322InitSequence() []byte
func DrawPixelCommand(x, y, color byte) []byte
func FillScreenCommand(color byte) []byte
func ContrastCommand(level byte) []byte
func InversionCommand(inverted bool) []byte
func PowerCommand(on bool) []byte
```

## Color Values

```
0x00 - Black
0x01 - 0x07 - Dark shades
0x08 - 0x0E - Medium shades
0x0F - White
```

## Enumerations

### PixelFormat

```go
const (
    HorizontalNibble PixelFormat = iota  // SSD1322 native
    VerticalByte                         // SSD1306 style
    RGB888                               // 24-bit color
)
```

### TextAlignment

```go
const (
    AlignLeft TextAlignment = iota
    AlignCenter
    AlignRight
)
```
