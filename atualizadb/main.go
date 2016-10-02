package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/danielfireman/contratospublicos/model"
	"gopkg.in/mgo.v2"
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

const (
	DB = "heroku_q6gnv76m"
)

type dadosFornecedor struct {
	fornecedor *model.DadosFornecedor
	resumo     map[string]*resumoFornecedor
}

type resumoFornecedor struct {
	valor      float64
	num        int64
	municipios map[string]*model.Municipio
	partidos   map[string]*model.Partido
}

// Mapeia fornecedores referentes a uma determinada legislatura
var dadosFornecedores = make(map[string]*dadosFornecedor)

func main() {
	dbURI := os.Getenv("MONGODB_URI")
	if dbURI == "" {
		log.Fatalf("Variável de ambiente MONGHQ_URL obrigatória.")
	}

	flag.Parse()

	nLinhas := 0
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	for {
		linha, err := r.Read()
		if err == io.EOF {
			break
		}
		// Ignorando cabeçalho.
		if nLinhas == 0 {
			nLinhas++
			continue
		}

		legislatura := linha[ANO_ELEICAO]
		id := linha[CPFCNPJ]
		dados, ok := dadosFornecedores[id]
		if !ok {
			dados = &dadosFornecedor{
				fornecedor: &model.DadosFornecedor{
					ID:   id,
					Nome: linha[NOME_FORNECEDOR],
				},
				resumo: make(map[string]*resumoFornecedor),
			}
			dadosFornecedores[id] = dados
		}

		resumo, ok := dados.resumo[legislatura]
		if !ok {
			resumo = &resumoFornecedor{
				municipios: make(map[string]*model.Municipio),
				partidos:   make(map[string]*model.Partido),
			}
			dados.resumo[legislatura] = resumo
		}

		qtContratos, err := strconv.ParseInt(linha[QT_EMPENHOS], 10, 64)
		if err != nil {
			log.Fatalf("Error processando a linha '%v': %q", linha, err)
		}

		valorContratos, err := strconv.ParseFloat(linha[VL_EMPENHOS], 64)
		if err != nil {
			log.Fatalf("Error processando a linha '%v': %q", linha, err)
		}
		resumo.num += qtContratos
		resumo.valor += valorContratos

		codMunicipio := linha[COD_MUNICIPIO]
		siglaPartido := linha[SIGLA_PARTIDO]

		m, ok := resumo.municipios[codMunicipio]
		if !ok {
			m = &model.Municipio{
				Cod:          codMunicipio,
				SiglaPartido: siglaPartido,
				//Nome: nomeMunicipio,
			}
			resumo.municipios[codMunicipio] = m
		}
		m.ResumoContratos.Quantidade += qtContratos
		m.ResumoContratos.Valor += valorContratos

		// Pegando informações sobre o partido e resumindo.
		p, ok := resumo.partidos[siglaPartido]
		if !ok {
			p = &model.Partido{
				Sigla: siglaPartido,
			}
			resumo.partidos[siglaPartido] = p
		}
		p.ResumoContratos.Quantidade += qtContratos
		p.ResumoContratos.Valor += valorContratos
		nLinhas++
	}

	session, err := mgo.Dial(dbURI)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.SetMode(mgo.Eventual, true)

	fornecedores := make([]interface{}, 0, len(dadosFornecedores))
	resumos := make(map[string][]interface{})
	for _, d := range dadosFornecedores {
		fornecedores = append(fornecedores, d.fornecedor)
		for l, r := range d.resumo {
			resumo := &model.ResumoContratosFornecedor{
				ID:             d.fornecedor.ID,
				ValorContratos: r.valor,
				NumContratos:   r.num,
			}
			for _, m := range r.municipios {
				resumo.Municipios = append(resumo.Municipios, m)
			}
			for _, p := range r.partidos {
				resumo.Partidos = append(resumo.Partidos, p)
			}
			resumos[l] = append(resumos[l], resumo)
		}
	}

	// Inserindo fornecedores
	fmt.Printf("Inserindo %d fornecedores.\n", len(fornecedores))
	c := session.DB(DB).C("fornecedores")
	fornecedoresIndex := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(fornecedoresIndex)
	if err != nil {
		log.Fatalf("Erro criando índice: %q", err)
	}
	if err := c.Insert(fornecedores...); err != nil {
		log.Fatalf("Erro inserindo fornecedores: %q", err)
	}
	fmt.Printf("%d fornecedores inseridos com sucesso.\n", len(fornecedores))

	// Inserindo Resumos
	for l, r := range resumos {
		c := session.DB(DB).C(l)
		resumoIndex := mgo.Index{
			Key:        []string{"id"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}
		err = c.EnsureIndex(resumoIndex)
		if err != nil {
			log.Fatalf("Erro criando índice: %q", err)
		}
		fmt.Printf("[Legislatura %s] Inserindo %d resumos.\n", l, len(resumos[l]))
		if err := c.Insert(r...); err != nil {
			log.Fatalf("[Legislatura %s] Erro inserindo resumos: %q", l, err)
		}
		fmt.Printf("[Legislatura %s] %d resumos inseridos com sucesso.\n", l, len(resumos[l]))
	}
}
