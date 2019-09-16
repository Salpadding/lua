package value

import (
	"math"
	"strings"

	"github.com/Salpadding/lua/types/value/types"
)

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

func Sub(a, b Value) (Value, bool) {
	ai, ok := a.(Integer)
	bi, ok2 := b.(Integer)
	if ok && ok2 {
		return ai - bi, true
	}
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return af - bf, true
}

func Mul(a, b Value) (Value, bool) {
	ai, ok := a.(Integer)
	bi, ok2 := a.(Integer)
	if ok && ok2 {
		return ai * bi, true
	}
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return af * bf, true
}

func Mod(a, b Value) (Value, bool) {
	ai, ok := a.(Integer)
	bi, ok2 := a.(Integer)
	if ok && ok2 {
		return IMod(ai, bi), true
	}
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return FMod(af, bf), true
}

func IDiv(a, b Value) (Value, bool) {
	ai, ok := a.(Integer)
	bi, ok2 := a.(Integer)
	if ok && ok2 {
		return IFloorDiv(ai, bi), true
	}
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return FFloorDiv(af, bf), true
}

func Negative(a Value) (Value, bool) {
	ai, ok := a.(Integer)
	if ok {
		return -ai, true
	}
	af, ok := a.ToFloat()
	if ok {
		return -af, true
	}
	return nil, false
}

func sigMod(f Float) int {
	if f < 0 {
		return -1
	}
	if f > 0 {
		return 1
	}
	return 0
}

func Compare(a, b Value) (int, bool) {
	switch x := a.(type) {
	case String:
		bs, ok := b.(String)
		if ok {
			return strings.Compare(string(x), string(bs)), true
		}
	case Integer:
		switch y := b.(type) {
		case Integer:
			return int(x - y), true
		case Float:
			return sigMod(Float(x) - y), true
		}
	case Float:
		switch y := b.(type) {
		case Integer:
			res, ok := Compare(b, a)
			return -res, ok
		case Float:
			return sigMod(x - y), true
		}
	}
	return 0, false
}

func Equal(a, b Value) bool {
	switch x := a.(type) {
	case *Nil:
		return b.Type() == types.Nil
	case Boolean:
		bb, ok := b.(Boolean)
		return ok && x == bb
	case String:
		bs, ok := b.(String)
		return ok && bs == x
	case Integer, Float:
		cmp, _ := Compare(a, b)
		return cmp == 0
	default:
		return a == b
	}
}
