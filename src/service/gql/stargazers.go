package gql

import "github.com/graphql-go/graphql"

var Stargazer = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"pageInfo": &graphql.Field{
				Type: pageInfo,
			},
		},
	},
)
