package main

import (
	"Steam-API/routeur"
	"fmt"
	"net/http"
	"time"
)

func main() {
	r := routeur.New()

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

func formatDate(v interface )