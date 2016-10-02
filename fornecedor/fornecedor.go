package fornecedor

import (
	"context"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"

	"gopkg.in/mgo.v2/bson"
)

const (
	db    = "heroku_q6gnv76m"
	table = "fornecedores"
)

type FornecedorDB struct {
	BSONID bson.ObjectId `bson:"_id,omitempty"`
	ID     string        `bson:"id,omitempty"`
	Nome   string        `bson:"nome,omitempty"`
}

type fornecedorFetcher struct {
	store store.Store
}

func Fetcher(s store.Store) fetcher.Fetcher {
	return &fornecedorFetcher{s}
}

func (f *fornecedorFetcher) Fetch(ctx context.Context, id string) (interface{}, error) {
	// Usar sync.Pool.
	ret := &FornecedorDB{}
	err := f.store.FindByID(db, table, id, ret)
	return ret, err
}

func AtualizaFornecedor(f *model.Fornecedor, i interface{}) {
	fDB := i.(*FornecedorDB)
	f.ID = fDB.ID
	f.Nome = fDB.Nome
}
