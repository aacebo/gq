package gq

type Resolver func(params *ResolveParams) Result
type Middleware func(params *ResolveParams, next Resolver) Result
