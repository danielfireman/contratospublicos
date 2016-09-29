package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT não encontrada.")
	}

	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprintf(w, "Hellow World.")
	})
	log.Println("Serviço inicializado na porta ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
