package animation

import "math"

// EasingFunc defines an easing function type
// Input t is normalized time (0 to 1)
// Output is the eased value (0 to 1)
type EasingFunc func(t float64) float64

// Clamp normalizes t to [0, 1]
func clamp(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// Linear easing (no acceleration)
func Linear(t float64) float64 {
	return clamp(t)
}

// EaseInQuad accelerating from zero velocity
func EaseInQuad(t float64) float64 {
	t = clamp(t)
	return t * t
}

// EaseOutQuad decelerating to zero velocity
func EaseOutQuad(t float64) float64 {
	t = clamp(t)
	return -t * (t - 2)
}

// EaseInOutQuad acceleration until halfway, then deceleration
func EaseInOutQuad(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseInCubic accelerating from zero velocity (cubic)
func EaseInCubic(t float64) float64 {
	t = clamp(t)
	return t * t * t
}

// EaseOutCubic decelerating to zero velocity (cubic)
func EaseOutCubic(t float64) float64 {
	t = clamp(t)
	t--
	return t*t*t + 1
}

// EaseInOutCubic acceleration until halfway, then deceleration (cubic)
func EaseInOutCubic(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return 4 * t * t * t
	}
	t--
	return 1 + t*t*t*4
}

// EaseInQuart accelerating from zero velocity (quartic)
func EaseInQuart(t float64) float64 {
	t = clamp(t)
	return t * t * t * t
}

// EaseOutQuart decelerating to zero velocity (quartic)
func EaseOutQuart(t float64) float64 {
	t = clamp(t)
	t--
	return -(t*t*t*t - 1)
}

// EaseInOutQuart acceleration until halfway, then deceleration (quartic)
func EaseInOutQuart(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return 8 * t * t * t * t
	}
	t--
	return -(8*t*t*t*t - 1)
}

// EaseInQuint accelerating from zero velocity (quintic)
func EaseInQuint(t float64) float64 {
	t = clamp(t)
	return t * t * t * t * t
}

// EaseOutQuint decelerating to zero velocity (quintic)
func EaseOutQuint(t float64) float64 {
	t = clamp(t)
	t--
	return t*t*t*t*t + 1
}

// EaseInOutQuint acceleration until halfway, then deceleration (quintic)
func EaseInOutQuint(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return 16 * t * t * t * t * t
	}
	t--
	return 1 + 16*t*t*t*t*t
}

// EaseInSine accelerating from zero velocity (sinusoidal)
func EaseInSine(t float64) float64 {
	t = clamp(t)
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// EaseOutSine decelerating to zero velocity (sinusoidal)
func EaseOutSine(t float64) float64 {
	t = clamp(t)
	return math.Sin(math.Pi*t) / 2
}

// EaseInOutSine acceleration until halfway, then deceleration (sinusoidal)
func EaseInOutSine(t float64) float64 {
	t = clamp(t)
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// EaseInExpo accelerating from zero velocity (exponential)
func EaseInExpo(t float64) float64 {
	t = clamp(t)
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*t-10)
}

// EaseOutExpo decelerating to zero velocity (exponential)
func EaseOutExpo(t float64) float64 {
	t = clamp(t)
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

// EaseInOutExpo acceleration until halfway, then deceleration (exponential)
func EaseInOutExpo(t float64) float64 {
	t = clamp(t)
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	if t < 0.5 {
		return math.Pow(2, 20*t-10) / 2
	}
	return (2 - math.Pow(2, -20*t+10)) / 2
}

// EaseInCirc accelerating from zero velocity (circular)
func EaseInCirc(t float64) float64 {
	t = clamp(t)
	return -(math.Sqrt(1-t*t) - 1)
}

// EaseOutCirc decelerating to zero velocity (circular)
func EaseOutCirc(t float64) float64 {
	t = clamp(t)
	t--
	return math.Sqrt(1 - t*t)
}

// EaseInOutCirc acceleration until halfway, then deceleration (circular)
func EaseInOutCirc(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return (1 - math.Sqrt(1-(2*t)*(2*t))) / 2
	}
	return (math.Sqrt(1-((-2*t+2)*(-2*t+2))) + 1) / 2
}

// EaseInBack accelerating with overshoot
func EaseInBack(t float64) float64 {
	t = clamp(t)
	c1 := 1.70158
	c3 := c1 + 1
	return c3*t*t*t - c1*t*t
}

// EaseOutBack decelerating with overshoot
func EaseOutBack(t float64) float64 {
	t = clamp(t)
	c1 := 1.70158
	c3 := c1 + 1
	return 1 + c3*(t-1)*(t-1)*(t-1) + c1*(t-1)*(t-1)
}

// EaseInOutBack acceleration and deceleration with overshoot
func EaseInOutBack(t float64) float64 {
	t = clamp(t)
	c1 := 1.70158
	c2 := c1 * 1.525
	if t < 0.5 {
		return (math.Pow(2*t, 2) * ((c2+1)*2*t - c2)) / 2
	}
	return (math.Pow(2*t-2, 2)*(c2+1)*(t*2-2)+c2)/2 + 1
}

// EaseInElastic elastically accelerating from zero velocity
func EaseInElastic(t float64) float64 {
	t = clamp(t)
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	c4 := (2 * math.Pi) / 3
	return -(math.Pow(2, 10*t-10)) * math.Sin((t*10-10.75)*c4)
}

// EaseOutElastic elastically decelerating to zero velocity
func EaseOutElastic(t float64) float64 {
	t = clamp(t)
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	c4 := (2 * math.Pi) / 3
	return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*c4) + 1
}

// EaseInOutElastic elastic acceleration and deceleration
func EaseInOutElastic(t float64) float64 {
	t = clamp(t)
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	c5 := (2 * math.Pi) / 4.5
	if t < 0.5 {
		return -(math.Pow(2, 20*t-10) * math.Sin((20*t-11.125)*c5)) / 2
	}
	return (math.Pow(2, -20*t+10)*math.Sin((20*t-11.125)*c5))/2 + 1
}

// EaseInBounce bouncing acceleration from zero velocity
func EaseInBounce(t float64) float64 {
	return 1 - EaseOutBounce(1-clamp(t))
}

// EaseOutBounce bouncing deceleration to zero velocity
func EaseOutBounce(t float64) float64 {
	t = clamp(t)
	n1 := 7.5625
	d1 := 2.75

	if t < 1/d1 {
		return n1 * t * t
	} else if t < 2/d1 {
		t -= 1.5 / d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t -= 2.25 / d1
		return n1*t*t + 0.9375
	} else {
		t -= 2.625 / d1
		return n1*t*t + 0.984375
	}
}

// EaseInOutBounce bouncing acceleration and deceleration
func EaseInOutBounce(t float64) float64 {
	t = clamp(t)
	if t < 0.5 {
		return (1 - EaseOutBounce(1-2*t)) / 2
	}
	return (1 + EaseOutBounce(2*t-1)) / 2
}
