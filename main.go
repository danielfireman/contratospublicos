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
	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent"

	"gopkg.in/mgo.v2"
)

const (
	DB = "heroku_q6gnv76m"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}

	dbURI := os.Getenv("MONGODB_URI")
	if dbURI == "" {
		log.Fatal("Variável de ambiente $MONGHQ_URL obrigatória.")
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

	mainSession, err := mgo.Dial(dbURI)
	if err != nil {
		log.Fatalf("Erro carregando mapa de municípios: %q", err)
	}

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

		// Usando nosso BD como fonte autoritativa para buscas. Se não existe lá, nós
		// não conhecemos. Por isso, essa chamada é síncrona.
		fSeg := newrelic.StartSegment(txn, "fornecedores_collection_query")
		f, err := fornecedor.FetcherFromMongoDB(mainSession).Fetch(ctx, id)
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
		fSeg.End()

		resultado := &model.Fornecedor{
			ID:          f.(*model.DadosFornecedor).ID,
			Nome:        f.(*model.DadosFornecedor).Nome,
			Legislatura: legislatura,
		}

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
			dr := v.(*receitaws.DadosReceitaWS)
			res.AtividadePrincipal = dr.AtividadePrincipal
			res.DataSituacao = dr.DataSituacao
			res.Tipo = dr.Tipo
			res.AtividadesSecundarias = dr.AtividadesSecundarias
			res.Situacao = dr.Situacao
			res.NomeReceita = dr.Nome
			res.Telefone = dr.Telefone
			res.Cnpj = dr.Cnpj
			res.Municipio = dr.Municipio
			res.UF = dr.UF
			res.DataAbertura = dr.DataAbertura
			res.NaturezaJuridica = dr.NaturezaJuridica
			res.NomeFantasia = dr.NomeFantasia
			res.UltimaAtualizacaoReceita = dr.UltimaAtualizacao
			res.Bairro = dr.Bairro
			res.Logradouro = dr.Logradouro
			res.Numero = dr.CEP
			res.CEP = dr.CEP
		}(resultado)
		wg.Add(1)
		go func(res *model.Fornecedor) {
			defer wg.Done()
			defer newrelic.StartSegment(txn, legislatura+"_collection_query").End()
			r, err := resumo.FetcherFromMongoDB(mainSession, legislatura).Fetch(ctx, id)
			if err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
				return
			}
			rcf := r.(*model.ResumoContratosFornecedor)
			res.ValorContratos = rcf.ValorContratos
			res.NumContratos = rcf.NumContratos
			res.Municipios = rcf.Municipios
			res.Partidos = rcf.Partidos
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
