package main

import (
	"enigma-protocol-go/pkg/api"
	"fmt"
	"net/http"
)

func main() {
	router, err := api.NewRouter(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
