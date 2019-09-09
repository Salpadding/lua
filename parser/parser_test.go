package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func testParser(t *testing.T, fname string) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Error(err)
	}
	p, err := New(bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
	blk, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blk.String())
}

func TestParser1(t *testing.T) {
	testParser(t, "testdata/p1.lua")
}

func TestParser2(t *testing.T) {
	testParser(t, "testdata/p2.lua")
}

func TestParser3(t *testing.T) {
	testParser(t, "testdata/p3.lua")
}

func TestParser4(t *testing.T) {
	testParser(t, "testdata/p4.lua")
}

func TestParser5(t *testing.T) {
	testParser(t, "testdata/p5.lua")
}
