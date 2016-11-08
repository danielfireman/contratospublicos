package fornecedor

import (
	"context"
	"strconv"
	"testing"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestColetorDB_ColetaDados(t *testing.T) {
	var dados []*model.Fornecedor
	for i := 1; i < 10; i++ {
		dados = append(dados, &model.Fornecedor{ID: strconv.Itoa(i)})
	}
	// Executando chamadas em paralelo para verificar se os pools de objetos estão funcionando corretamente.
	for _, want := range dados {
		s := &mockStore{}
		s.On("FindByID", table, want.ID, mock.Anything).Run(func(args mock.Arguments) {
			args[2].(*FornecedorDB).ID = want.ID
			args[2].(*FornecedorDB).Nome = want.Nome
		}).Return(nil)

		coletor := ColetorBD(s)
		got := &model.Fornecedor{ID: want.ID}
		err := coletor.ColetaDados(context.TODO(), got)
		assert.Nil(t, err)
		assert.Exactly(t, want, got)
		s.AssertExpectations(t)
	}
}

type mockStore struct {
	mock.Mock
}

func (m *mockStore) FindByID(c string, id string, ret interface{}) error {
	m.Called(c, id, ret)
	return nil
}
