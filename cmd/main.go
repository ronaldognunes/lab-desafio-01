package main

import (
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-chi/chi/v5"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/cep"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/temperatura"
	web "github.com/ronaldognunes/lab-desafio-01/internal/infra/webserver"
)

func main() {

	servicoCep := cep.NewCepService("https://viacep.com.br/ws/")
	servicoTemperatura := temperatura.NewTemperaturaService("http://api.weatherapi.com/v1/current.json?q=", "6c71bade554742f0b46142657250703")
	consultaHandler := web.NewConsultaHandler(servicoCep, servicoTemperatura)

	router := chi.NewRouter()
	router.Get("/consulta-temperatura-por-cep", consultaHandler.ConsultarCepHandler)

	http.ListenAndServe(":8080", router)
}
