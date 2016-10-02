package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/fornecedor"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/receitaws"
	"github.com/danielfireman/contratospublicos/resumo"
	"github.com/danielfireman/contratospublicos/store"
	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent"
)

const (
	DB = "heroku_q6gnv76m"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}

	mongoDBStore, err := store.MongoDB(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal(err)
	}

	nrLicence := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if nrLicence == "" {
		log.Fatal("$NEW_RELIC_LICENSE_KEY must be set")
	}
	config := newrelic.NewConfig("ciframe-api", nrLicence)
	newRelicApp, err := newrelic.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Monitoramento NewRelic configurado com sucesso.")

	router := httprouter.New()
	router.GET("/api/v1/fornecedor/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		txn := newRelicApp.StartTransaction("fornecedor", w, r)
		defer txn.End()

		ctx := context.Background()
		id := p.ByName("id")

		legislatura := r.URL.Query().Get("legislatura")
		if legislatura == "" {
			legislatura = "2012"
		}

		resultado := &model.Fornecedor{
			ID:          id,
			Legislatura: legislatura,
		}

		// Usando nosso BD como fonte autoritativa para buscas. Se não existe lá, nós
		// não conhecemos. Por isso, essa chamada é síncrona.
		fSeg := newrelic.StartSegment(txn, "fornecedores_collection_query")
		f, err := fornecedor.Fetcher(mongoDBStore).Fetch(ctx, id)
		if err != nil {
			if fetcher.IsNotFound(err) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				log.Println("Err id:'%s' err:'%q'", id, err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			fSeg.End()
			return
		}
		fornecedor.AtualizaFornecedor(resultado, f)
		fSeg.End()

		// Pegando dados remotos de forma concorrente.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(res *model.Fornecedor) {
			defer wg.Done()
			defer newrelic.StartSegment(txn, "receitaws_query").End()
			v, err := receitaws.Fetch(ctx, id)
			if err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
				return
			}
			receitaws.AtualizaFornecedor(res, v)

		}(resultado)
		wg.Add(1)
		go func(res *model.Fornecedor) {
			defer wg.Done()
			defer newrelic.StartSegment(txn, legislatura+"_collection_query").End()
			r, err := resumo.Fetcher(mongoDBStore, legislatura).Fetch(ctx, id)
			if err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
				return
			}
			resumo.AtualizaFornecedor(res, r)
		}(resultado)
		wg.Wait()

		marshallSeg := newrelic.StartSegment(txn, "marshall_results")
		b, err := json.Marshal(resultado)
		marshallSeg.End()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, string(b))
	})
	log.Println("Serviço inicializado na porta ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
