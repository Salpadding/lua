package vm

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T){
	f, err := os.Open("testdata/test3.o")
	assert.NoError(t, err)
	var vm LuaVM
	assert.NoError(t, vm.Load(f))
	assert.Error(t, vm.Execute()) // assertion fail error, native function
}

func TestUpValue(t *testing.T){
	f, err := os.Open("testdata/test4.o")
	assert.NoError(t, err)
	var vm LuaVM
	counter := &gasCounter{}
	vm.hooks = []Hook{counter.hook}
	assert.NoError(t, vm.Load(f))
	assert.NoError(t, vm.Execute())
	fmt.Print(counter.gas)

}