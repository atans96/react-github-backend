package gql

import "github.com/graphql-go/graphql"

var repositoryDiscussionCommentsNodes = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"upvoteCount": &graphql.Field{
				Type: graphql.Int,
			},
			"discussion": &graphql.Field{
				Type: graphql.NewObject(
					graphql.ObjectConfig{
						Fields: graphql.Fields{
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
				),
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
		},
	},
)

var repositoryDiscussionComments = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"pageInfo": &graphql.Field{
				Type: pageInfo,
			},
			"nodes": &graphql.Field{
				Type: repositoryDiscussionCommentsNodes,
			},
		},
	},
)
