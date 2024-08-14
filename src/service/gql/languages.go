package gql

import "github.com/graphql-go/graphql"

var languages = graphql.NewObject(graphql.ObjectConfig{
	Fields: graphql.Fields{
		"edges": &graphql.Field{Type: graphql.NewObject(graphql.ObjectConfig{
			Fields: graphql.Fields{
				"node": &graphql.Field{Type: graphql.NewObject(graphql.ObjectConfig{
					Fields: graphql.Fields{
						"name": &graphql.Field{
							Type: graphql.String,
						},
					},
				})},
				"size": &graphql.Field{Type: graphql.Int},
			},
		})},
	},
})
