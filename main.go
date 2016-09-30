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
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/receitaws"
	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent"
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

	munTxn := newRelicApp.StartTransaction("load_cities", nil, nil)
	municipios, err := carregaMunicipios(dadosMunicipiosPath)
	if err != nil {
		log.Fatalf("Erro carregando mapa de municípios: %q", err)
	}
	munTxn .End()
	fmt.Println("Municípios carregados com sucesso.")

	mainSession, err := mgo.Dial(dbURI)
	if err != nil {
		log.Fatalf("Erro carregando mapa de municípios: %q", err)
	}

	router := httprouter.New()
	router.GET("/api/v1/fornecedor/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		txn := newRelicApp.StartTransaction("fornecedor", w, r)
		defer txn.End()

		session := mainSession.Copy()
		defer session.Close()

		id := p.ByName("id")

		legislatura := r.URL.Query().Get("legislatura")
		if legislatura == "" {
			legislatura = "2012"
		}

		fornecedor := &model.DadosFornecedor{}
		resumo := &model.ResumoContratosFornecedor{}
		dadosReceitaWs := &receitaws.DadosReceitaWS{}

		fSeg := newrelic.StartSegment(txn, "fornecedores_collection_query")
		c := session.DB(DB).C("fornecedores")
		if err = c.Find(bson.M{"id": id}).One(&fornecedor); err != nil {
			log.Println("Err id:'%s' err:'%q'", id, err)
			if err == mgo.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			fSeg.End()

			// Usando nosso BD como fonte autoritativa para buscas. Se não existe lá, nós
			// não conhecemos.
			return
		}
		fSeg.End()

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer newrelic.StartSegment(txn, "receitaws_query").End()
			if err := receitaws.GetData(id, dadosReceitaWs); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer newrelic.StartSegment(txn, legislatura + "_collection_query").End()
			c := session.DB(DB).C(legislatura)
			if err = c.Find(bson.M{"id": id}).One(resumo); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}()
		wg.Wait()

		// Adicionando nomes aos municipios.
		for _, m := range resumo.Municipios {
			nome, ok := municipios[m.Cod]
			if ok {
				m.Nome = nome
			}
		}

		resultado := &model.Fornecedor{
			ID:             fornecedor.ID,
			Nome:           fornecedor.Nome,
			Legislatura:    legislatura,
			ValorContratos: resumo.ValorContratos,
			NumContratos:   resumo.NumContratos,
			Municipios:     resumo.Municipios,
			Partidos:       resumo.Partidos,
		}

		if resumo != nil {
			resultado.ValorContratos = resumo.ValorContratos
			resultado.NumContratos = resumo.NumContratos
			resultado.Municipios = resumo.Municipios
			resultado.Partidos = resumo.Partidos
		}

		if dadosReceitaWs != nil {
			resultado.AtividadePrincipal = dadosReceitaWs.AtividadePrincipal
			resultado.DataSituacao = dadosReceitaWs.DataSituacao
			resultado.Tipo = dadosReceitaWs.Tipo
			resultado.AtividadesSecundarias = dadosReceitaWs.AtividadesSecundarias
			resultado.Situacao = dadosReceitaWs.Situacao
			resultado.NomeReceita = dadosReceitaWs.Nome
			resultado.Telefone = dadosReceitaWs.Telefone
			resultado.Cnpj = dadosReceitaWs.Cnpj
			resultado.Municipio = dadosReceitaWs.Municipio
			resultado.UF = dadosReceitaWs.UF
			resultado.DataAbertura = dadosReceitaWs.DataAbertura
			resultado.NaturezaJuridica = dadosReceitaWs.NaturezaJuridica
			resultado.NomeFantasia = dadosReceitaWs.NomeFantasia
			resultado.UltimaAtualizacaoReceita = dadosReceitaWs.UltimaAtualizacao
			resultado.Bairro = dadosReceitaWs.Bairro
			resultado.Logradouro = dadosReceitaWs.Logradouro
			resultado.Numero = dadosReceitaWs.CEP
			resultado.CEP = dadosReceitaWs.CEP
		}

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
	log.Fatal(http.ListenAndServe(":" + port, router))
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
