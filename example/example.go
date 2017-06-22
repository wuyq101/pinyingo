package main

import (
	"fmt"

	"github.com/wuyq101/pinyingo"
)

func main() {
	py := pinyingo.NewPy(pinyingo.StyleNormal, pinyingo.UseSegment)
	s := "苹果"
	r := py.Convert(s)
	fmt.Println(r)
}
