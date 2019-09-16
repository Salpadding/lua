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
	if a > 0 && math.IsInf(float64(b), 1) || a < 0 && math.IsInf(float64(b), -1) {
		return a
	}
	if a > 0 && math.IsInf(float64(b), -1) || a < 0 && math.IsInf(float64(b), 1) {
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

func toComparison(i int) types.Comparison {
	if i < 0 {
		return types.LessThan
	}
	if i > 0 {
		return types.GreaterThan
	}
	return types.Equal
}

func Compare(a, b Value) (types.Comparison, bool) {
	switch x := a.(type) {
	case String:
		bs, ok := b.(String)
		if ok {
			return toComparison(strings.Compare(string(x), string(bs))), true
		}
	case Integer:
		switch y := b.(type) {
		case Integer:
			return toComparison(int(x - y)), true
		case Float:
			return toComparison(sigMod(Float(x) - y)), true
		}
	case Float:
		switch y := b.(type) {
		case Integer:
			res, ok := Compare(b, a)
			return -res, ok
		case Float:
			return toComparison(sigMod(x - y)), true
		}
	}
	return 0, false
}

func Pow(a, b Value) (Value, bool) {
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return Float(math.Pow(float64(af), float64(bf))), true
}

func Div(a, b Value) (Value, bool) {
	af, ok := a.ToFloat()
	bf, ok2 := b.ToFloat()
	if !ok || !ok2 {
		return nil, false
	}
	return af / bf, true
}

func BitwiseAnd(a, b Value) (Value, bool) {
	ai, ok := a.ToInteger()
	bi, ok2 := b.ToInteger()
	if ok && ok2 {
		return ai & bi, true
	}
	return nil, false
}

func BitwiseXor(a, b Value) (Value, bool) {
	ai, ok := a.ToInteger()
	bi, ok2 := b.ToInteger()
	if ok && ok2 {
		return ai ^ bi, true
	}
	return nil, false
}

func ShiftLeft(a, b Value) (Value, bool) {
	ai, ok := a.ToInteger()
	bi, ok2 := b.ToInteger()
	if !ok || !ok2 {
		return nil, false
	}
	if bi >= 0 {
		return ai << uint64(bi), true
	}
	return ai >> uint64(-bi), true
}

func ShiftRight(a, b Value) (Value, bool) {
	ai, ok := a.ToInteger()
	bi, ok2 := b.ToInteger()
	if !ok || !ok2 {
		return nil, false
	}
	if bi >= 0 {
		return ai >> uint64(bi), true
	}
	return ai << uint64(-bi), true
}

func BitwiseOr(a, b Value) (Value, bool) {
	ai, ok := a.ToInteger()
	bi, ok2 := b.ToInteger()
	if ok && ok2 {
		return ai | bi, true
	}
	return nil, false
}

func equal(a, b Value) bool {
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

func Equal(a, b Value) (types.Comparison, bool) {
	ok := equal(a, b)
	if ok {
		return types.Equal, true
	}
	return 0, true
}

func Len(a Value) (Value, bool) {
	switch x := a.(type) {
	case String:
		return Integer(len(x)), true
	}
	return nil, false
}

func UnaryMinus(a Value) (Value, bool) {
	switch x := a.(type) {
	case Float:
		return -x, true
	case Integer:
		return -x, true
	default:
		return nil, false
	}
}

func BitwiseNot(a Value) (Value, bool) {
	ai, ok := a.ToInteger()
	if ok {
		return ^ai, true
	}
	return nil, false
}
