package chunk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type Prototype struct {
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

func (b *ByteCodeReader) ReadPrototype() (*Prototype, error) {
	return nil, nil
}

func (b *ByteCodeReader) checkHeader() error {
	if sig, err := b.ReadBytes(4); err != nil || ! bytes.Equal(sig, []byte(LuaSignature)) {
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
