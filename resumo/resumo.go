package resumo

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"
)

const (
	db                  = "heroku_q6gnv76m"
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
	store       store.Store
	legislatura string
}

func Fetcher(s store.Store, legislatura string) fetcher.Fetcher {
	return &resumo{s, legislatura}
}

func (r *resumo) Fetch(ctx context.Context, id string) (interface{}, error) {
	// TODO(danielfireman): Usar um sync.Pool
	ret := &model.ResumoContratosFornecedor{}
	if err := r.store.FindByID(db, r.legislatura, id, ret); err != nil {
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
