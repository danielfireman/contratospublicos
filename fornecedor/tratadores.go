package fornecedor

import (
	"log"
	"net/http"
	"os"
	"strings"
	"html/template"
	"strconv"
	"time"

	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"
	"github.com/leekchan/accounting"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
)

const (
	DB = "heroku_q6gnv76m"
)

type T struct {
	buscador *Buscador
	ac       accounting.Accounting
}

func Tratadores() (*T, error) {
	s, err := store.MongoDB(os.Getenv("MONGODB_URI"), DB)
	if err != nil {
		return nil, err
	}
	return &T{
		buscador: &Buscador{
			principal: ColetorBD(s),
			secundarios: []ColetorDadosFornecedor{
				ColetorReceitaWs(),
				ColetorResumoContratos(s),
			},
		},
		ac: accounting.Accounting{
			Symbol:    "",
			Precision: 2,
			Format:    "%s %v",
			Decimal:   ",",
			Thousand:  ".",
		},
	}, err
}

func (t *T) TrataAPICall() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id, legislatura := extraiParametros(r)
		if id == "" {
			http.Error(w, "CNPJ inválido.", http.StatusBadRequest)
			return
		}
		resultado, err := t.buscador.ColetaDados(id, legislatura)
		if err != nil {
			if NaoEncontrado(err) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Println("Err id:'%s' err:'%q'", id, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&resultado); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type fornecedorVO struct {
	Nome                string
	NomeFantasia        string
	Cnpj                string
	EnderecoParte1      string
	EnderecoParte2      string
	UltimaAtualizacao   string
	Legislatura         string
	Telefone            string
	Tipo                string
	DataAbertura        string
	EstaAtiva           bool
	Email               string
	AtividadePrimaria   []*model.Atividade
	AtividadeSecundaria []*model.Atividade
	ResumoContratos     *resumoContratosFornecedorVO
}

type resumoContratosFornecedorVO struct {
	ValorContratos string
	NumContratos   string
	Municipios     []*municipioVO
	Partidos       []*partidoVO
}

type municipioVO struct {
	Nome            string
	ResumoContratos resumoContratosVO `bson:"resumo_contratos,omitempty" json:"resumo_contratos,omitempty"`
	SiglaPartido    string
}

type resumoContratosVO struct {
	Quantidade string
	Valor      string
}

type partidoVO struct {
	Sigla           string
	ResumoContratos resumoContratosVO
}

var fornecedorTmpl = template.Must(template.New("fornecedor").ParseFiles("fornecedor/fornecedor.html"))

func (t *T) TrataPaginaFornecedor() func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id, legislatura := extraiParametros(r)
		if id == "" {
			http.Error(w, "CNPJ inválido.", http.StatusBadRequest)
			return
		}
		f, err := t.buscador.ColetaDados(id, legislatura)
		if err != nil {
			if NaoEncontrado(err) {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			log.Println("Err id:'%s' err:'%q'", id, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		fVO := fornecedorVO{}
		fVO.Cnpj = f.Cnpj
		fVO.Nome = f.Nome
		fVO.NomeFantasia = f.NomeFantasia
		switch f.Legislatura {
		case "2008":
			fVO.Legislatura = "2008-2012"
		case "2012":
			fVO.Legislatura = "2012-2016"
		}
		if f.UltimaAtualizacaoReceita != "" {
			t, _ := time.Parse(time.RFC3339, f.UltimaAtualizacaoReceita)
			fVO.UltimaAtualizacao = strconv.Itoa(int(time.Now().Sub(t).Hours() / float64(24)))
		}
		fVO.EnderecoParte1 = f.Logradouro
		if f.Numero == "" {
			fVO.EnderecoParte1 += ", S/N"
		} else {
			fVO.EnderecoParte1 += ", " + f.Numero
		}
		fVO.EnderecoParte1 += ", " + f.Bairro
		fVO.EnderecoParte2 += f.Municipio + "-" + f.UF + ", " + f.CEP
		fVO.Telefone = f.Telefone
		fVO.Tipo = f.Tipo
		fVO.DataAbertura = f.DataAbertura
		fVO.EstaAtiva = f.Situacao == "ATIVA"
		fVO.AtividadePrimaria = f.AtividadesPrincipais
		fVO.AtividadeSecundaria = f.AtividadesSecundarias

		// Formatando os resumos dos contratos.
		fVO.ResumoContratos = &resumoContratosFornecedorVO{
			ValorContratos: t.ac.FormatMoney(f.ResumoContratos.ValorContratos),
			NumContratos:   strconv.Itoa(int(f.ResumoContratos.NumContratos)),
		}
		for _, m := range f.ResumoContratos.Municipios {
			fVO.ResumoContratos.Municipios = append(fVO.ResumoContratos.Municipios, &municipioVO{
				Nome:         m.Nome,
				SiglaPartido: m.SiglaPartido,
				ResumoContratos: resumoContratosVO{
					Valor:      t.ac.FormatMoney(m.ResumoContratos.Valor),
					Quantidade: strconv.Itoa(int(m.ResumoContratos.Quantidade)),
				},
			})
		}
		for _, m := range f.ResumoContratos.Partidos {
			fVO.ResumoContratos.Partidos = append(fVO.ResumoContratos.Partidos, &partidoVO{
				Sigla: m.Sigla,
				ResumoContratos: resumoContratosVO{
					Valor:      t.ac.FormatMoney(m.ResumoContratos.Valor),
					Quantidade: strconv.Itoa(int(m.ResumoContratos.Quantidade)),
				},
			})
		}
		if err := fornecedorTmpl.Execute(w, &fVO); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func extraiParametros(r *http.Request) (string, string) {
	query := r.URL.Query()
	legislatura := query.Get("legislatura")
	if legislatura == "" {
		legislatura = "2012"
	}

	// NOTA: Utilizando parametros de consulta para permitir que usuários copiem e colem CNPJs
	// completos.
	id := query.Get("cnpj")
	if id == "" {
		return "", legislatura
	}

	// Removendo caracteres especiais que existem no CPF e CNPJ.
	// Isso permite que os usuários copiem e colem CPFs e CNPJs de sites na internet e outras
	// fontes.
	id = strings.Replace(id, ".", "", -1)
	id = strings.Replace(id, "-", "", -1)
	id = strings.Replace(id, "/", "", -1)
	return id, legislatura
}
