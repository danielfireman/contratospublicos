package fetcher

import (
	"context"
)

type notFound struct {
	error
}

func NotFound(err error) error {
	return notFound{err}
}

func IsNotFound(err error) bool {
	_, ok := err.(notFound)
	return ok
}

type Fetcher interface {
	Fetch(ctx context.Context, id string) (interface{}, error)
}
