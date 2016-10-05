package main

import (
	"log"
	"os"

	"net/http"

	"github.com/danielfireman/contratospublicos/fornecedor"
	"github.com/danielfireman/contratospublicos/templates"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/yvasiyarov/gorelic"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variável de ambiente $PORT obrigatória.")
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
	agent.CollectHTTPStatuses = true
	agent.CollectMemoryStat = true
	agent.NewrelicPollInterval = 120
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("Monitoramento NewRelic configurado com sucesso.")

	// Configuração do roteador echo.
	e := echo.New()
	e.Static("/", "public/index.html")
	e.Static("/", "public")
	e.SetRenderer(templates.T)
	e.GET("/repo", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "http://github.com/danielfireman/contratospublicos")
	})

	// Lida com fornecedores
	fTratadores, err := fornecedor.Tratadores()
	if err != nil {
		log.Fatal(err)
	}
	e.GET("/api/v1/fornecedor", fTratadores.TrataAPICall())
	e.GET("/fornecedor", fTratadores.TrataPaginaFornecedor())

	log.Println("Serviço inicializado na porta ", port)
	log.Fatal(e.Run(fasthttp.New(":" + port)))
}
