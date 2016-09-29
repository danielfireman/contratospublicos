package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/danielfireman/contratospublicos/model"
)

// A declaração de constantes abaixo tem que manter a mesma ordem do cabeçalho CSV.
const (
	CPFCNPJ = iota
	NOME_FORNECEDOR
	CODIGO_MUNICIPIO_FORNECEDOR
	COD_MUNICIPIO
	ANO_ELEICAO
	QT_EMPENHOS
	VL_EMPENHOS
	SIGLA_PARTIDO
)

type dadosFornecedor struct {
	fornecedor *model.Fornecedor
	municipios map[string]*model.Municipio
	partidos   map[string]*model.Partido
}

var fornecedores = make(map[string]*dadosFornecedor)

func main() {
	nLinhas := 0
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	for {
		linha, err := r.Read()
		if err == io.EOF {
			break
		}
		if nLinhas == 0 {
			nLinhas++
			continue
		}

		id := linha[CPFCNPJ]
		dados, ok := fornecedores[id]
		if !ok {
			dados = &dadosFornecedor{
				fornecedor: &model.Fornecedor{
					ID:   id,
					Nome: linha[NOME_FORNECEDOR],
				},
				municipios: make(map[string]*model.Municipio),
				partidos:   make(map[string]*model.Partido),
			}
			fornecedores[id] = dados
		}

		qtContratos, err := strconv.ParseInt(linha[QT_EMPENHOS], 10, 64)
		if err != nil {
			log.Fatal("Error processando a linha '%v': %q", linha, err)
		}

		valorContratos, err := strconv.ParseFloat(linha[VL_EMPENHOS], 64)
		if err != nil {
			log.Fatal("Error processando a linha '%v': %q", linha, err)
		}
		dados.fornecedor.NumTotalContratos += qtContratos
		dados.fornecedor.ValorTotalContratos += valorContratos

		codMunicipio := linha[COD_MUNICIPIO]
		siglaPartido := linha[SIGLA_PARTIDO]

		// Pegando informações sobre o município e resumindo.
		m, ok := dados.municipios[codMunicipio]
		if !ok {
			m = &model.Municipio{
				Cod:          codMunicipio,
				SiglaPartido: siglaPartido,
			}
			dados.municipios[codMunicipio] = m
		}
		m.ResumoContratos.Quantidade += qtContratos
		m.ResumoContratos.Valor += valorContratos

		// Pegando informações sobre o partido e resumindo.
		p, ok := dados.partidos[siglaPartido]
		if !ok {
			p = &model.Partido{
				Sigla: siglaPartido,
			}
			dados.partidos[siglaPartido] = p
		}
		p.ResumoContratos.Quantidade += qtContratos
		p.ResumoContratos.Valor += valorContratos
		nLinhas++
	}

	for _, d := range fornecedores {
		for _, m := range d.municipios {
			d.fornecedor.Municipios = append(d.fornecedor.Municipios, m)
		}
		for _, p := range d.partidos {
			d.fornecedor.Partidos = append(d.fornecedor.Partidos, p)
		}
	}
	fmt.Print(nLinhas)
}
