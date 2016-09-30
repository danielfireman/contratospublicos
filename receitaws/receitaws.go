package receitaws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func GetData(id string, ret *DadosReceitaWS) error {
	resp, err := http.Get("http://receitaws.com.br/v1/cnpj/" + id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, ret); err != nil {
		return err
	}
	if ret.Status != "OK" {
		return fmt.Errorf("Error calling receitaws: '%s'", ret.Message)
	}
	return nil
}
