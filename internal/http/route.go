package http

import (
	"fmt"
	"net/http"
	"refactoring/internal/service"
)

type Config struct {
	Port int `yaml:"port"`
}

type HttpRouter struct {
	port   int
	server *http.ServeMux
	client *service.Service
}

func NewHttpRouter(config Config, client *service.Service) *HttpRouter {
	serveMux := http.NewServeMux()
	fs := http.FileServer(http.Dir("web"))
	serveMux.Handle("/", fs)
	httpRouter := HttpRouter{
		port:   config.Port,
		server: serveMux,
		client: client,
	}
	//serveMux.HandleFunc("/order", GetOrder(&httpRouter))
	return &httpRouter
}

func (r *HttpRouter) Start() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.server)
	if err != nil {
		return err
	}
	return nil
}
