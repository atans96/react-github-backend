package gql

import "github.com/graphql-go/graphql"

var starredNode = graphql.NewObject(graphql.ObjectConfig{
	Fields: graphql.Fields{
		"languages": &graphql.Field{
			Type: languages,
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type:         graphql.Int,
					DefaultValue: 5,
				},
			},
		},
	},
})
var starredRepository = graphql.NewObject(graphql.ObjectConfig{
	Fields: graphql.Fields{
		"pageInfo": &graphql.Field{
			Type: pageInfo,
		},
		"nodes": &graphql.Field{
			Type: starredNode,
		},
	},
})
