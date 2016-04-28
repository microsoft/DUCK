package main

import (
	"net/http"
	"fmt"
	"flag"
)

func main() {

	var webDir string;

	flag.StringVar(&webDir, "webdir", "frontend", "The root directory for serving web content")
	flag.Parse()

	fmt.Println("Web root: " + webDir)

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.ListenAndServe(":3000", nil)
}
