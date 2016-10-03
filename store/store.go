package store

type Store interface {
	FindByID(c string, id string, ret interface{}) error
}

type NaoEncontradoErr struct {
	error
}
