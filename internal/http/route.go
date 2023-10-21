package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"refactoring/internal/service"
	"time"
)

type Config struct {
	Port int `yaml:"port"`
}

type HttpRouter struct {
	port   int
	server *http.Server
	client *service.Service
	logger *log.Logger
}

func NewHttpRouter(config Config, client *service.Service, logger *log.Logger) *HttpRouter {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	httpRouter := HttpRouter{
		port: config.Port,
		server: &http.Server{
			Handler: r,
			Addr:    fmt.Sprintf(":%d", config.Port),
		},
		client: client,
		logger: logger,
	}

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", httpRouter.searchUsers)
				r.Post("/", httpRouter.createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", httpRouter.getUser)
					r.Patch("/", httpRouter.updateUser)
					r.Delete("/", httpRouter.deleteUser)
				})
			})
		})
	})

	return &httpRouter
}

func (h *HttpRouter) Start() error {
	go func() {
		h.logger.Printf("start http-server")
		err := h.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("fail to serve the server on the port %d", h.port)
		}
	}()
	return nil
}

func (h *HttpRouter) Stop(ctx context.Context) error {
	fmt.Println("http router is stopping")
	return h.server.Shutdown(ctx)
}
