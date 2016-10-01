package fetcher

import (
	"context"
)

type Fetcher interface {
	Fetch(ctx context.Context, id string) (interface{}, error)
}
