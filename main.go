package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"bufio"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/danielfireman/contratospublicos/supplier"
	"github.com/julienschmidt/httprouter"
)

var fornecedorTmpl = template.Must(template.New("fornecedor").ParseFiles("./fornecedor.tmpl.html"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public/"))
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.GET("/repo", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.Redirect(w, r, "http://github.com/danielfireman/contratospublicos", http.StatusMovedPermanently)
	})

	// Supplier-related Handlers.
	cities, err := loadCityNames()
	if err != nil {
		log.Fatal(err)
	}
	fetcher, err := supplier.NewDataFetcher(os.Getenv("MONGODB_URI"), cities)
	if err != nil {
		log.Fatal(err)
	}
	router.GET("/api/v1/fornecedor", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id, legislature := extractQueryParams(r)
		f, err := fetcher.Summary(context.Background(), id, legislature)
		if err != nil {
			if err == supplier.NotFoundErr {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Printf("Err id:%s leg:%s err:%q\n", id, legislature, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	router.GET("/fornecedor", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id, legislature := extractQueryParams(r)
		fVO, err := fetchSupplierVirtualObject(context.Background(), fetcher, id, legislature)
		if err != nil {
			if err == supplier.NotFoundErr {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Printf("Err id:%s leg:%s err:%q\n", id, legislature, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := fornecedorTmpl.Execute(w, fVO); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Service listening at port ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func extractQueryParams(r *http.Request) (string, string) {
	query := r.URL.Query()
	legislatura := query.Get("legislatura")
	if legislatura == "" {
		legislatura = "2012"
	}

	// NOTA: Using query parameters to allow users to copy and paste the permalink.
	// CNPJs have "/", which makes it a bad candidate for being part of the URL per se
	id := query.Get("cnpj")
	if id == "" {
		return "", legislatura
	}

	// Removing special chars. This allow users to copy and paste ids from other sources
	id = strings.Replace(id, ".", "", -1)
	id = strings.Replace(id, "-", "", -1)
	id = strings.Replace(id, "/", "", -1)
	return id, legislatura
}

// NOTE: To save DB space we do nome hold the city names in the DB. We keep them in memory instead.
func loadCityNames() (map[string]string, error) {
	f, err := os.Open("dados_municipios.csv")
	if err != nil {
		return nil, fmt.Errorf("Erro ao carregar arquivo de municípios: %q", err)
	}
	r := csv.NewReader(bufio.NewReader(f))
	cities := make(map[string]string)
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		cities[l[0]] = l[1]
	}
	fmt.Println("Cities map loaded successfully.")
	return cities, nil
}
