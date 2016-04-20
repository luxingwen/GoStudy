package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("err :", err)
		}
		defer file.Close()
		fmt.Fprintf(w, "%s", "ok")
		f, err := os.OpenFile("lxw/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("err2: ", err)
		}
		defer f.Close()
		io.Copy(f, file)
	} else if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, 0)
	}
}

func main() {
	http.HandleFunc("/", upload)
	http.Handle("/lxw/", http.StripPrefix("/lxw/", http.FileServer(http.Dir("lxw/"))))
	http.ListenAndServe(":9090", nil)
}
