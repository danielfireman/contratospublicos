package main

import (
	"context"
	"strconv"
	"time"

	"github.com/danielfireman/contratospublicos/supplier"
	"github.com/leekchan/accounting"
)

var ac = accounting.Accounting{
	Symbol:    "",
	Precision: 2,
	Format:    "%s %v",
	Decimal:   ",",
	Thousand:  ".",
}

// fetchSupplierVirtualObject fetches the data holder of data needed to render the supplier page tamplte.
func fetchSupplierVirtualObject(ctx context.Context, fetcher *supplier.DataFecher, id, legislature string) (*fornecedorVO, error) {
	f, err := fetcher.Summary(ctx, id, legislature)
	if err != nil {
		return nil, err
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

	if f.Logradouro != "" {
		fVO.EnderecoParte1 = f.Logradouro
		if f.Numero == "" {
			fVO.EnderecoParte1 += ", S/N"
		} else {
			fVO.EnderecoParte1 += ", " + f.Numero
		}
		fVO.EnderecoParte1 += ", " + f.Bairro
		fVO.EnderecoParte2 += f.Municipio + "-" + f.UF + ", " + f.CEP
	}
	fVO.Telefone = f.Telefone
	fVO.Tipo = f.Tipo
	fVO.DataAbertura = f.DataAbertura
	fVO.EstaAtiva = f.Situacao == "ATIVA"
	for _, a := range f.AtividadesPrincipais {
		fVO.AtividadePrimaria = append(fVO.AtividadePrimaria, &atividade{
			Code: a.Code,
			Text: a.Text,
		})
	}
	for _, a := range f.AtividadesSecundarias {
		fVO.AtividadeSecundaria = append(fVO.AtividadePrimaria, &atividade{
			Code: a.Code,
			Text: a.Text,
		})
	}

	// Formatando os resumos dos contratos.
	fVO.ResumoContratos = &resumoContratosFornecedorVO{
		ValorContratos: ac.FormatMoney(f.ResumoContratos.ValorContratos),
		NumContratos:   strconv.Itoa(int(f.ResumoContratos.NumContratos)),
	}
	for _, m := range f.ResumoContratos.Municipios {
		fVO.ResumoContratos.Municipios = append(fVO.ResumoContratos.Municipios, &municipioVO{
			Nome:         m.Nome,
			SiglaPartido: m.SiglaPartido,
			ResumoContratos: resumoContratosVO{
				Valor:      ac.FormatMoney(m.ResumoContratos.Valor),
				Quantidade: strconv.Itoa(int(m.ResumoContratos.Quantidade)),
			},
		})
	}
	for _, m := range f.ResumoContratos.Partidos {
		fVO.ResumoContratos.Partidos = append(fVO.ResumoContratos.Partidos, &partidoVO{
			Sigla: m.Sigla,
			ResumoContratos: resumoContratosVO{
				Valor:      ac.FormatMoney(m.ResumoContratos.Valor),
				Quantidade: strconv.Itoa(int(m.ResumoContratos.Quantidade)),
			},
		})
	}
	return &fVO, nil
}

// ## Template data structs ##

type atividade struct {
	Text string `json:"text"`
	Code string `json:"code"`
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
	AtividadePrimaria   []*atividade
	AtividadeSecundaria []*atividade
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
