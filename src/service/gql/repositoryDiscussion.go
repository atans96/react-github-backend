package gql

import "github.com/graphql-go/graphql"

var repositoryDiscussions = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"pageInfo": &graphql.Field{
				Type: pageInfo,
			},
			"nodes": &graphql.Field{
				Type: repositoryDiscussionsNodes,
			},
		},
	},
)
var repositoryDiscussionsNodes = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"upvoteCount": &graphql.Field{
				Type: graphql.Int,
			},
			"reactions": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 5,
					},
				},
				Type: reactions,
			},
			"repository": &graphql.Field{
				Type: repository,
			},
			"category": &graphql.Field{
				Type: graphql.NewObject(
					graphql.ObjectConfig{
						Fields: graphql.Fields{
							"name": &graphql.Field{
								Type: graphql.String,
							},
						},
					},
				),
			},
		},
	},
)
var repository = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
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
	},
)
