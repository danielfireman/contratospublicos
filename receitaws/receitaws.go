package receitaws

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/danielfireman/contratospublicos/model"
)

type DadosReceitaWS struct {
	AtividadePrincipal    []*Atividade `json:"atividade_principal"`
	DataSituacao          string       `json:"data_situacao"`
	Tipo                  string       `json:"tipo"`
	Nome                  string       `json:"nome"`
	Telefone              string       `json:"telefone"`
	AtividadesSecundarias []*Atividade `json:"atividades_secundarias"`
	Situacao              string       `json:"situacao"`
	Cnpj                  string       `json:"cnpj"`
	Bairro                string       `json:"bairro"`
	Logradouro            string       `json:"logradouro"`
	Numero                string       `json:"numero"`
	CEP                   string       `json:"cep"`
	Municipio             string       `json:"municipio"`
	UF                    string       `json:"uf"`
	DataAbertura          string       `json:"abertura"`
	NaturezaJuridica      string       `json:"natureza_juridica"`
	NomeFantasia          string       `json:"fantasia"`
	UltimaAtualizacao     string       `json:"ultima_atualizacao"`

	// Error
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Atividade struct {
	Text string `json:"text"`
	Code string `json:"code"`
}

const (
	timeout = 500 * time.Millisecond
	url     = "http://receitaws.com.br/v1/cnpj/"
)

var client = http.DefaultClient

func Fetch(ctx context.Context, id string) (interface{}, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest("GET", url+id, nil)
	req.WithContext(timeoutCtx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// TODO(danielfireman): Usar sync.Pool
	ret := &DadosReceitaWS{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	if ret.Status != "OK" {
		return nil, fmt.Errorf("Error calling receitaws: '%s'", ret.Message)
	}
	return ret, nil
}

func AtualizaFornecedor(f *model.Fornecedor, i interface{}) {
	dr := i.(*DadosReceitaWS)
	f.DataSituacao = dr.DataSituacao
	f.Tipo = dr.Tipo
	f.Situacao = dr.Situacao
	f.NomeReceita = dr.Nome
	f.Telefone = dr.Telefone
	f.Cnpj = dr.Cnpj
	f.Municipio = dr.Municipio
	f.UF = dr.UF
	f.DataAbertura = dr.DataAbertura
	f.NaturezaJuridica = dr.NaturezaJuridica
	f.NomeFantasia = dr.NomeFantasia
	f.UltimaAtualizacaoReceita = dr.UltimaAtualizacao
	f.Bairro = dr.Bairro
	f.Logradouro = dr.Logradouro
	f.Numero = dr.CEP
	f.CEP = dr.CEP
	for _, a := range dr.AtividadePrincipal {
		f.AtividadePrincipal = append(f.AtividadePrincipal, &model.Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
	for _, a := range dr.AtividadesSecundarias {
		f.AtividadesSecundarias = append(f.AtividadesSecundarias, &model.Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
}
