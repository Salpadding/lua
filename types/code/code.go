package code

const (
	MaxArgBx  = 1<<18 - 1     // 2^18 - 1 = 262143
	MaxArgsBx = MaxArgBx >> 1 // 262143 / 2 = 131071
)

type OpMode int

/* OpMode */
/* basic instruction format */
const (
	IABC  OpMode = iota // [  B:9  ][  C:9  ][ A:8  ][OP:6]
	IABx                // [      Bx:18     ][ A:8  ][OP:6]
	IAsBx               // [     sBx:18     ][ A:8  ][OP:6]
	IAx                 // [           Ax:26        ][OP:6]
)

type OpArgMask int

/* OpArgMask */
const (
	OpArgN OpArgMask = iota /* argument is not used */
	OpArgU                  /* argument is used */
	OpArgR                  /* argument is a register or a jump offset */
	OpArgK                  /* argument is a constant or register/constant */
)

type Type int

/* OpCode */
const (
	Move Type = iota
	LoadK
	LoadKX
	LoadBool
	LoadNil
	GETUPVAL
	GETTABUP
	GETTABLE
	SETTABUP
	SETUPVAL
	SETTABLE
	NEWTABLE
	SELF
	Add
	Sub
	Mul
	Mod
	Pow
	Div
	IDiv
	BitwiseAnd
	BitwiseOr
	BitwiseXor
	ShiftLeft
	ShiftRight
	UnaryMinus
	BitwiseNot
	LogicalNot
	Len
	Concat
	Jmp
	Equal
	LessThan
	LessThanOrEqual
	Test
	TestSet
	CALL
	TAILCALL
	Return
	ForLoop
	ForPrep
	TFORCALL
	TFORLOOP
	SETLIST
	CLOSURE
	VARARG
	EXTRAARG
)

func init() {
	for i := range OpCodes {
		OpCodes[i].Type = Type(i)
	}
}

type OpCode struct {
	Type     Type
	TestFlag byte
	SetAFlag byte
	ArgBMode OpArgMask
	ArgCMode OpArgMask
	OpMode   OpMode
	Name     string
}

