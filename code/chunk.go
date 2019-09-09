package chunk

import (
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

func (b *ByteCodeReader) ReadString() (string, error) {
	size, err := b.ReadByte()
	if err != nil {
		return "", err
	}
	if size == 0 {
		return "", nil
	}
	if size != 0xff {
		bytes, err := b.ReadBytes(int(size - 1))
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	longSize, err := b.ReadUint64()
	if longSize > math.MaxInt64 || longSize == 0{
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
