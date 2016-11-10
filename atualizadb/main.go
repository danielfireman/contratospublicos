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

	"gopkg.in/mgo.v2"
	"github.com/danielfireman/contratospublicos/supplier"
)

var dbURI = flag.String("dburi", "", "URI completa do mongo db que vai ser populado.")

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

// TODO(danielfireman): Translate eveyrthing to English.
type dadosFornecedor struct {
	resumo     map[string]*resumoFornecedor
}

type resumoFornecedor struct {
	id string
	valor      float64
	num        int32
	municipios map[string]*supplier.City
	partidos   map[string]*supplier.Party
}

// Mapeia fornecedores referentes a uma determinada legislatura
var dadosFornecedores = make(map[string]*dadosFornecedor)

func main() {
	flag.Parse()

	if *dbURI == "" {
		log.Fatalf("--dburi flag é obrigatória.")
	}
	mgoInfo, err := mgo.ParseURL(*dbURI)
	if err != nil {
		log.Fatalf("Erro processando URI:%s err:%q\n", *dbURI, err)
	}

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
				resumo: make(map[string]*resumoFornecedor),
			}
			dadosFornecedores[id] = dados
		}

		resumo, ok := dados.resumo[legislatura]
		if !ok {
			resumo = &resumoFornecedor{
				id: id,
				municipios: make(map[string]*supplier.City),
				partidos:   make(map[string]*supplier.Party),
			}
			dados.resumo[legislatura] = resumo
		}

		qtContratos, err := strconv.ParseInt(linha[QT_EMPENHOS], 10, 32)
		if err != nil {
			log.Fatalf("Error processando a linha '%v': %q", linha, err)
		}

		valorContratos, err := strconv.ParseFloat(linha[VL_EMPENHOS], 64)
		if err != nil {
			log.Fatalf("Error processando a linha '%v': %q", linha, err)
		}
		resumo.num += int32(qtContratos)
		resumo.valor += valorContratos

		codMunicipio := linha[COD_MUNICIPIO]
		siglaPartido := linha[SIGLA_PARTIDO]

		m, ok := resumo.municipios[codMunicipio]
		if !ok {
			m = &supplier.City{
				ID:          codMunicipio,
				PartyInitials: siglaPartido,
			}
			resumo.municipios[codMunicipio] = m
		}
		m.NumContracts += int32(qtContratos)
		m.AmountCountracts += valorContratos

		// Pegando informações sobre o partido e resumindo.
		p, ok := resumo.partidos[siglaPartido]
		if !ok {
			p = &supplier.Party{
				Initials: siglaPartido,
			}
			resumo.partidos[siglaPartido] = p
		}
		p.NumContracts += int32(qtContratos)
		p.AmountCountracts += valorContratos
		nLinhas++
	}

	fmt.Println("Processed ", nLinhas, " rows")
	resumos := make(map[string][]interface{})
	for _, d := range dadosFornecedores {
		for l, r := range d.resumo {
			resumo := &supplier.ContractsSummary{
				ID:             r.id,
				AmountContracts: r.valor,
				NumContracts:   r.num,
			}
			for _, m := range r.municipios {
				resumo.Cities = append(resumo.Cities, m)
			}
			for _, p := range r.partidos {
				resumo.Parties = append(resumo.Parties, p)
			}
			resumos[l] = append(resumos[l], resumo)
		}
	}

	session, err := mgo.DialWithInfo(mgoInfo)
	if err != nil {
		log.Fatal(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	for l, r := range resumos {
		c := session.DB(mgoInfo.Database).C(l)
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
