package supplier

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	timeout = 1000 * time.Millisecond
	url     = "http://receitaws.com.br/v1/cnpj/"
)

// PopulateSupplierInfo populates the passed-in supplier with data that comes from receitaws web service.
func FetchReceitaWSData(ctx context.Context, errChan chan error, id string, supplier *Fornecedor) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequest("GET", url+id, nil)
	req.WithContext(timeoutCtx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errChan <- err
		return
	}

	dr := &DadosReceitaWS{}
	if err := json.Unmarshal(body, dr); err != nil {
		errChan <- err
		return
	}
	if dr.Status != "OK" {
		if strings.Contains(dr.Message, "CNPJ inválido") {
			errChan <- NotFoundErr
		} else {
			errChan <- err
		}
		return
	}
	supplier.DataSituacao = dr.DataSituacao
	supplier.Nome = dr.Nome
	supplier.Tipo = dr.Tipo
	supplier.Situacao = dr.Situacao
	supplier.NomeReceita = dr.Nome
	supplier.Telefone = dr.Telefone
	supplier.Cnpj = dr.Cnpj
	supplier.Municipio = dr.Municipio
	supplier.UF = dr.UF
	supplier.DataAbertura = dr.DataAbertura
	supplier.NaturezaJuridica = dr.NaturezaJuridica
	supplier.NomeFantasia = dr.NomeFantasia
	supplier.Bairro = dr.Bairro
	supplier.Logradouro = dr.Logradouro
	supplier.Numero = dr.Numero
	supplier.CEP = dr.CEP
	supplier.UltimaAtualizacaoReceita = dr.UltimaAtualizacao
	for _, a := range dr.AtividadePrincipal {
		supplier.AtividadesPrincipais = append(supplier.AtividadesPrincipais, &Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
	for _, a := range dr.AtividadesSecundarias {
		supplier.AtividadesSecundarias = append(supplier.AtividadesSecundarias, &Atividade{
			Text: a.Text,
			Code: a.Code,
		})
	}
}

// ## RECEITAWS DATA MODEL ##

// DadosReceitaWS holds supplier data returned from receitaws web service.
type DadosReceitaWS struct {
	AtividadePrincipal    []*AtividadeReceitaWS `json:"atividade_principal"`
	DataSituacao          string                `json:"data_situacao"`
	Tipo                  string                `json:"tipo"`
	Nome                  string                `json:"nome"`
	Telefone              string                `json:"telefone"`
	AtividadesSecundarias []*AtividadeReceitaWS `json:"atividades_secundarias"`
	Situacao              string                `json:"situacao"`
	Cnpj                  string                `json:"cnpj"`
	Bairro                string                `json:"bairro"`
	Logradouro            string                `json:"logradouro"`
	Numero                string                `json:"numero"`
	CEP                   string                `json:"cep"`
	Municipio             string                `json:"municipio"`
	UF                    string                `json:"uf"`
	DataAbertura          string                `json:"abertura"`
	NaturezaJuridica      string                `json:"natureza_juridica"`
	NomeFantasia          string                `json:"fantasia"`
	UltimaAtualizacao     string                `json:"ultima_atualizacao"`

	// Error
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AtividadeReceitaWS represents an supplier activity area, returned from receitaws web service.
type AtividadeReceitaWS struct {
	Text string `json:"text"`
	Code string `json:"code"`
}
