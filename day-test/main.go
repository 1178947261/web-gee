package main

import "fmt"

type name struct {
	R map[string]string
}

func main() {

	n := &name{
		R: make(map[string]string),
	}

	e := &name{}
	n.R["aa"] = "aaaa"
	e.R["bbb"] = "11111"
	fmt.Println(n.R)
	fmt.Println(e.R)
}
