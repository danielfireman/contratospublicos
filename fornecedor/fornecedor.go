package fornecedor

import (
	"context"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"
)

const (
	db    = "heroku_q6gnv76m"
	table = "fornecedores"
)

type fornecedorFetcher struct {
	store store.Store
}

func Fetcher(s store.Store) fetcher.Fetcher {
	return &fornecedorFetcher{s}
}

func (f *fornecedorFetcher) Fetch(ctx context.Context, id string) (interface{}, error) {
	// Usar sync.Pool.
	ret := &model.DadosFornecedor{}
	err := f.store.FindByID(db, table, id, ret)
	return ret, err
}
