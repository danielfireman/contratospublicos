package fornecedor

import (
	"context"
	"testing"

	"fmt"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test(t *testing.T) {
	id := "id"
	want := &model.DadosFornecedor{ID: id, Nome: "foo"}

	s := &mockStore{}
	s.On("FindByID", db, table, id, mock.Anything).Run(func(args mock.Arguments) {
		fmt.Print("Chamou")
		args[3].(*model.DadosFornecedor).ID = want.ID
		args[3].(*model.DadosFornecedor).Nome = want.Nome
	}).Return(nil)

	f := fornecedorFetcher{s}
	got, err := f.Fetch(context.TODO(), id)
	assert.Nil(t, err)
	assert.Exactly(t, want, got)
	s.AssertExpectations(t)
}

type mockStore struct {
	mock.Mock
}

func (m *mockStore) FindByID(db, c string, id string, ret interface{}) error {
	m.Called(db, c, id, ret)
	return nil
}
