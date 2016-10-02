package model

type Fornecedor struct {
	ID                       string                     `json:"id"`
	Nome                     string                     `json:"nome"`
	Legislatura              string                     `json:"legislatura"`
	AtividadePrincipal       []*Atividade               `json:"atividade_principal"`
	Cnpj                     string                     `json:"cnpj"`
	Bairro                   string                     `json:"bairro"`
	Logradouro               string                     `json:"logradouro"`
	Numero                   string                     `json:"numero"`
	CEP                      string                     `json:"cep"`
	Municipio                string                     `json:"municipio"`
	UF                       string                     `json:"uf"`
	DataAbertura             string                     `json:"abertura"`
	NaturezaJuridica         string                     `json:"natureza_juridica"`
	NomeFantasia             string                     `json:"nome_fantasia"`
	DataSituacao             string                     `json:"data_situacao"`
	Tipo                     string                     `json:"tipo"`
	AtividadesSecundarias    []*Atividade               `json:"atividades_secundarias"`
	Situacao                 string                     `json:"situacao"`
	NomeReceita              string                     `json:"nome_receita"`
	Telefone                 string                     `json:"telefone"`
	UltimaAtualizacaoReceita string                     `json:"ultima_atualizacao_receita"`
	ResumoContratos          *ResumoContratosFornecedor `json:"resumo_contratos"`
}

type Atividade struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

type ResumoContratosFornecedor struct {
	ValorContratos float64      `json:"valor_contratos"`
	NumContratos   int64        `json:"num_contratos"`
	Municipios     []*Municipio `json:"municipios"`
	Partidos       []*Partido   `json:"partidos"`
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
