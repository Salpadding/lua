package value

import "math"

func FloatToInteger(f Float) (Integer, bool) {
	// todo: correct?
	i := Integer(f)
	return i, Float(i) == f
}

// a % b == a - ((a // b) * b)
func IMod(a, b Integer) Integer {
	return a - IFloorDiv(a, b)*b
}

// a % b == a - ((a // b) * b)
// lua-5.3.4/src/llimits.h#luai_nummod
func FMod(a, b Float) Float {
	if a > 0 && math.IsInf(float64(b), 1) || a < 0 && math.IsInf(b, -1) {
		return a
	}
	if a > 0 && math.IsInf(float64(b), -1) || a < 0 && math.IsInf(b, 1) {
		return b
	}
	return a - Float(math.Floor(float64(a/b)))*b
}

func IFloorDiv(a, b Integer) Integer {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

func FFloorDiv(a, b Float) Float {
	return Float(math.Floor(float64(a / b)))
}

func ShiftLeft(a, n Integer) Integer {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}

func ShiftRight(a, n Integer) Integer {
	if n >= 0 {
		return Integer(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

func Add(a, b Value) (Value, bool) {
	ai, ok := a.(Integer)
	bi, ok2 := b.(Integer)
	if ok && ok2 {
		return ai + bi, true
	}
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return af + bf, true
}
