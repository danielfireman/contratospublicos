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
	"github.com/danielfireman/contratospublicos/receitaws"
	"github.com/julienschmidt/httprouter"
	"sync"
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

		id := p.ByName("id")

		legislatura := r.URL.Query().Get("legislatura")
		if legislatura == "" {
			legislatura = "2012"
		}

		fornecedor := &model.DadosFornecedor{}
		resumo := &model.ResumoContratosFornecedor{}
		dadosReceitaWs := &receitaws.DadosReceitaWS{}

		var fornecedoresColErr error

		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := session.DB(DB).C("fornecedores")
			if fornecedoresColErr = c.Find(bson.M{"id": id}).One(&fornecedor); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := receitaws.GetData(id, dadosReceitaWs); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := session.DB(DB).C(legislatura)
			if err = c.Find(bson.M{"id": id}).One(resumo); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}()
		wg.Wait()

		if fornecedoresColErr != nil {
			if fornecedoresColErr == mgo.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

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

		b, err := json.Marshal(resultado)
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
