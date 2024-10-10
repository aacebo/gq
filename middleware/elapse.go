package middleware

import (
	"time"

	"github.com/aacebo/gq"
)

func Elapse(params *gq.ResolveParams, next gq.Resolver) gq.Result {
	now := time.Now()
	res := next(params)
	res.Meta["$elapse"] = time.Now().Sub(now).Milliseconds()
	return res
}
