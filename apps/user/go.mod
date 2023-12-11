module github.com/ilya-mezentsev/micro-dep/user

go 1.21.3

replace github.com/ilya-mezentsev/micro-dep/shared => ../shared

require (
	github.com/frankenbeanies/uuid4 v0.0.0-20180313125435-68b799ec299a
	github.com/ilya-mezentsev/micro-dep/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
