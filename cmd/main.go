package main

import (
	"fmt"
	"net/http"
	"rwa/internal/realworld"
)

// сюда код писать не надо

func main() {
	addr := ":8080"
	h := realworld.GetApp()
	fmt.Println("start server at", addr)
	http.ListenAndServe(addr, h)
}
