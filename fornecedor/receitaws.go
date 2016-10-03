package fornecedor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
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

type coletorReceitaWS struct {
	cliente *http.Client
	pool    sync.Pool
}

func ColetorReceitaWs() ColetorDadosFornecedor {
	return &coletorReceitaWS{
		cliente: http.DefaultClient,
		pool: sync.Pool{
			New: func() interface{} {
				return &DadosReceitaWS{}
			},
		},
	}
}

func (c *coletorReceitaWS) ColetaDados(ctx context.Context, fornecedor *model.Fornecedor) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest("GET", url+fornecedor.ID, nil)
	req.WithContext(timeoutCtx)
	resp, err := c.cliente.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	dr := c.pool.Get().(*DadosReceitaWS)
	defer c.pool.Put(dr)
	defer func(dr *DadosReceitaWS) {
		p := reflect.ValueOf(dr).Elem()
		p.Set(reflect.Zero(p.Type()))
	}(dr)

	if err := json.Unmarshal(body, dr); err != nil {
		return err
	}
	if dr.Status != "OK" {
		return fmt.Errorf("Error calling receitaws: '%s'", dr.Message)
	}
	fornecedor.DataSituacao = dr.DataSituacao
	fornecedor.Tipo = dr.Tipo
	fornecedor.Situacao = dr.Situacao
	fornecedor.NomeReceita = dr.Nome
	fornecedor.Telefone = dr.Telefone
	fornecedor.Cnpj = dr.Cnpj
	fornecedor.Municipio = dr.Municipio
	fornecedor.UF = dr.UF
	fornecedor.DataAbertura = dr.DataAbertura
	fornecedor.NaturezaJuridica = dr.NaturezaJuridica
	fornecedor.NomeFantasia = dr.NomeFantasia
	fornecedor.UltimaAtualizacaoReceita = dr.UltimaAtualizacao
	fornecedor.Bairro = dr.Bairro
	fornecedor.Logradouro = dr.Logradouro
	fornecedor.Numero = dr.CEP
	fornecedor.CEP = dr.CEP
	for _, a := range dr.AtividadePrincipal {
		fornecedor.AtividadesPrincipais = append(fornecedor.AtividadesPrincipais, &model.Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
	for _, a := range dr.AtividadesSecundarias {
		fornecedor.AtividadesSecundarias = append(fornecedor.AtividadesSecundarias, &model.Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
	return nil
}
