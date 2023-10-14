module github.com/ilya-mezentsev/micro-dep/store

go 1.21.3

replace github.com/ilya-mezentsev/micro-dep/shared => ../shared

require github.com/ilya-mezentsev/micro-dep/shared v0.0.0-00010101000000-000000000000

require (
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
