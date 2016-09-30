package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/danielfireman/contratospublicos/receitaws"
)

type Fornecedor struct {
	ID                    string       `json:"id"`
	Nome                  string       `json:"nome"`
	Legislatura           string       `json:"legislatura"`
	ValorContratos        float64      `json:"valor_contratos"`
	NumContratos          int64        `json:"num_contratos"`
	AtividadePrincipal    []*receitaws.Atividade        `json:"atividade_principal"`
	Cnpj string `json:"cnpj"`
	Bairro string `json:"bairro"`
	Logradouro string `json:"logradouro"`
	Numero string `json:"numero"`
	CEP string `json:"cep"`
	Municipio string `json:"municipio"`
	UF string `json:"uf"`
	DataAbertura string `json:"abertura"`
	NaturezaJuridica string `json:"natureza_juridica"`
	NomeFantasia string `json:"nome_fantasia"`
	DataSituacao          string `json:"data_situacao"`
	Tipo                  string `json:"tipo"`
	AtividadesSecundarias []*receitaws.Atividade `json:"atividades_secundarias"`
	Situacao              string `json:"situacao"`
	NomeReceita           string `json:"nome_receita"`
	Telefone              string `json:"telefone"`
	UltimaAtualizacaoReceita string `json:"ultima_atualizacao_receita"`
	Municipios            []*Municipio `json:"municipios"`
	Partidos              []*Partido   `json:"partidos"`
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
