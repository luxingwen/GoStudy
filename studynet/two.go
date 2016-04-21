package main

import (
	"fmt"
	"net"
	"os"
)

//获取掩码信息
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Uage: %s dotted-ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	dotAddr := os.Args[1]
	addr := net.ParseIP(dotAddr)
	if addr == nil {
		fmt.Println("Invalid address")
		os.Exit(1)
	}
	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	ones, bits := mask.Size()
	fmt.Println("Address is ", addr.String())
	fmt.Println("Default mask length is ", bits)
	fmt.Println("Leading ones count is ", ones)
	fmt.Println("Mask is (hex)", mask.String())
	fmt.Println("Network is ", network.String())
	os.Exit(0)
}
