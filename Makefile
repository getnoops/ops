install:
	npm i -g get-graphql-schema

graphql:
	get-graphql-schema http://localhost:8200/graphql > pkg/queries/schema.graphql

generate:
	go generate ./...

build:
	go build -o ops main.go
	mv ops ~/go/bin/