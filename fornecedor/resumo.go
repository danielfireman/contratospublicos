package fornecedor

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"

	"gopkg.in/mgo.v2/bson"
)

const (
	dadosMunicipiosPath = "dados_municipios.csv"
)

type ResumoContratosDB struct {
	BSONID         bson.ObjectId      `bson:"_id,omitempty"`
	ID             string             `bson:"id,omitempty"`
	ValorContratos float64            `bson:"valor_contratos,omitempty"`
	NumContratos   int64              `bson:"num_contratos,omitempty"`
	Municipios     []*model.Municipio `bson:"municipios,omitempty"`
	Partidos       []*model.Partido   `bson:"partidos,omitempty"`
}

type coletorResumoContratos struct {
	store      store.Store
	pool       sync.Pool
	municipios map[string]string
}

func ColetorResumoContratos(s store.Store) ColetorDadosFornecedor {
	f, err := os.Open(dadosMunicipiosPath)
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo de municípios: %q", err)
	}
	r := csv.NewReader(bufio.NewReader(f))
	municipios := make(map[string]string)
	for {
		l, err := r.Read()
		if err == io.EOF {
			break
		}
		municipios[l[0]] = l[1]
	}
	fmt.Println("Municípios carregados com sucesso.")
	return &coletorResumoContratos{
		store: s,
		pool: sync.Pool{
			New: func() interface{} {
				return &ResumoContratosDB{}
			},
		},
		municipios: municipios,
	}
}

func (r *coletorResumoContratos) ColetaDados(ctx context.Context, fornecedor *model.Fornecedor) error {
	rDB := r.pool.Get().(*ResumoContratosDB)
	defer r.pool.Put(rDB)
	defer func(rDB *ResumoContratosDB) {
		p := reflect.ValueOf(rDB).Elem()
		p.Set(reflect.Zero(p.Type()))
	}(rDB)
	if err := r.store.FindByID(fornecedor.Legislatura, fornecedor.ID, rDB); err != nil {
		return err
	}
	for _, m := range rDB.Municipios {
		nome, ok := r.municipios[m.Cod]
		if ok {
			m.Nome = nome
		}
	}
	fornecedor.ResumoContratos = &model.ResumoContratosFornecedor{
		ValorContratos: rDB.ValorContratos,
		NumContratos:   rDB.NumContratos,
		Municipios:     rDB.Municipios,
		Partidos:       rDB.Partidos,
	}
	return nil
}
