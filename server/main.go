package main

import (
	"os"
	"log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT não encontrada.")
	}
	log.Println("Porta utilizada: ", port)

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintf(w, "Hellow World.")
	})
}
