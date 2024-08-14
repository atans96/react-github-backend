package gql

import "github.com/graphql-go/graphql"

var pageInfo = graphql.NewObject(graphql.ObjectConfig{
	Fields: graphql.Fields{
		"endCursor":   &graphql.Field{Type: graphql.String},
		"hasNextPage": &graphql.Field{Type: graphql.Boolean},
		"startCursor": &graphql.Field{Type: graphql.String},
	},
})
