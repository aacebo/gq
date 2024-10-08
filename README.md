<p align="center">
	<img src="./assets/icon.png" width="120px" style="border-radius:20%" />
</p>
 
<p align="center">
	a zero dependency performant graph query resolver
</p>

<p align="center">
	<a href="https://opensource.org/licenses/MIT" target="_blank" alt="License">
		<img src="https://img.shields.io/badge/License-MIT-blue.svg" />
	</a>
	<a href="https://pkg.go.dev/github.com/aacebo/gq" target="_blank" alt="Go Reference">
		<img src="https://pkg.go.dev/badge/github.com/aacebo/gq.svg" />
	</a>
	<a href="https://goreportcard.com/report/github.com/aacebo/gq" target="_blank" alt="Go Report Card">
		<img src="https://goreportcard.com/badge/github.com/aacebo/gq" />
	</a>
	<a href="https://github.com/aacebo/owl/actions/workflows/ci.yml" target="_blank" alt="Build">
		<img src="https://github.com/aacebo/owl/actions/workflows/ci.yml/badge.svg?branch=main" />
	</a>
	<a href="https://codecov.io/gh/aacebo/gq" > 
		<img src="https://codecov.io/gh/aacebo/gq/graph/badge.svg?token=9XETRUUQUY" /> 
	</a>
</p>

# Install

```bash
go get github.com/aacebo/gq
```

# Usage

```go
schema := owl.String().Required()

if err := schema.Validate("..."); err != nil { // nil
	panic(err)
}
```

