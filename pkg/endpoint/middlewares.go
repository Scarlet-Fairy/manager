package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"time"
)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())

			return next(ctx, request)
		}
	}
}

func UnwrapErrorMiddleware() endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			resp, err := e(ctx, request)
			if err != nil {
				return nil, err
			}

			if f, ok := resp.(endpoint.Failer); ok {
				if f.Failed() != nil {
					return nil, f.Failed()
				}
			}

			return resp, nil
		}
	}
}
