module github.com/ilya-mezentsev/micro-dep/store

go 1.21.3

replace github.com/ilya-mezentsev/micro-dep/shared => ../shared

require (
	github.com/ilya-mezentsev/micro-dep/shared v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.3.5
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
