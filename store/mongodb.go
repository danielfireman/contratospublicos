package store

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBStore struct {
	session *mgo.Session
	db      string
}

func (s *mongoDBStore) FindByID(c string, id string, ret interface{}) error {
	session := s.session.Copy()
	defer session.Close()
	err := session.DB(s.db).C(c).Find(bson.M{"id": id}).One(ret)
	if err == mgo.ErrNotFound {
		return NaoEncontradoErr(err)
	}
	return err
}

func MongoDB(uri string) (Store, error) {
	if uri == "" {
		return nil, fmt.Errorf("MongoDB URI inválida.")
	}
	info, err := mgo.ParseURL(uri)
	if err != nil {
		return nil, fmt.Errorf("Erro processando URI:%s err:%q\n", uri, err)
	}
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	s.SetMode(mgo.Monotonic, true)
	return &mongoDBStore{s, info.Database}, nil
}
