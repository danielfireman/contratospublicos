package main

import (
	"log"
	"os"

	"html/template"
	"net/http"

	"github.com/danielfireman/contratospublicos/fornecedor"
	"github.com/julienschmidt/httprouter"
)

var fornecedorTmpl = template.Must(template.New("fornecedor").ParseFiles("./fornecedor/fornecedor.html"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}

	// Configuração do roteador echo.
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public/"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/repo", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.Redirect(w, r, "http://github.com/danielfireman/contratospublicos", http.StatusMovedPermanently)
	})

	// Lida com fornecedores
	fTratadores, err := fornecedor.Tratadores(fornecedorTmpl)
	if err != nil {
		log.Fatal(err)
	}
	router.GET("/api/v1/fornecedor", fTratadores.TrataAPICall())
	router.GET("/fornecedor", fTratadores.TrataPaginaFornecedor())

	log.Println("Serviço inicializado na porta ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
