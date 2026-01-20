package emulator

import (
	"image"
	"image/color"

	"github.com/flavioheleno/oled-emulator/device"
	"github.com/hajimehoshi/ebiten/v2"
)

// Palette defines color mapping for different grayscale levels
type Palette struct {
	Colors [16]color.Color
}

// NewGrayscalePalette creates a standard grayscale palette
func NewGrayscalePalette() *Palette {
	p := &Palette{}

	// Create grayscale levels from black to white
	for i := 0; i < 16; i++ {
		level := uint8((i * 255) / 15)
		// OLED-style: yellow tint for bright pixels, slight blue tint for dim
		if i < 8 {
			// Darker pixels: slight blue tint
			p.Colors[i] = color.RGBA{
				R: level * 200 / 255,
				G: level * 150 / 255,
				B: level * 255 / 255,
				A: 255,
			}
		} else {
			// Brighter pixels: yellow tint (characteristic of OLEDs)
			p.Colors[i] = color.RGBA{
				R: level,
				G: level * 200 / 255,
				B: level * 100 / 255,
				A: 255,
			}
		}
	}

	// Ensure color 0 is pure black for off pixels
	p.Colors[0] = color.RGBA{R: 20, G: 20, B: 20, A: 255}

	return p
}

// VRAMRenderer converts device VRAM to a renderable image
type VRAMRenderer struct {
	device        device.Device
	palette       *Palette
	scale         int
	lastDirtyX0   int
	lastDirtyY0   int
	lastDirtyX1   int
	lastDirtyY1   int
	backgroundColor color.Color
}

// NewVRAMRenderer creates a new VRAM renderer
func NewVRAMRenderer(dev device.Device, scale int) *VRAMRenderer {
	return &VRAMRenderer{
		device:          dev,
		palette:         NewGrayscalePalette(),
		scale:           scale,
		backgroundColor: color.RGBA{R: 20, G: 20, B: 20, A: 255},
	}
}

// SetPalette sets a custom palette
func (vr *VRAMRenderer) SetPalette(p *Palette) {
	vr.palette = p
}

// SetBackgroundColor sets the background color (off pixel color)
func (vr *VRAMRenderer) SetBackgroundColor(c color.Color) {
	vr.backgroundColor = c
}

// RenderToImage converts VRAM to an ebiten.Image
func (vr *VRAMRenderer) RenderToImage() *ebiten.Image {
	width := vr.device.Width()
	height := vr.device.Height()

	// Create image with scaled dimensions
	img := ebiten.NewImage(width*vr.scale, height*vr.scale)

	// Get dirty region for optimization
	dirtyX0, dirtyY0, dirtyX1, dirtyY1 := vr.device.GetDirtyRegion()

	// If no dirty region, render full screen
	if dirtyX0 == -1 {
		dirtyX0 = 0
		dirtyY0 = 0
		dirtyX1 = width - 1
		dirtyY1 = height - 1
	}

	// Render pixels in dirty region
	for y := dirtyY0; y <= dirtyY1; y++ {
		for x := dirtyX0; x <= dirtyX1; x++ {
			pixel, err := vr.device.GetPixel(x, y)
			if err != nil {
				pixel = 0
			}

			// Ensure pixel is 4-bit
			pixel = pixel & 0x0F

			// Get color from palette
			pixelColor := vr.palette.Colors[pixel]

			// Draw scaled pixel
			rect := image.Rect(
				x*vr.scale, y*vr.scale,
				(x+1)*vr.scale, (y+1)*vr.scale,
			)

			for py := rect.Min.Y; py < rect.Max.Y; py++ {
				for px := rect.Min.X; px < rect.Max.X; px++ {
					img.Set(px, py, pixelColor)
				}
			}
		}
	}

	return img
}

// RenderFullScreen renders the entire VRAM regardless of dirty state
func (vr *VRAMRenderer) RenderFullScreen() *ebiten.Image {
	width := vr.device.Width()
	height := vr.device.Height()

	img := ebiten.NewImage(width*vr.scale, height*vr.scale)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel, err := vr.device.GetPixel(x, y)
			if err != nil {
				pixel = 0
			}

			pixel = pixel & 0x0F
			pixelColor := vr.palette.Colors[pixel]

			rect := image.Rect(
				x*vr.scale, y*vr.scale,
				(x+1)*vr.scale, (y+1)*vr.scale,
			)

			for py := rect.Min.Y; py < rect.Max.Y; py++ {
				for px := rect.Min.X; px < rect.Max.X; px++ {
					img.Set(px, py, pixelColor)
				}
			}
		}
	}

	return img
}
