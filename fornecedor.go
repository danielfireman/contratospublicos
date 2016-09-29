package contratospublicos

type Fornecedor struct {
	ID                  string
	Nome                string
	ValorTotalContratos float64
	NumTotalContratos   int64
	Municipios          []Municipio
}

type ResumoContratos struct {
	Quantidade int64
	Valor      float64
}

type Municipio struct {
	Cod             string;
	Nome            string;
	ResumoContratos ResumoContratos;
	SiglaPartido    string;
}

type Partido struct {
	Sigla           string;
	ResumoContratos ResumoContratos;
}