package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Hello, World!")
		fmt.Println("Number of bytes written:", n)
		if err != nil {
			fmt.Println("Error:", err)
		}
	})

	_ = http.ListenAndServe(":8080", nil)
}
