package graphics

// Simplified TrueType font support - currently delegates to bitmap fonts
// Full TrueType rendering can be added later with golang.org/x/image/font

// TrueTypeFont is a placeholder for TrueType font support
// For now, delegates to bitmap font to keep implementation simple
type TrueTypeFont struct {
	bitmapFont *BitmapFont
	height     int
}

// NewTrueTypeFont creates a new TrueType font renderer
// This is a simplified implementation that uses a bitmap font
func NewTrueTypeFont(height int) *TrueTypeFont {
	bf := DefaultBitmapFont()

	return &TrueTypeFont{
		bitmapFont: bf,
		height:     height,
	}
}

// Height returns the font height
func (ttf *TrueTypeFont) Height() int {
	return ttf.bitmapFont.Height()
}

// DrawString draws text at the specified position
func (ttf *TrueTypeFont) DrawString(fb *FrameBuffer, x, y int, text string, color byte) (int, error) {
	return ttf.bitmapFont.DrawString(fb, x, y, text, color)
}

// MeasureString returns the width and height of text
func (ttf *TrueTypeFont) MeasureString(text string) (width, height int, err error) {
	return ttf.bitmapFont.MeasureString(text)
}

// GetGlyph returns glyph data for a character
func (ttf *TrueTypeFont) GetGlyph(ch rune) (GlyphData, error) {
	return ttf.bitmapFont.GetGlyph(ch)
}
