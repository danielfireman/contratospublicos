package model

import "gopkg.in/mgo.v2/bson"

type Fornecedor struct {
	ID             string       `json:"id"`
	Nome           string       `json:"nome"`
	Legislatura    string       `json:"legislatura"`
	ValorContratos float64      `json:"valor_contratos"`
	NumContratos   int64        `json:"num_contratos"`
	Municipios     []*Municipio `json:"municipios"`
	Partidos       []*Partido   `json:"partidos"`
}

type DadosFornecedor struct {
	BSONID bson.ObjectId `bson:"_id,omitempty"`
	ID     string        `bson:"id,omitempty"`
	Nome   string        `bson:"nome,omitempty"`
}

type ResumoContratosFornecedor struct {
	BSONID         bson.ObjectId `bson:"_id,omitempty"`
	ID             string        `bson:"id,omitempty"`
	ValorContratos float64       `bson:"valor_contratos,omitempty"`
	NumContratos   int64         `bson:"num_contratos,omitempty"`
	Municipios     []*Municipio  `bson:"municipios,omitempty"`
	Partidos       []*Partido    `bson:"partidos,omitempty"`
}

type ResumoContratos struct {
	Quantidade int64   `bson:"quantidade,omitempty" json:"quantidade"`
	Valor      float64 `bson:"valor,omitempty" json:"valor"`
}

type Municipio struct {
	Cod             string          `bson:"cod,omitempty" json:"cod"`
	Nome            string          `bson:"nome,omitempty" json:"nome"`
	ResumoContratos ResumoContratos `bson:"resumo_contratos,omitempty" json:"resumo_contratos"`
	SiglaPartido    string          `bson:"sigla,omitempty" json:"sigla"`
}

type Partido struct {
	Sigla           string          `bson:"sigla,omitempty" json:"sigla"`
	ResumoContratos ResumoContratos `bson:"resumo_contratos,omitempty" json:"resumo_contratos"`
}
