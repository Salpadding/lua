package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/Salpadding/lua/types/code"
	"github.com/Salpadding/lua/types/value"

	"github.com/Salpadding/lua/types/tag"
)

type Prototype struct {
	Source          string // debug
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []code.Instruction
	Constants       []value.Value
	UpValues        []UpValue
	Prototypes      []*Prototype
	LineInfo        []uint32         // debug
	LocalVariables  []*LocalVariable // debug
	UpvalueNames    []string         // debug
}

type Chunk struct {
}

type ByteCodeReader struct {
	io.Reader
}

func (b *ByteCodeReader) ReadBytes(n int) ([]byte, error) {
	res := make([]byte, n)
	read, err := b.Reader.Read(res)
	if err != nil {
		return nil, err
	}
	if read < n {
		return nil, errors.New("unexpected eof")
	}
	return res, nil
}

func (b *ByteCodeReader) ReadInt() (int64, error) {
	i, err := b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

func (b *ByteCodeReader) ReadUint32() (uint32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(data), nil
}

func (b *ByteCodeReader) ReadUint64() (uint64, error) {
	data, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(data), nil
}

func (b *ByteCodeReader) ReadByte() (byte, error) {
	data, err := b.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (b *ByteCodeReader) ReadFloat() (float64, error) {
	u, err := b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u), nil
}

func (b *ByteCodeReader) ReadString() (string, error) {
	size, err := b.ReadByte()
	if err != nil {
		return "", err
	}
	if size == 0 {
		return "", nil
	}
	if size != 0xff {
		data, err := b.ReadBytes(int(size - 1))
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	longSize, err := b.ReadUint64()
	if longSize > math.MaxInt64 || longSize == 0 {
		return "", errors.New("the string is too large or two small")
	}
	if err != nil {
		return "", err
	}
	str, err := b.ReadBytes(int(longSize - 1))
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func (b *ByteCodeReader) Load() (*Prototype, error) {
	if err := b.checkHeader(); err != nil {
		return nil, err
	}
	if _, err := b.ReadByte(); err != nil {
		return nil, err
	} // size_upvalues

	return b.ReadPrototype()
}

func (b *ByteCodeReader) ReadPrototype() (*Prototype, error) {
	res := &Prototype{}
	var err error
	if res.Source, err = b.ReadString(); err != nil {
		return nil, err
	}
	if res.LineDefined, err = b.ReadUint32(); err != nil {
		return nil, err
	}
	if res.LastLineDefined, err = b.ReadUint32(); err != nil {
		return nil, err
	}
	if res.NumParams, err = b.ReadByte(); err != nil {
		return nil, err
	}
	if res.IsVararg, err = b.ReadByte(); err != nil {
		return nil, err
	}
	if res.MaxStackSize, err = b.ReadByte(); err != nil {
		return nil, err
	}
	if res.Code, err = b.readCode(); err != nil {
		return nil, err
	}
	if res.Constants, err = b.readConstants(); err != nil {
		return nil, err
	}
	if res.UpValues, err = b.readUpValues(); err != nil {
		return nil, err
	}
	if res.Prototypes, err = b.readPrototypes(); err != nil {
		return nil, err
	}
	if res.LineInfo, err = b.readLineInfo(); err != nil {
		return nil, err
	}
	if res.LocalVariables, err = b.readLocalVariables(); err != nil {
		return nil, err
	}
	if res.UpvalueNames, err = b.readUpValueNames(); err != nil {
		return nil, err
	}
	return res, nil
}

func (b *ByteCodeReader) readCode() ([]code.Instruction, error) {
	rawCodes, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	codes := make([]code.Instruction, rawCodes)
	for i := range codes {
		c, err := b.ReadUint32()
		if err != nil {
			return nil, err
		}
		codes[i] = code.Instruction(c)
	}
	return codes, nil
}

func (b *ByteCodeReader) readConstants() ([]value.Value, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	constants := make([]value.Value, size)
	for i := range constants {
		constants[i], err = b.readConstant()
		if err != nil {
			return nil, err
		}
	}
	return constants, nil
}

func (b *ByteCodeReader) readConstant() (value.Value, error) {
	t, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	switch t {
	case tag.Nil:
		return value.GetNil(), nil
	case tag.Boolean:
		n, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		return value.Boolean(n != 0), nil
	case tag.Number:
		i, err := b.ReadFloat()
		if err != nil {
			return nil, err
		}
		return value.Float(i), nil
	case tag.Integer:
		i, err := b.ReadInt()
		if err != nil {
			return nil, err
		}
		return value.Integer(i), nil
	case tag.ShortString, tag.LongString:
		str, err := b.ReadString()
		if err != nil {
			return nil, err
		}
		return value.String(str), nil
	default:
		return nil, errors.New("unsupported constant type")
	}
}

func (b *ByteCodeReader) checkHeader() error {
	if sig, err := b.ReadBytes(4); err != nil || !bytes.Equal(sig, []byte(LuaSignature)) {
		return errors.New("signature check fail")
	}
	if v, err := b.ReadByte(); err != nil || v != LuaVersion {
		return errors.New("version not match")
	}
	if f, err := b.ReadByte(); err != nil || f != LuaFormat {
		return errors.New("format mismatch")
	}
	if data, err := b.ReadBytes(6); err != nil || !bytes.Equal(data, []byte(LuaData)) {
		return errors.New("corrupted")
	}
	if s, err := b.ReadByte(); err != nil || s != CIntSize {
		return errors.New("int size mismatch")
	}
	if s, err := b.ReadByte(); err != nil || s != CSizeTSize {
		return errors.New("size_t size mismatch")
	}
	if s, err := b.ReadByte(); err != nil || s != InstructionSize {
		return errors.New("instruction size mismatch")
	}
	if s, err := b.ReadByte(); err != nil || s != LuaIntegerSize {
		return errors.New("instruction size mismatch")
	}
	if s, err := b.ReadByte(); err != nil || s != LuaNumberSize {
		return errors.New("instruction size mismatch")
	}
	if i, err := b.ReadInt(); err != nil || i != LuaCInt {
		return errors.New("endianness mismatch")
	}
	if f, err := b.ReadFloat(); err != nil || f != LuaCNumber {
		return errors.New("float format mismatch")
	}
	return nil
}

func (b *ByteCodeReader) readPrototypes() ([]*Prototype, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	prototypes := make([]*Prototype, size)
	for i := range prototypes {
		prototypes[i], err = b.ReadPrototype()
		if err != nil {
			return nil, err
		}
	}
	return prototypes, nil
}

func (b *ByteCodeReader) readUpValues() ([]UpValue, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	upValues := make([]UpValue, size)
	for i := range upValues {
		val, err := b.ReadBytes(2)
		if err != nil {
			return nil, err
		}
		copy(upValues[i][:], val[:])
	}
	return upValues, nil
}

func (b *ByteCodeReader) readLineInfo() ([]uint32, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	lineInfo := make([]uint32, size)
	for i := range lineInfo {
		lineInfo[i], err = b.ReadUint32()
		if err != nil {
			return nil, err
		}
	}
	return lineInfo, nil
}

func (b *ByteCodeReader) readLocalVariables() ([]*LocalVariable, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	localVariables := make([]*LocalVariable, size)
	for i := range localVariables {
		localVariables[i] = &LocalVariable{}
		if localVariables[i].Name, err = b.ReadString(); err != nil {
			return nil, err
		}
		if localVariables[i].StartPC, err = b.ReadUint32(); err != nil {
			return nil, err
		}
		if localVariables[i].EndPC, err = b.ReadUint32(); err != nil {
			return nil, err
		}
	}
	return localVariables, nil
}

func (b *ByteCodeReader) readUpValueNames() ([]string, error) {
	size, err := b.ReadUint32()
	if err != nil {
		return nil, err
	}
	names := make([]string, size)
	for i := range names {
		names[i], err = b.ReadString()
		if err != nil {
			return nil, err
		}
	}
	return names, nil
}

// UpValue = instack + idx
type UpValue [2]byte

type LocalVariable struct {
	Name    string
	StartPC uint32
	EndPC   uint32
}

func ReadPrototype(rd io.Reader) (*Prototype, error){
	return (&ByteCodeReader{
		Reader: rd,
	}).Load()
}