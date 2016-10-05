package fornecedor

import (
	"context"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"
)

type ColetorDadosFornecedor interface {
	ColetaDados(context.Context, *model.Fornecedor) error
}

func NaoEncontrado(err error) bool {
	_, ok := err.(store.NaoEncontradoErr)
	return ok
}
