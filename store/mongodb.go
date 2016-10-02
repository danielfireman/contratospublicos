package store

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/danielfireman/contratospublicos/fetcher"
)

type mongoDBStore struct {
	session *mgo.Session
}

func (s *mongoDBStore) FindByID(db, c string, id string, ret interface{}) error {
	session := s.session.Copy()
	defer session.Close()
	err := session.DB(db).C(c).Find(bson.M{"id": id}).One(ret)
	if err == mgo.ErrNotFound {
		return fetcher.NotFound(err)
	}
	return err
}

func MongoDB(uri string) (Store, error) {
	if uri == "" {
		return nil, fmt.Errorf("MongoDB URI inválida.")
	}
	s, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &mongoDBStore{s}, nil
}
