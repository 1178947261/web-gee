package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexPage)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		return
	}
}

//主页
func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}
