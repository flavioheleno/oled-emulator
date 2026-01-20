package graphics

import (
	"testing"

	"github.com/flavioheleno/oled-emulator/device"
)

func TestBitmapFontCreation(t *testing.T) {
	bf := NewBitmapFont(5, 7, 6)

	if bf.Height() != 7 {
		t.Errorf("expected height 7, got %d", bf.Height())
	}
}

func TestBitmapFontMeasure(t *testing.T) {
	bf := NewBitmapFont(5, 7, 6)

	width, height, err := bf.MeasureString("Hello")
	if err != nil {
		t.Fatalf("measure failed: %v", err)
	}

	if width != 30 { // 5 characters * 6 advance
		t.Errorf("expected width 30, got %d", width)
	}

	if height != 7 {
		t.Errorf("expected height 7, got %d", height)
	}
}

func TestDefaultBitmapFont(t *testing.T) {
	bf := DefaultBitmapFont()

	if bf.Height() != 7 {
		t.Errorf("expected height 7, got %d", bf.Height())
	}

	// Test that we can get glyphs
	glyph, err := bf.GetGlyph('A')
	if err != nil {
		t.Fatalf("failed to get glyph for 'A': %v", err)
	}

	if glyph.Width != 5 || glyph.Height != 7 {
		t.Errorf("expected glyph size 5x7, got %dx%d", glyph.Width, glyph.Height)
	}
}

func TestBitmapFontDrawString(t *testing.T) {
	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	bf := DefaultBitmapFont()
	width, err := bf.DrawString(fb, 10, 10, "H", 0x0F)

	if err != nil {
		t.Fatalf("draw string failed: %v", err)
	}

	if width != 6 {
		t.Errorf("expected width 6, got %d", width)
	}
}

func TestTextRenderer(t *testing.T) {
	bf := DefaultBitmapFont()
	tr := NewTextRenderer(bf)

	opts := DefaultTextOptions()
	opts.Color = 0x08
	tr.SetOptions(opts)

	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	width, err := tr.DrawText(fb, 20, 20, "Test")
	if err != nil {
		t.Fatalf("draw text failed: %v", err)
	}

	if width != 24 { // 4 characters * 6 advance
		t.Errorf("expected width 24, got %d", width)
	}
}

func TestMultilineText(t *testing.T) {
	bf := DefaultBitmapFont()
	tr := NewTextRenderer(bf)

	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	text := "Line1\nLine2\nLine3"
	if err := tr.DrawMultilineText(fb, 10, 10, text); err != nil {
		t.Fatalf("draw multiline text failed: %v", err)
	}
}

func TestMeasureMultilineText(t *testing.T) {
	bf := DefaultBitmapFont()
	tr := NewTextRenderer(bf)

	width, height, err := tr.MeasureMultilineText("Line1\nLine2")
	if err != nil {
		t.Fatalf("measure multiline text failed: %v", err)
	}

	if height != 14 { // 2 lines * 7 height (no line spacing)
		t.Errorf("expected height 14, got %d", height)
	}

	if width <= 0 {
		t.Errorf("expected positive width, got %d", width)
	}
}

func TestAlignedTextDrawer(t *testing.T) {
	bf := DefaultBitmapFont()
	atd := NewAlignedTextDrawer(bf)

	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	// Test left align
	if err := atd.DrawAlignedText(fb, 10, 10, "Left", AlignLeft, 0x0F); err != nil {
		t.Fatalf("left align failed: %v", err)
	}

	// Test center align
	if err := atd.DrawAlignedText(fb, 128, 20, "Center", AlignCenter, 0x0F); err != nil {
		t.Fatalf("center align failed: %v", err)
	}

	// Test right align
	if err := atd.DrawAlignedText(fb, 250, 30, "Right", AlignRight, 0x0F); err != nil {
		t.Fatalf("right align failed: %v", err)
	}
}

func TestCenteredText(t *testing.T) {
	bf := DefaultBitmapFont()
	atd := NewAlignedTextDrawer(bf)

	dev := device.NewSSD1322(256, 64)
	fb := NewFrameBuffer(dev)

	if err := atd.DrawCenteredText(fb, 128, 32, "Center", 0x0F); err != nil {
		t.Fatalf("centered text failed: %v", err)
	}
}
