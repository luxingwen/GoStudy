package main

import (
	"strings"
)
import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now().Unix()
	fmt.Println(strings.Join(os.Args[1:], " "))
	end := time.Now().Unix()
	fmt.Println("join run time: ", end-start)
	start = time.Now().Unix()
	var s, sep string
	for _, v := range os.Args[1:] {
		s += sep + v
		sep = " "
	}
	fmt.Println(s)
	end = time.Now().Unix()
	fmt.Println("range run time: ", end-start)
}