var OpCodes = []*OpCode{
	/*     T  A    B       C     mode         name       action */
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "Move    "},              // R(A) := R(B)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgN, OpMode: IABx /* */, Name: "LoadK   "},              // R(A) := Kst(Bx)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgN, ArgCMode: OpArgN, OpMode: IABx /* */, Name: "LoadKX  "},              // R(A) := Kst(extra arg)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "LoadBool"},              // R(A) := (bool)B; if (C) pc++
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "LoadNil "},              // R(A), R(A+1), ..., R(A+B) := nil
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "GETUPVAL"},              // R(A) := UpValue[B]
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "GETTABUP"},              // R(A) := UpValue[B][RK(C)]
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "GETTABLE"},              // R(A) := R(B)[RK(C)]
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "SETTABUP"},              // UpValue[A][RK(B)] := RK(C)
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "SETUPVAL"},              // UpValue[B] := R(A)
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "SETTABLE"},              // R(A)[RK(B)] := RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "NEWTABLE"},              // R(A) := {} (size = B,C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "SELF    "},              // R(A+1) := R(B); R(A) := R(B)[RK(C)]
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Add     "},              // R(A) := RK(B) + RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Sub     "},              // R(A) := RK(B) - RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Mul     "},              // R(A) := RK(B) * RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Mod     "},              // R(A) := RK(B) % RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Pow     "},              // R(A) := RK(B) ^ RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Div     "},              // R(A) := RK(B) / RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "IDiv    "},              // R(A) := RK(B) // RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "BitwiseAnd    "},        // R(A) := RK(B) & RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "BitwiseOr     "},        // R(A) := RK(B) | RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "BitwiseXor    "},        // R(A) := RK(B) ~ RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "ShiftLeft     "},        // R(A) := RK(B) << RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "ShiftRight     "},       // R(A) := RK(B) >> RK(C)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "UnaryMinus     "},       // R(A) := -R(B)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "BitwiseNot    "},        // R(A) := ~R(B)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "LogicalNot     "},       // R(A) := not R(B)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "Len     "},              // R(A) := length of R(B)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgR, OpMode: IABC /* */, Name: "Concat  "},              // R(A) := R(B).. ... ..R(C)
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IAsBx /**/, Name: "Jmp     "},              // pc+=sBx; if (A) close all upvalues >= R(A - 1)
	{TestFlag: 1, SetAFlag: 0, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "Equal      "},           // if ((RK(B) == RK(C)) ~= A) then pc++
	{TestFlag: 1, SetAFlag: 0, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "LessThan      "},        // if ((RK(B) <  RK(C)) ~= A) then pc++
	{TestFlag: 1, SetAFlag: 0, ArgBMode: OpArgK, ArgCMode: OpArgK, OpMode: IABC /* */, Name: "LessThanOrEqual      "}, // if ((RK(B) <= RK(C)) ~= A) then pc++
	{TestFlag: 1, SetAFlag: 0, ArgBMode: OpArgN, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "Test    "},              // if not (R(A) <=> C) then pc++
	{TestFlag: 1, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "TestSet "},              // if (R(B) <=> C) then R(A) := R(B) else pc++
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "CALL    "},              // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "TAILCALL"},              // return R(A)(R(A+1), ... ,R(A+B-1))
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "Return  "},              // return R(A), ... ,R(A+B-2)
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IAsBx /**/, Name: "ForLoop "},              // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IAsBx /**/, Name: "ForPrep "},              // R(A)-=R(A+2); pc+=sBx
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgN, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "TFORCALL"},              // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgR, ArgCMode: OpArgN, OpMode: IAsBx /**/, Name: "TFORLOOP"},              // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IABC /* */, Name: "SETLIST "},              // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABx /* */, Name: "CLOSURE "},              // R(A) := closure(KPROTO[Bx])
	{TestFlag: 0, SetAFlag: 1, ArgBMode: OpArgU, ArgCMode: OpArgN, OpMode: IABC /* */, Name: "VARARG  "},              // R(A), R(A+1), ..., R(A+B-2) = vararg
	{TestFlag: 0, SetAFlag: 0, ArgBMode: OpArgU, ArgCMode: OpArgU, OpMode: IAx /*  */, Name: "EXTRAARG"},              // extra (larger) argument for previous opcode
}

/*
 31       22       13       5    0
  +-------+^------+-^-----+-^-----
  |b=9bits |c=9bits |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    bx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |   sbx=18bits    |a=8bits|op=6|
  +-------+^------+-^-----+-^-----
  |    ax=26bits            |op=6|
  +-------+^------+-^-----+-^-----
 31      23      15       7      0
*/
type Instruction uint32

func (ins Instruction) Opcode() *OpCode {
	return OpCodes[ins&0x3f]
}

func (ins Instruction) OpName() string {
	return ins.Opcode().Name
}

func (ins Instruction) OpMode() OpMode {
	return ins.Opcode().OpMode
}

func (ins Instruction) ArgBMode() OpArgMask {
	return ins.Opcode().ArgBMode
}

func (ins Instruction) ArgCMode() OpArgMask {
	return ins.Opcode().ArgCMode
}

func (ins Instruction) ABC() (a, b, c int) {
	a = int(ins >> 6 & 0xFF)
	c = int(ins >> 14 & 0x1FF)
	b = int(ins >> 23 & 0x1FF)
	return
}

func (ins Instruction) ABx() (a, bx int) {
	a = int(ins >> 6 & 0xFF)
	bx = int(ins >> 14)
	return
}

func (ins Instruction) AsBx() (a, sbx int) {
	a, bx := ins.ABx()
	return a, bx - MaxArgsBx
}

func (ins Instruction) Ax() int {
	return int(ins >> 6)
}
