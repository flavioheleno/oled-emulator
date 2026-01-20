package graphics

import (
	"fmt"
)

// Font defines the interface for text rendering
type Font interface {
	// DrawString draws text at the specified position
	// Returns the width of the drawn text
	DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error)

	// MeasureString returns the width and height of text without drawing
	MeasureString(text string) (width, height int, err error)

	// Height returns the font height in pixels
	Height() int

	// GetGlyph returns the glyph data for a character
	GetGlyph(ch rune) (GlyphData, error)
}

// GlyphData contains information about a single character
type GlyphData struct {
	Width    int    // Glyph width in pixels
	Height   int    // Glyph height in pixels
	AdvanceX int    // Pixels to advance after drawing
	BearingX int    // Offset from cursor position to glyph left
	BearingY int    // Offset from cursor position to glyph top
	Data     []byte // Glyph bitmap data (1 bit per pixel, packed horizontally)
}

// TextAlignment defines text alignment modes
type TextAlignment int

const (
	AlignLeft TextAlignment = iota
	AlignCenter
	AlignRight
)

// TextOptions holds text rendering options
type TextOptions struct {
	Alignment   TextAlignment
	LineSpacing int
	CharSpacing int
	Color       byte
}

// DefaultTextOptions returns default text rendering options
func DefaultTextOptions() TextOptions {
	return TextOptions{
		Alignment:   AlignLeft,
		LineSpacing: 0,
		CharSpacing: 0,
		Color:       0x0F,
	}
}

// TextRenderer provides high-level text rendering with layout support
type TextRenderer struct {
	font Font
	opts TextOptions
}

// NewTextRenderer creates a new text renderer
func NewTextRenderer(font Font) *TextRenderer {
	return &TextRenderer{
		font: font,
		opts: DefaultTextOptions(),
	}
}

// SetOptions sets text rendering options
func (tr *TextRenderer) SetOptions(opts TextOptions) {
	tr.opts = opts
}

// DrawText draws text with current options
func (tr *TextRenderer) DrawText(fb *FrameBuffer, x, y int, text string) (int, error) {
	return tr.font.DrawString(fb, x, y, text, tr.opts.Color)
}

// DrawMultilineText draws multiple lines of text
func (tr *TextRenderer) DrawMultilineText(fb *FrameBuffer, x, y int, text string) error {
	// Split text by newlines
	lines := splitLines(text)
	currentY := y

	for _, line := range lines {
		if _, err := tr.font.DrawString(fb, x, currentY, line, tr.opts.Color); err != nil {
			return fmt.Errorf("failed to draw line: %w", err)
		}

		currentY += tr.font.Height() + tr.opts.LineSpacing
	}

	return nil
}

// MeasureMultilineText measures the bounding box of multiline text
func (tr *TextRenderer) MeasureMultilineText(text string) (width, height int, err error) {
	lines := splitLines(text)
	if len(lines) == 0 {
		return 0, 0, nil
	}

	maxWidth := 0
	for _, line := range lines {
		w, _, err := tr.font.MeasureString(line)
		if err != nil {
			return 0, 0, err
		}
		if w > maxWidth {
			maxWidth = w
		}
	}

	totalHeight := tr.font.Height()*len(lines) + tr.opts.LineSpacing*(len(lines)-1)

	return maxWidth, totalHeight, nil
}

// Helper function to split text by newlines
func splitLines(text string) []string {
	var lines []string
	var currentLine string

	for _, ch := range text {
		if ch == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else if ch != '\r' {
			currentLine += string(ch)
		}
	}

	if currentLine != "" || len(text) > 0 && text[len(text)-1] == '\n' {
		lines = append(lines, currentLine)
	}

	if len(lines) == 0 {
		lines = []string{""}
	}

	return lines
}

// AlignedTextDrawer handles text alignment and positioning
type AlignedTextDrawer struct {
	renderer *TextRenderer
}

// NewAlignedTextDrawer creates a new aligned text drawer
func NewAlignedTextDrawer(font Font) *AlignedTextDrawer {
	return &AlignedTextDrawer{
		renderer: NewTextRenderer(font),
	}
}

// DrawAlignedText draws text with alignment
func (atd *AlignedTextDrawer) DrawAlignedText(fb *FrameBuffer, x, y int, text string, alignment TextAlignment, color byte) error {
	width, _, err := atd.renderer.font.MeasureString(text)
	if err != nil {
		return err
	}

	var drawX int
	switch alignment {
	case AlignLeft:
		drawX = x
	case AlignCenter:
		drawX = x - width/2
	case AlignRight:
		drawX = x - width
	}

	opts := atd.renderer.opts
	opts.Color = color
	atd.renderer.SetOptions(opts)

	_, err = atd.renderer.font.DrawString(fb, drawX, y, text, color)
	return err
}

// DrawCenteredText is a convenience function for center-aligned text
func (atd *AlignedTextDrawer) DrawCenteredText(fb *FrameBuffer, x, y int, text string, color byte) error {
	return atd.DrawAlignedText(fb, x, y, text, AlignCenter, color)
}

// DrawRightAlignedText is a convenience function for right-aligned text
func (atd *AlignedTextDrawer) DrawRightAlignedText(fb *FrameBuffer, x, y int, text string, color byte) error {
	return atd.DrawAlignedText(fb, x, y, text, AlignRight, color)
}
