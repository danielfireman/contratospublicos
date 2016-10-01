package receitaws

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

	ret := &DadosReceitaWS{}
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	if ret.Status != "OK" {
		return nil, fmt.Errorf("Error calling receitaws: '%s'", ret.Message)
	}
	return ret, nil
}
