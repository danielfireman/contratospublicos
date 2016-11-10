package supplier

import (
	"context"
	"fmt"
	"sync"
)

// NewDataFetcher creates a brand new supplier data fetcher.
func NewDataFetcher(dbURI string, cities map[string]string) (*DataFecher, error) {
	db, err := dialDB(dbURI, cities)
	if err != nil {
		return nil, err
	}
	return &DataFecher{db}, nil
}

// dataFetcher is responsible for fetching supplier-related information.
type DataFecher struct {
	db *db
}

// NotFoundErr is returned in all supplier-related ops when the supplier is not found.
var NotFoundErr = fmt.Errorf("nf")

// Summary fetches the supplier's summary informatio.
func (f *DataFecher) Summary(ctx context.Context, id, legislature string) (*Fornecedor, error) {
	// Making sure calls to receitaws are cancelled.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	supplier := &Fornecedor{}

	// Asynchronously calling both DB and ReceitaWS and waits for either all of them successfully return or
	// the first error.
	var wg sync.WaitGroup
	errChan := make(chan error)
	wg.Add(2)
	go func() {
		defer wg.Done()
		f.db.FetchSummaryData(errChan, id, legislature, supplier)
	}()
	go func() {
		defer wg.Done()
		FetchReceitaWSData(ctx, errChan, id, supplier)
	}()
	go func() {
		wg.Wait()
		close(errChan)
	}()
	if err := <-errChan; err != nil {
		return nil, err
	}
	return supplier, nil
}
