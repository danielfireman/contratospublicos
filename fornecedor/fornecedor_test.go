package fornecedor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test(t *testing.T) {
	id := "id"
	want := &FornecedorDB{ID: id, Nome: "foo"}

	s := &mockStore{}
	s.On("FindByID", db, table, id, mock.Anything).Run(func(args mock.Arguments) {
		args[3].(*FornecedorDB).ID = want.ID
		args[3].(*FornecedorDB).Nome = want.Nome
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
