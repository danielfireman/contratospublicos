package fornecedor

import (
	"context"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielfireman/contratospublicos/fetcher"
	"github.com/danielfireman/contratospublicos/model"
)

const (
	DB    = "heroku_q6gnv76m"
	table = "fornecedores"
)

type fornecedor struct {
	session *mgo.Session
}

func FetcherFromMongoDB(session *mgo.Session) fetcher.Fetcher {
	return &fornecedor{session}
}

func (f *fornecedor) Fetch(ctx context.Context, id string) (interface{}, error) {
	reqSession := f.session.Copy()
	defer reqSession.Close()

	// TODO(danielfireman): Usar um sync.Pool
	ret := model.Fornecedor{}
	c := reqSession.DB(DB).C(table)
	if err := c.Find(bson.M{"id": id}).One(&ret); err != nil {
		if err == mgo.ErrNotFound { // Aqui consolidamos o tratamento de errors do tipo: não encontrado.
			return nil, fetcher.NotFound(err)
		}
		return nil, err
	}
	return ret, nil
}
