package vm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T){
	f, err := os.Open("testdata/test3.o")
	assert.NoError(t, err)
	var vm LuaVM
	assert.NoError(t, vm.Load(f))
	assert.NoError(t, vm.Execute())
}