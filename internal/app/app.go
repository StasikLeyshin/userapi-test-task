package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
)

type App struct {
	logger     *logrus.Logger
	components []component
}

type component interface {
	Start() error
	Stop(ctx context.Context) error
}

func NewApp(logger *logrus.Logger, components ...component) *App {
	return &App{
		logger:     logger,
		components: components,
	}
}

func (a *App) Run(ctx context.Context) {
	componentsCtx, componentsStopCtx := signal.NotifyContext(ctx, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer componentsStopCtx()

	for _, comp := range a.components {
		err := comp.Start()
		if err != nil {
			a.logger.Printf("error when starting the component %v", err)
		}
	}

	<-componentsCtx.Done()

	for _, comp := range a.components {
		err := comp.Stop(componentsCtx)
		if err != nil {
			a.logger.Printf("error when stopping the component %v", err)
		}
	}

}
