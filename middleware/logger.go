package middleware

import (
	"log"
	"time"

	"github.com/aacebo/gq"
)

func Logger(logger *log.Logger) gq.Middleware {
	if logger == nil {
		logger = log.Default()
	}

	return func(params *gq.ResolveParams, next gq.Resolver) gq.Result {
		requestId := ""

		if params.Context != nil {
			v, exists := params.Context.Value("X-Request-Id").(string)

			if exists {
				requestId = v
			}
		}

		now := time.Now()
		res := next(params)

		if requestId != "" {
			logger.Printf("%s [%s] => %dms", requestId, params.Key, time.Now().Sub(now).Milliseconds())
		} else {
			logger.Printf("[%s] => %dms", params.Key, time.Now().Sub(now).Milliseconds())
		}

		return res
	}
}
