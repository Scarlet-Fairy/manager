package logger

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

type Loggers struct {
	InfoLogger  log.Logger
	WarnLogger  log.Logger
	ErrorLogger log.Logger

	TransportLayerLogger log.Logger
	EndpointLayerLogger  log.Logger
	CoreLayerLogger      log.Logger

	ServiceComponentLogger    log.Logger
	MessageComponentLogger    log.Logger
	RepositoryComponentLogger log.Logger
	SchedulerComponentLogger  log.Logger
}

func NewLogger() Loggers {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = level.NewInjector(logger, level.InfoValue())

	return Loggers{
		InfoLogger:  logger,
		WarnLogger:  level.Warn(logger),
		ErrorLogger: level.Error(logger),

		TransportLayerLogger: log.With(logger, "layer", "transport"),
		EndpointLayerLogger:  log.With(logger, "layer", "endpoint"),
		CoreLayerLogger:      log.With(logger, "layer", "core"),

		ServiceComponentLogger:    log.With(logger, "component", "service"),
		MessageComponentLogger:    log.With(logger, "component", "message"),
		RepositoryComponentLogger: log.With(logger, "component", "repository"),
		SchedulerComponentLogger:  log.With(logger, "component", "scheduler"),
	}
}
