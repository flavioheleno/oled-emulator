package graphics

import (
	"fmt"
	"image"
	"image/color"
)

// DrawImage draws an image to the framebuffer at the specified position
func DrawImage(fb *FrameBuffer, x, y int, img image.Image) error {
	if img == nil {
		return fmt.Errorf("image is nil")
	}

	bounds := img.Bounds()

	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			r, g, b, a := img.At(px, py).RGBA()

			// Skip fully transparent pixels
			if a == 0 {
				continue
			}

			// Convert RGB to grayscale
			gray := byte(((r >> 8) * 77 + (g >> 8) * 150 + (b >> 8) * 29) / 256)

			// Convert to 4-bit grayscale
			level := gray >> 4

			if level > 0 {
				screenX := x + px - bounds.Min.X
				screenY := y + py - bounds.Min.Y
				fb.SetPixel(screenX, screenY, level)
			}
		}
	}

	return nil
}

// DrawImageScaled draws a scaled image to the framebuffer
func DrawImageScaled(fb *FrameBuffer, x, y, w, h int, img image.Image) error {
	if img == nil {
		return fmt.Errorf("image is nil")
	}

	if w <= 0 || h <= 0 {
		return fmt.Errorf("invalid image dimensions: %dx%d", w, h)
	}

	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcWidth <= 0 || srcHeight <= 0 {
		return fmt.Errorf("source image has invalid dimensions")
	}

	// Use nearest-neighbor scaling
	for py := 0; py < h; py++ {
		for px := 0; px < w; px++ {
			// Calculate source pixel coordinates
			srcX := (px * srcWidth) / w
			srcY := (py * srcHeight) / h

			// Get pixel from source image
			r, g, b, a := img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY).RGBA()

			// Skip fully transparent pixels
			if a == 0 {
				continue
			}

			// Convert RGB to grayscale
			gray := byte(((r >> 8) * 77 + (g >> 8) * 150 + (b >> 8) * 29) / 256)

			// Convert to 4-bit grayscale
			level := gray >> 4

			if level > 0 {
				screenX := x + px
				screenY := y + py
				fb.SetPixel(screenX, screenY, level)
			}
		}
	}

	return nil
}

// ImageTiler provides tiling/repeating functionality for images
type ImageTiler struct {
	img image.Image
	w   int
	h   int
}

// NewImageTiler creates a new image tiler
func NewImageTiler(img image.Image) *ImageTiler {
	bounds := img.Bounds()
	return &ImageTiler{
		img: img,
		w:   bounds.Dx(),
		h:   bounds.Dy(),
	}
}

// DrawTiled draws a tiled pattern of the image
func (it *ImageTiler) DrawTiled(fb *FrameBuffer, x, y, w, h int) error {
	if it.w <= 0 || it.h <= 0 {
		return fmt.Errorf("tile dimensions invalid: %dx%d", it.w, it.h)
	}

	bounds := it.img.Bounds()

	for py := 0; py < h; py++ {
		for px := 0; px < w; px++ {
			// Get tiled coordinates
			tileX := px % it.w
			tileY := py % it.h

			// Get pixel from source image
			r, g, b, a := it.img.At(bounds.Min.X+tileX, bounds.Min.Y+tileY).RGBA()

			// Skip fully transparent pixels
			if a == 0 {
				continue
			}

			// Convert RGB to grayscale
			gray := byte(((r >> 8) * 77 + (g >> 8) * 150 + (b >> 8) * 29) / 256)

			// Convert to 4-bit grayscale
			level := gray >> 4

			if level > 0 {
				screenX := x + px
				screenY := y + py
				fb.SetPixel(screenX, screenY, level)
			}
		}
	}

	return nil
}

// ConvertToGrayscale converts an image to 4-bit grayscale
func ConvertToGrayscale(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()

			// Convert to grayscale using standard luminosity formula
			gray := uint8(((r >> 8) * 77 + (g >> 8) * 150 + (b >> 8) * 29) / 256)

			// Convert to 4-bit and back to 8-bit for display
			level := (gray >> 4) << 4

			dst.Set(x, y, color.RGBA{R: level, G: level, B: level, A: uint8(a >> 8)})
		}
	}

	return dst
}

// ConvertToBitmap converts an image to 1-bit black and white using threshold
func ConvertToBitmap(src image.Image, threshold uint8) image.Image {
	bounds := src.Bounds()
	dst := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()

			// Convert to grayscale
			gray := uint8(((r >> 8) * 77 + (g >> 8) * 150 + (b >> 8) * 29) / 256)

			// Apply threshold
			if gray > threshold {
				dst.Set(x, y, color.White)
			} else {
				dst.Set(x, y, color.Black)
			}
		}
	}

	return dst
}
