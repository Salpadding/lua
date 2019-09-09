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

func Test2(t *testing.T) {
	f, err := os.Open("testdata/hello_world.o")
	if err != nil {
		t.Error(err)
	}
	_, err = (&ByteCodeReader{
		Reader: f,
	}).Load()
	if err != nil{
		t.Error(err)
	}
}
