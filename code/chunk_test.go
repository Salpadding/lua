package chunk

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	f, err := os.Open("testdata/hello_world.o")
	if err != nil {
		t.Error(err)
	}
	err = (&ByteCodeReader{
		Reader: f,
	}).checkHeader()
	if err != nil{
		t.Error(err)
	}
}
