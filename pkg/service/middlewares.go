package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(service Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(service Service) Service {
		return &loggingMiddlware{
			next:   service,
			logger: logger,
		}
	}
}

type loggingMiddlware struct {
	next   Service
	logger log.Logger
}

func (l *loggingMiddlware) Deploy(ctx context.Context, gitRepo string, name string, envs map[string]string) (deploy *Deploy, err error) {
	defer func() {
		l.logger.Log(
			"method", "Deploy",
			"gitRepo", gitRepo,
			"name", name,
			"envs", envs,
			"deploy", deploy,
			"err", err,
		)
	}()

	return l.next.Deploy(ctx, gitRepo, name, envs)
}

func (l *loggingMiddlware) Destroy(ctx context.Context, deployId string) (err error) {
	defer func() {
		l.logger.Log(
			"method", "Destroy",
			"deployId", deployId,
			"err", err,
		)
	}()

	return l.next.Destroy(ctx, deployId)
}
