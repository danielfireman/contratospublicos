package fornecedor

import (
	"context"
	"log"
	"sync"

	"github.com/danielfireman/contratospublicos/model"
)

type Buscador struct {
	principal   ColetorDadosFornecedor
	secundarios []ColetorDadosFornecedor
}

func (b *Buscador) ColetaDados(id, legislatura string) (*model.Fornecedor, error) {
	ctx := context.Background()
	resultado := &model.Fornecedor{
		ID:          id,
		Legislatura: legislatura,
	}
	// Usando nosso BD como fonte autoritativa para buscas. Se não existe lá, nós
	// não conhecemos. Por isso, essa chamada é síncrona.
	if err := b.principal.ColetaDados(ctx, resultado); err != nil {
		return nil, err
	}
	// Busca dados dos coletores remotos de forma concorrente.
	wg := sync.WaitGroup{}
	for _, coletor := range b.secundarios {
		wg.Add(1)
		go func(coletor ColetorDadosFornecedor, res *model.Fornecedor) {
			defer wg.Done()
			if err := coletor.ColetaDados(ctx, res); err != nil {
				log.Println("Err id:'%s' err:'%q'", id, err)
			}
		}(coletor, resultado)
	}
	wg.Wait()
	return resultado, nil
}
