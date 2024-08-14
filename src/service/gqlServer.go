package service

import (
	"backend/src/service/gql"
	"github.com/graphql-go/graphql"
)

var repositoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			"stargazers": &graphql.Field{
				Type: gql.Stargazer,
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 4,
					},
				},
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			"repository": &graphql.Field{
				Type: repositoryType,
				Args: graphql.FieldConfigArgument{
					"owner": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
			},
		},
	})
