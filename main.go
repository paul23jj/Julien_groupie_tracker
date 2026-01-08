package main

import (
	"Steam-API/routeur"
	"fmt"
	"net/http"
)

func main() {
	r := routeur.New()

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
