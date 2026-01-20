package main

import (
	"log"
	"time"

	"github.com/flavioheleno/oled-emulator/animation"
	"github.com/flavioheleno/oled-emulator/device"
	"github.com/flavioheleno/oled-emulator/emulator"
	"github.com/flavioheleno/oled-emulator/graphics"
)

// Example: Animation - Demonstrates animated transitions and effects
func main() {
	// Create device
	dev := device.NewSSD1322(256, 64)

	// Create emulator
	emu := emulator.NewEmulator(dev, 2)
	emu.SetWindowTitle("OLED Emulator - Animation")
	emu.ShowDebugInfo(false)

	// Create framebuffer
	fb := graphics.NewFrameBuffer(dev)

	// Create animator
	animator := animation.NewAnimator(30)

	// Variables to track animation
	var x, y, radius float64 = 30, 32, 5

	// Animation 1: Moving circle
	tween1 := animation.NewTween(30, 226, 2*time.Second, animation.EaseInOutQuad)
	tween1.SetOnUpdate(func(value float64) {
		x = value
	})

	// Animation 2: Growing circle
	tween2 := animation.NewTween(5, 20, 2*time.Second, animation.EaseInOutCubic)
	tween2.SetOnUpdate(func(value float64) {
		radius = value
	})

	// Create parallel animation (both tweens run together)
	parallel := animation.NewParallelTween(tween1, tween2)

	// Add animation to animator
	frameCounter := 0
	animationFunc := func(frame int, dt float64) bool {
		frameCounter++

		// Clear display
		fb.Clear(0x00)

		// Draw title
		font := graphics.DefaultBitmapFont()
		font.DrawString(fb, 80, 5, "Animation", 0x0F)

		// Draw animated circle
		color := byte((int(x) + int(radius)) % 16)
		fb.DrawCircle(int(x), int(y), int(radius), color, true)

		// Draw progress bar (simple progress based on frame count)
		progress := float64(frameCounter) / 60.0 // 60 frames for full animation
		if progress > 1.0 {
			progress = 1.0
		}
		barWidth := int(200 * progress)
		fb.DrawRect(28, 50, barWidth, 3, 0x0A, true)
		fb.DrawRect(28, 50, 200, 3, 0x05, false)

		fb.Flush()

		// Update parallel animation
		return parallel.Update(dt)
	}

	animator.AddAnimation(animationFunc)
	animator.SetFrameRate(30)
	animator.Start()

	// Create a channel to handle emulator window events
	go func() {
		// This would normally be done through ebiten events,
		// but for now we'll just let the animator run
		<-time.After(10 * time.Second)
	}()

	// Run emulator
	if err := emu.Run(); err != nil {
		log.Fatalf("emulator error: %v", err)
	}

	animator.Stop()
}
