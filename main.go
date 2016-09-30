package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/julienschmidt/httprouter"
)

const (
	DB = "heroku_q6gnv76m"
)

const dadosMunicipiosPath = "dados_municipios.csv"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}

	dbURI := os.Getenv("MONGODB_URI")
	if dbURI == "" {
		log.Fatal("Variável de ambiente $MONGHQ_URL obrigatória.")
	}

	municipios, err := carregaMunicipios(dadosMunicipiosPath)
	if err != nil {
		log.Fatalf("Erro carregando mapa de municípios: %q", err)
	}
	fmt.Println("Municípios carregados com sucesso.")

	mainSession, err := mgo.Dial(dbURI)
	if err != nil {
		log.Fatalf("Erro carregando mapa de municípios: %q", err)
	}

	router := httprouter.New()
	router.GET("/api/v1/fornecedor/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session := mainSession.Copy()
		defer session.Close()

		fornecedor := &model.DadosFornecedor{}
		id := p.ByName("id")
		func() {
			c := session.DB(DB).C("fornecedores")
			if err = c.Find(bson.M{"id": id}).One(&fornecedor); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Err id:'%s' err:''", id, err)
				return
			}
		}()

		legislatura := r.URL.Query().Get("legislatura")
		if legislatura == "" {
			legislatura = "2012"
		}

		resumo := model.ResumoContratosFornecedor{}
		func() {
			c := session.DB(DB).C(legislatura)
			if err = c.Find(bson.M{"id": id}).One(&resumo); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Err id:'%s' err:''", id, err)
				return
			}
		}()

		// Adicionando nomes aos municipios.
		for _, m := range resumo.Municipios {
			nome, ok := municipios[m.Cod]
			if ok {
				m.Nome = nome
			}
		}
		b, err := json.Marshal(&model.Fornecedor{
			ID:             fornecedor.ID,
			Nome:           fornecedor.Nome,
			Legislatura:    legislatura,
			ValorContratos: resumo.ValorContratos,
			NumContratos:   resumo.NumContratos,
			Municipios:     resumo.Municipios,
			Partidos:       resumo.Partidos,
		})
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

func carregaMunicipios(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	returned := make(map[string]string)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		returned[l[0]] = l[1]
	}
	return returned, nil
}
