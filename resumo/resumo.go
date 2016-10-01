package resumo

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/model"
)

const (
	DB                  = "heroku_q6gnv76m"
	dadosMunicipiosPath = "dados_municipios.csv"
)

var municipios = make(map[string]string)

func init() {
	f, err := os.Open(dadosMunicipiosPath)
	if err != nil {
		log.Fatal("Erro ao carregar arquivo de municípios: %q", err)
	}
	r := csv.NewReader(bufio.NewReader(f))
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		municipios[l[0]] = l[1]
	}
	fmt.Println("Municípios carregados com sucesso.")
}

type resumo struct {
	session     *mgo.Session
	legislatura string
}

func FetcherFromMongoDB(session *mgo.Session, legislatura string) fetcher.Fetcher {
	return &resumo{session, legislatura}
}

func (r *resumo) Fetch(ctx context.Context, id string) (interface{}, error) {
	reqSession := r.session.Copy()
	defer reqSession.Close()

	// TODO(danielfireman): Usar um sync.Pool
	ret := &model.ResumoContratosFornecedor{}
	c := reqSession.DB(DB).C(r.legislatura)
	if err := c.Find(bson.M{"id": id}).One(ret); err != nil {
		return nil, err
	}

	// Adicionando nomes aos municipios.
	for _, m := range ret.Municipios {
		nome, ok := municipios[m.Cod]
		if ok {
			m.Nome = nome
		}
	}
	return ret, nil
}
