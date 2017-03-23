package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/model"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	model  *model.Model
	router *mux.Router
}

func New(model *model.Model) *Api {
	c := controller{model: model}
	a := Api{
		router: router(&c),
	}
	return &a
}

func (a *Api) ListenAndServe() {
	err := http.ListenAndServe(":8080", a.router)
	if err != nil {
		fmt.Println(err)
	}
}
