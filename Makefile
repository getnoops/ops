install:
	npm i -g get-graphql-schema

graphql:
	get-graphql-schema http://localhost:8200/api/graphql > pkg/queries/schema.graphql

generate:
	go generate ./...