package gql

import "github.com/graphql-go/graphql"

var reactions = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type: graphql.Int,
			},
			"nodes": &graphql.Field{
				Type: graphql.NewObject(
					graphql.ObjectConfig{
						Fields: graphql.Fields{
							"content": &graphql.Field{
								Type: graphql.String,
							},
						},
					},
				),
			},
		},
	},
)
