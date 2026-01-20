package graphics

import (
	"math"
)

// DrawLineBresenham draws a line using Bresenham's algorithm
func DrawLineBresenham(fb *FrameBuffer, x0, y0, x1, y1 int, color byte, setPixel func(int, int, byte)) {
	// Handle line clipping and drawing
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := sign(x1 - x0)
	sy := sign(y1 - y0)
	err := dx - dy

	x, y := x0, y0

	for {
		setPixel(x, y, color)

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// DrawCircle draws a circle using midpoint algorithm
func DrawCircle(fb *FrameBuffer, cx, cy, r int, color byte, filled bool, setPixel func(int, int, byte)) {
	if filled {
		DrawFilledCircle(fb, cx, cy, r, color, setPixel)
		return
	}

	DrawCircleOutline(fb, cx, cy, r, color, setPixel)
}

// DrawCircleOutline draws the outline of a circle
func DrawCircleOutline(fb *FrameBuffer, cx, cy, r int, color byte, setPixel func(int, int, byte)) {
	if r <= 0 {
		return
	}

	x := 0
	y := r
	d := 3 - 2*r

	for x <= y {
		// Draw 8 symmetric points
		setPixel(cx+x, cy+y, color)
		setPixel(cx-x, cy+y, color)
		setPixel(cx+x, cy-y, color)
		setPixel(cx-x, cy-y, color)
		setPixel(cx+y, cy+x, color)
		setPixel(cx-y, cy+x, color)
		setPixel(cx+y, cy-x, color)
		setPixel(cx-y, cy-x, color)

		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// DrawFilledCircle draws a filled circle
func DrawFilledCircle(fb *FrameBuffer, cx, cy, r int, color byte, setPixel func(int, int, byte)) {
	if r <= 0 {
		return
	}

	x := 0
	y := r
	d := 3 - 2*r

	for x <= y {
		// Draw horizontal lines
		drawHorizontalLine(cx-x, cx+x, cy+y, color, setPixel)
		drawHorizontalLine(cx-x, cx+x, cy-y, color, setPixel)
		drawHorizontalLine(cx-y, cx+y, cy+x, color, setPixel)
		drawHorizontalLine(cx-y, cx+y, cy-x, color, setPixel)

		if d < 0 {
			d = d + 4*x + 6
		} else {
			d = d + 4*(x-y) + 10
			y--
		}
		x++
	}
}

// drawHorizontalLine draws a horizontal line from x1 to x2 at y
func drawHorizontalLine(x1, x2, y int, color byte, setPixel func(int, int, byte)) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		setPixel(x, y, color)
	}
}

// DrawRect draws a rectangle
func DrawRect(fb *FrameBuffer, x, y, w, h int, color byte, filled bool, setPixel func(int, int, byte)) {
	if w < 0 || h < 0 {
		return
	}

	x1 := x + w
	y1 := y + h

	if filled {
		for py := y; py < y1; py++ {
			drawHorizontalLine(x, x1-1, py, color, setPixel)
		}
	} else {
		// Top and bottom lines
		drawHorizontalLine(x, x1-1, y, color, setPixel)
		drawHorizontalLine(x, x1-1, y1-1, color, setPixel)

		// Left and right lines
		for py := y; py < y1; py++ {
			setPixel(x, py, color)
			setPixel(x1-1, py, color)
		}
	}
}

// DrawEllipse draws an ellipse using midpoint algorithm
func DrawEllipse(fb *FrameBuffer, cx, cy, rx, ry int, color byte, filled bool, setPixel func(int, int, byte)) {
	if rx <= 0 || ry <= 0 {
		return
	}

	if filled {
		DrawFilledEllipse(fb, cx, cy, rx, ry, color, setPixel)
		return
	}

	DrawEllipseOutline(fb, cx, cy, rx, ry, color, setPixel)
}

// DrawEllipseOutline draws the outline of an ellipse
func DrawEllipseOutline(fb *FrameBuffer, cx, cy, rx, ry int, color byte, setPixel func(int, int, byte)) {
	x := rx
	y := 0
	dx := ry * ry * (1 - 2*rx)
	dy := rx * rx
	decision := ry*ry - rx*rx*ry + rx*rx/4

	for x >= y {
		// Draw 4 symmetric points
		setPixel(cx+x, cy+y, color)
		setPixel(cx-x, cy+y, color)
		setPixel(cx+x, cy-y, color)
		setPixel(cx-x, cy-y, color)

		if decision < 0 {
			decision += dy
			dy += 2 * rx * rx
		} else {
			decision += dx + dy
			dx += 2 * ry * ry
			x--
		}
		y++
	}
}

// DrawFilledEllipse draws a filled ellipse
func DrawFilledEllipse(fb *FrameBuffer, cx, cy, rx, ry int, color byte, setPixel func(int, int, byte)) {
	x := rx
	y := 0
	dx := ry * ry * (1 - 2*rx)
	dy := rx * rx
	decision := ry*ry - rx*rx*ry + rx*rx/4

	for x >= y {
		// Draw horizontal lines
		drawHorizontalLine(cx-x, cx+x, cy+y, color, setPixel)
		drawHorizontalLine(cx-x, cx+x, cy-y, color, setPixel)

		if decision < 0 {
			decision += dy
			dy += 2 * rx * rx
		} else {
			decision += dx + dy
			dx += 2 * ry * ry
			x--
		}
		y++
	}
}

// DrawTriangle draws a triangle
func DrawTriangle(fb *FrameBuffer, x1, y1, x2, y2, x3, y3 int, color byte, filled bool, setPixel func(int, int, byte)) {
	if filled {
		DrawFilledTriangle(fb, x1, y1, x2, y2, x3, y3, color, setPixel)
		return
	}

	DrawLineBresenham(fb, x1, y1, x2, y2, color, setPixel)
	DrawLineBresenham(fb, x2, y2, x3, y3, color, setPixel)
	DrawLineBresenham(fb, x3, y3, x1, y1, color, setPixel)
}

// DrawFilledTriangle draws a filled triangle using barycentric coordinates
func DrawFilledTriangle(fb *FrameBuffer, x1, y1, x2, y2, x3, y3 int, color byte, setPixel func(int, int, byte)) {
	// Find bounding box
	minX := min(x1, min(x2, x3))
	maxX := max(x1, max(x2, x3))
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	// Compute vectors
	v0x := x3 - x1
	v0y := y3 - y1
	v1x := x2 - x1
	v1y := y2 - y1

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot11 := v1x*v1x + v1y*v1y
	invDenom := float64(1) / float64(dot00*dot11-dot01*dot01)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			v2x := x - x1
			v2y := y - y1

			dot02 := v0x*v2x + v0y*v2y
			dot12 := v1x*v2x + v1y*v2y

			u := (float64(dot11*dot02-dot01*dot12)) * invDenom
			v := (float64(dot00*dot12-dot01*dot02)) * invDenom

			if u >= 0 && v >= 0 && u+v < 1 {
				setPixel(x, y, color)
			}
		}
	}
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp clamps value between min and max
func Clamp(value, minVal, maxVal int) int {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Lerp performs linear interpolation
func Lerp(a, b float64, t float64) float64 {
	return a + (b-a)*t
}

// Map maps a value from one range to another
func Map(value, inMin, inMax, outMin, outMax float64) float64 {
	return Lerp(outMin, outMax, (value-inMin)/(inMax-inMin))
}

// Distance calculates the distance between two points
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
