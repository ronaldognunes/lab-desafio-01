package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ronaldognunes/lab-desafio-01/internal/entity"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/cep"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/temperatura"
)

type ResponseDto struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type ConsultaHandler struct {
	CepService         cep.CepService
	TemperaturaService temperatura.TemperaturaService
}

func NewConsultaHandler(cepService cep.CepService, temperaturaService temperatura.TemperaturaService) *ConsultaHandler {
	return &ConsultaHandler{CepService: cepService, TemperaturaService: temperaturaService}
}

func (c *ConsultaHandler) ConsultarCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP não informado", http.StatusBadRequest)
		return
	}

	if cepValido := entity.ValidaCEP(cep); cepValido == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	fmt.Println("Consultando CEP")
	dadosCep, err := c.CepService.ConsultarCep(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	if dadosCep.Cep == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	fmt.Println("Localidade:", dadosCep.Localidade)
	temp, err := c.TemperaturaService.ConsultarTemperatura(entity.RemoveAcentos(dadosCep.Localidade))

	if err != nil {
		http.Error(w, "Erro ao consultar temperatura", http.StatusBadRequest)
		return
	}

	if temp.TempC == 0 {
		http.Error(w, "Temperatura não encontrada", http.StatusBadRequest)
		return
	}

	fmt.Println("Temperatura:", temp)
	temp.CalcularTemperaturas()
	response := ResponseDto{
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: temp.TempK,
	}
	fmt.Println("Temperatura", response)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return

}
