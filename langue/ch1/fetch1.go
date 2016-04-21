package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: %v\n", err)
			os.Exit(1)
		}
		resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		//        dst,err:=os.Create("hahh.html")
		//        if err!=nil{
		//            fmt.Fprintf(os.Stderr,"can't creat htm   err:%v\n",err)
		//            os.Exit(1)
		//        }
		var buf [512]byte

		//       _,err= io.Copy(dst.W,b)
		//    if err!=nil{
		//        fmt.Fprintf("Fetch :io copy %s :%v\n",url,err)
		//        os.Exit(1)
		//    }
	}
}
