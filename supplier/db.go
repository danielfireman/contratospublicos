package supplier

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Dial establishes a database session.
// TODO(danielfireman): Set operation timeouts.
func dialDB(uri string, cities map[string]string) (*db, error) {
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
	return &db{
		session: s,
		name:    info.Database,
		cities:  cities,
	}, nil
}

type db struct {
	session *mgo.Session
	name    string
	cities  map[string]string
}

// FetchSummaryData populates the passed-in supplier info struct with information stored in  the database.
func (db *db) FetchSummaryData(errChan chan error, id, legislature string, supplier *Fornecedor) {
	cs, err := db.findSummary(id, legislature)
	if err != nil {
		// Consolidating not found errors.
		if err == mgo.ErrNotFound {
			errChan <- NotFoundErr
		} else {
			errChan <- err
		}
		return
	}
	// TODO(danielfireman): Fix this Enlgish PT hoge-podge in an next commit.
	supplier.ResumoContratos = &ResumoContratosFornecedor{
		ValorContratos:cs.AmountContracts,
		NumContratos:cs.NumContracts,
	}
	for _, c := range cs.Cities {
		supplier.ResumoContratos.Municipios = append(supplier.ResumoContratos.Municipios, &Municipio{
			Cod:          c.ID,
			Nome:         db.cities[c.ID],
			SiglaPartido: c.PartyInitials,
			ResumoContratos: ResumoContratos{
				Valor:      c.AmountCountracts,
				Quantidade: c.NumContracts,
			},
		})
	}
	for _, p := range cs.Parties {
		supplier.ResumoContratos.Partidos = append(supplier.ResumoContratos.Partidos, &Partido{
			Sigla: p.Initials,
			ResumoContratos: ResumoContratos{
				Valor:      p.AmountCountracts,
				Quantidade: p.NumContracts,
			},
		})
	}
}

// FindSupplier queries the database by the given supplier id.
func (db *db) findSummary(id, legislature string) (*ContractsSummary, error) {
	session := db.session.Copy()
	defer session.Close()

	dbS := &ContractsSummary{}
	err := session.DB(db.name).C(legislature).Find(bson.M{"id": id}).One(dbS)
	if err != nil {
		return nil, err
	}
	return dbS, nil
}

// NotFound returns true if the error means that the entity was not found
// at the database. Returns false otherwise.
func NotFound(err error) bool {
	return err == mgo.ErrNotFound
}

// ## Database Model ##
// All model structs are exported due to mgo serialization. Please don't use them outside this package.

// CityDataModel holds data about the relationship between a supplier and a certain city.
type City struct {
	ID               string  `bson:"cod,omitempty"`
	Name             string  `bson:"nome,omitempty"`
	NumContracts     int32   `bson:"qtd_contratos,omitempty"`
	AmountCountracts float64 `bson:"valor_contratos,omitempty"`
	PartyInitials    string  `bson:"sigla,omitempty"`
}

// PartyDataModel holds data about the relationship between a supplier and a certain city.
type Party struct {
	Initials         string  `bson:"sigla,omitempty"`
	NumContracts     int32  `bson:"qtd_contratos,omitempty"`
	AmountCountracts float64 `bson:"valor_contratos,omitempty"`
}

// SupplierContractsSummary holds summary information about the relationship between a certain state.
type ContractsSummary struct {
	ID              string   `bson:"id,omitempty"`
	AmountContracts float64  `bson:"valor_contratos,omitempty"`
	NumContracts    int32    `bson:"num_contratos,omitempty"`
	Cities          []*City  `bson:"municipios,omitempty"`
	Parties         []*Party `bson:"partidos,omitempty"`
}
