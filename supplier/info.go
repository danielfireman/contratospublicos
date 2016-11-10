package supplier

// TODO(danielfireman): Translate to English
type Fornecedor struct {
	ID                       string                     `json:"id,omitempty"`
	Nome                     string                     `json:"nome,omitempty"`
	Legislatura              string                     `json:"legislatura,omitempty"`
	Cnpj                     string                     `json:"cnpj,omitempty"`
	Bairro                   string                     `json:"bairro,omitempty"`
	Logradouro               string                     `json:"logradouro,omitempty"`
	Numero                   string                     `json:"numero,omitempty"`
	CEP                      string                     `json:"cep,omitempty"`
	Municipio                string                     `json:"municipio,omitempty"`
	UF                       string                     `json:"uf,omitempty"`
	DataAbertura             string                     `json:"abertura,omitempty"`
	NaturezaJuridica         string                     `json:"natureza_juridica,omitempty"`
	NomeFantasia             string                     `json:"nome_fantasia,omitempty"`
	DataSituacao             string                     `json:"data_situacao,omitempty"`
	Tipo                     string                     `json:"tipo,omitempty"`
	Situacao                 string                     `json:"situacao,omitempty"`
	NomeReceita              string                     `json:"nome_receita,omitempty"`
	Telefone                 string                     `json:"telefone,omitempty"`
	UltimaAtualizacaoReceita string                     `json:"ultima_atualizacao_receita,omitempty"`
	AtividadesPrincipais     []*Atividade               `json:"atividades_principais,omitempty"`
	AtividadesSecundarias    []*Atividade               `json:"atividades_secundarias,omitempty"`
	ResumoContratos          *ResumoContratosFornecedor `json:"resumo_contratos,omitempty"`
}

type Atividade struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

type ResumoContratosFornecedor struct {
	ValorContratos float64      `json:"valor,omitempty"`
	NumContratos   int32        `json:"quantidade,omitempty"`
	Municipios     []*Municipio `json:"municipios,omitempty"`
	Partidos       []*Partido   `json:"partidos,omitempty"`
}

// TODO(danielfireman): Remove this (make it flat like in ContractsSummary
type ResumoContratos struct {
	Quantidade int32   `json:"quantidade,omitempty"`
	Valor      float64 `json:"valor,omitempty"`
}

type Municipio struct {
	Cod             string          `json:"cod,omitempty"`
	Nome            string          `json:"nome,omitempty"`
	ResumoContratos ResumoContratos `json:"resumo_contratos,omitempty"`
	SiglaPartido    string          `json:"sigla,omitempty"`
}

type Partido struct {
	Sigla           string          `json:"sigla,omitempty"`
	ResumoContratos ResumoContratos `json:"resumo_contratos,omitempty"`
}
