package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/danielfireman/contratospublicos/fornecedor"
	"github.com/danielfireman/contratospublicos/model"
	"github.com/danielfireman/contratospublicos/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/yvasiyarov/gorelic"
)

const (
	DB = "heroku_q6gnv76m"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
	}

	mongoDBStore, err := store.MongoDB(os.Getenv("MONGODB_URI"), DB)
	if err != nil {
		log.Fatal(err)
	}

	nrLicence := os.Getenv("NEW_RELIC_LICENSE_KEY")
	if nrLicence == "" {
		log.Fatal("$NEW_RELIC_LICENSE_KEY must be set")
	}
	agent := gorelic.NewAgent()
	agent.Verbose = true
	agent.NewrelicLicense = nrLicence
	agent.NewrelicName = "contratospublicos"
	agent.CollectHTTPStat = true
	agent.CollectMemoryStat = true
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Monitoramento NewRelic configurado com sucesso.")

	// Configuração do roteador echo.
	e := echo.New()
	e.Static("/", "public/index.html")
	e.Static("/", "public")

	coletorBDPrincipal := fornecedor.ColetorBD(mongoDBStore)
	coletorReceitaWS := fornecedor.ColetorReceitaWs()
	coletorResumoContratos := fornecedor.ColetorResumoContratos(mongoDBStore)
	e.GET("/api/v1/fornecedor/:id", func(c echo.Context) error {
		ctx := context.Background()
		id := c.Param("id")

		if id == "" || id == "CNPJ" {
			return c.String(http.StatusBadRequest, "CNPJ do fornecedor obrigatório.")
		}

		legislatura := c.QueryParam("legislatura")
		if legislatura == "" {
			legislatura = "2012"
		}
		resultado := &model.Fornecedor{
			ID:          id,
			Legislatura: legislatura,
		}

		// Usando nosso BD como fonte autoritativa para buscas. Se não existe lá, nós
		// não conhecemos. Por isso, essa chamada é síncrona.
		tF := agent.Tracer.BeginTrace("fornecedores_collection_query")
		if err := coletorBDPrincipal.ColetaDados(ctx, resultado); err != nil {
			if fornecedor.NaoEncontrado(err) {
				tF.EndTrace()
				return c.NoContent(http.StatusNotFound)
			}
			log.Println("Err id:'%s' err:'%q'", id, err)
			tF.EndTrace()
			return c.NoContent(http.StatusInternalServerError)
		}

		// Pegando dados remotos de forma concorrente.
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(res *model.Fornecedor) {
			defer wg.Done()
			agent.Tracer.Trace("receitaws_query", func() {
				if coletorReceitaWS.ColetaDados(ctx, res); err != nil {
					log.Println("Err id:'%s' err:'%q'", id, err)
				}
			})
		}(resultado)
		wg.Add(1)
		go func(res *model.Fornecedor) {
			defer wg.Done()
			agent.Tracer.Trace(legislatura+"_collection_query", func() {
				if err := coletorResumoContratos.ColetaDados(ctx, res); err != nil {
					log.Println("Err id:'%s' err:'%q'", id, err)
				}
			})
		}(resultado)
		wg.Wait()
		return c.JSON(http.StatusOK, resultado)
	})
	log.Println("Serviço inicializado na porta ", port)
	log.Fatal(e.Run(fasthttp.New(":" + port)))
}
