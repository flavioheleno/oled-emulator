package emulator

import (
	"fmt"
	"image/color"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Emulator represents the display emulator window
type Emulator struct {
	device         device.Device
	renderer       *VRAMRenderer
	screenImage    *ebiten.Image
	scale          int
	frameRate      int
	windowTitle    string
	backgroundColor color.Color
	showDebugInfo  bool
	frameCount     int
	lastFPS        float64
}

// NewEmulator creates a new emulator window
func NewEmulator(dev device.Device, scale int) *Emulator {
	return &Emulator{
		device:         dev,
		renderer:       NewVRAMRenderer(dev, scale),
		scale:          scale,
		frameRate:      60,
		windowTitle:    "OLED Display Emulator",
		backgroundColor: color.RGBA{R: 20, G: 20, B: 20, A: 255},
		showDebugInfo:  false,
		frameCount:     0,
	}
}

// SetWindowTitle sets the window title
func (e *Emulator) SetWindowTitle(title string) {
	e.windowTitle = title
}

// SetFrameRate sets the target frame rate
func (e *Emulator) SetFrameRate(fps int) {
	e.frameRate = fps
	ebiten.SetMaxTPS(fps)
}

// ShowDebugInfo enables/disables debug information display
func (e *Emulator) ShowDebugInfo(show bool) {
	e.showDebugInfo = show
}

// SetBackgroundColor sets the background color
func (e *Emulator) SetBackgroundColor(c color.Color) {
	e.backgroundColor = c
	e.renderer.SetBackgroundColor(c)
}

// SetPalette sets a custom color palette
func (e *Emulator) SetPalette(p *Palette) {
	e.renderer.SetPalette(p)
}

// Update implements the ebiten.Game Update method
func (e *Emulator) Update() error {
	e.frameCount++

	// Update FPS calculation every 30 frames
	if e.frameCount%30 == 0 {
		e.lastFPS = ebiten.ActualFPS()
	}

	return nil
}

// Draw implements the ebiten.Game Draw method
func (e *Emulator) Draw(screen *ebiten.Image) {
	// Clear screen with background color
	screen.Fill(e.backgroundColor)

	// Render VRAM to image
	e.screenImage = e.renderer.RenderFullScreen()

	// Draw the display at (0, 0)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.screenImage, op)

	// Draw debug info if enabled
	if e.showDebugInfo {
		e.drawDebugInfo(screen)
	}
}

// Layout implements the ebiten.Game Layout method
func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	width := e.device.Width() * e.scale
	height := e.device.Height() * e.scale
	return width, height
}

// drawDebugInfo draws debug information on screen
func (e *Emulator) drawDebugInfo(screen *ebiten.Image) {
	debugText := fmt.Sprintf(
		"FPS: %.1f\nFrame: %d\nDevice: %dx%d\nScale: %dx",
		e.lastFPS,
		e.frameCount,
		e.device.Width(),
		e.device.Height(),
		e.scale,
	)

	// Draw debug text
	ebitenutil.DebugPrintAt(screen, debugText, 5, 5)
}

// Run starts the emulator window
func (e *Emulator) Run() error {
	ebiten.SetWindowTitle(e.windowTitle)
	ebiten.SetMaxTPS(e.frameRate)

	return ebiten.RunGame(e)
}

// GetDevice returns the underlying device
func (e *Emulator) GetDevice() device.Device {
	return e.device
}

// GetFrameCount returns the number of frames rendered
func (e *Emulator) GetFrameCount() int {
	return e.frameCount
}

// GetFPS returns the current FPS
func (e *Emulator) GetFPS() float64 {
	return e.lastFPS
}
