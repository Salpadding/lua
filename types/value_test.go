package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuaTable(t *testing.T){

	// float as integer index
	tb := NewTable()
	assert.NoError(t, tb.Set(Integer(1), Integer(1)))
	v, err := tb.Get(Float(1.0))
	assert.NoError(t, err)
	assert.Equal(t, Integer(1), v)


	assert.NoError(t, tb.Set(Integer(2), Integer(2)))
	assert.NoError(t, tb.Set(Float(3), Integer(3)))
	assert.Equal(t, 3, tb.Len())

	// float as float index
	assert.NoError(t, tb.Set(Float(5.2), Integer(3)))
	assert.Equal(t, 3, tb.Len())

	// set tail as hole shrink
	assert.NoError(t, tb.Set(Float(3), GetNil()))
	assert.Equal(t, 2, tb.Len())

	fmt.Println(tb)

	// expand
	assert.NoError(t, tb.Set(Integer(4), Integer(4)))
	assert.Equal(t, 2, tb.Len())

	assert.NoError(t, tb.Set(Integer(3), Integer(3)))
	assert.Equal(t, 4, tb.Len())

	// string key
	assert.NoError(t, tb.Set(String("4"), Integer(444)))
	v, err = tb.Get(String("4"))
	assert.NoError(t, err)
	assert.Equal(t, Integer(444), v)
}
