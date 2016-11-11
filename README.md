# Contratos Públicos

> Resumo dos gastos que governos estaduais e municipais tem com fornecedores

[![Build Status][build-badge]][build-status] [![Coverage Status][cov-badge]][cov-status]

## Objetivo Geral

Maior transparência sobre os gastos estaduais e municipais através do
resumo e consolidação do acesso a informação.

## Objetivos Específicos

1. Prover API estável que poderá ser usada por outros sistemas e 
projetos de análise de dados;
1. Prover um portal que permita buscar realizar buscas por CNPJ e ter
ver um resumo do relacionamento daquela pessoa jurídica com os governos
estaduas;
1. Prover uma página estável (permalink) para cada empresa: esse link
pode ser facilmente compartilhado em redes sociais e outros meios.

## Tecnologias Utilizadas

Existe um esforço grande no projeto para utilização de ferramentas de
livres, gratuitas e/ou código fonte aberto. Dentre as quais, destacamos:

* Backend: Escrito utilizando a linguagem de programação [Go](https://golang.org)
* DNS/CDN: [Cloudflare](https://cloudflare.com)
* Frontend: HTML+CSS (bootstrap)+JS
* Integração contínua: [TravisCI][build-status]
* Cobertura de testes: [Coveralls.io][cov-status]
* PaaS: [Heroku](https://heroku.com)
* Banco de dados: [MongoDB](https://www.mongodb.com)
* Monitoramento de desempenho da aplicação (APM): [New Relic](https://newrelic.com/)
* Monitoramento de utilização: [Google Analytics](https://analytics.google.com)


[build-badge]:https://travis-ci.org/danielfireman/contratospublicos.svg?branch=master
[build-status]:https://travis-ci.org/danielfireman/contratospublicos
[cov-badge]:https://coveralls.io/repos/github/danielfireman/contratospublicos/badge.svg?branch=master
[cov-status]:https://coveralls.io/github/danielfireman/contratospublicos?branch=master
