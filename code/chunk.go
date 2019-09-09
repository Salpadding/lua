package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
	f, err := os.Open("chunk.go")
	if err != nil{
		panic(err)
	}
	rd := bufio.NewReader(f)
	for{
		l, _, err := rd.ReadLine()
		if err != nil{
			break
		}
		fmt.Println(string(l))
	}
}