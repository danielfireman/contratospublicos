package fornecedor

import (
	"context"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"

	"gopkg.in/mgo.v2/bson"
	"reflect"
	"sync"
)

const (
	table = "fornecedores"
)

type FornecedorDB struct {
	BSONID bson.ObjectId `bson:"_id,omitempty"`
	ID     string        `bson:"id,omitempty"`
	Nome   string        `bson:"nome,omitempty"`
}

type coletorBDPrincipal struct {
	store store.Store
	pool  sync.Pool
}

func ColetorBD(s store.Store) ColetorDadosFornecedor {
	return &coletorBDPrincipal{
		store: s,
		pool: sync.Pool{
			New: func() interface{} {
				return &FornecedorDB{}
			},
		},
	}
}

func (f *coletorBDPrincipal) ColetaDados(ctx context.Context, fornecedor *model.Fornecedor) error {
	fDB := f.pool.Get().(*FornecedorDB)
	defer f.pool.Put(fDB)
	defer func(fDB *FornecedorDB) {
		p := reflect.ValueOf(fDB).Elem()
		p.Set(reflect.Zero(p.Type()))
	}(fDB)

	if err := f.store.FindByID(table, fornecedor.ID, fDB); err != nil {
		return err
	}
	fornecedor.ID = fDB.ID
	fornecedor.Nome = fDB.Nome
	return nil
}
