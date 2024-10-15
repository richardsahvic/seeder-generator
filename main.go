package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/richardsahvic/generate-query/logic"
)

func main() {
	http.HandleFunc("/generate/query", logic.GenerateQueryHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
