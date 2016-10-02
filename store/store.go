package store

type Store interface {
	FindByID(db, c string, id string, ret interface{}) error
}
