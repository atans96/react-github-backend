package service

import (
	userQuery "backend/src/service/userQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
	"strings"
)

func (f *Fetch) UserFetchStarredRepositoriesMore(node types.StargazersNode, userInfo UserInfo) {
	running := node.StarredRepositories.PageInfo.HasNextPage
	fetchMore := node.StarredRepositories.PageInfo.EndCursor
	f.UserFetchStarredRepositoriesStargazers(node.StarredRepositories.Nodes, userInfo)
	for running {
		req := graphql.NewRequest(userQuery.StarredRepositoriesMore)
		// set any variables
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_C"))

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryUserResponse
		_, err := f.FetchRateLimit(os.Getenv("TOKEN_C"))
		if err != nil {
			panic(err)
		}
		if err := GQLClient.Run(ctx, req, &respData); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "wrong") {
				fmt.Println("Retrying due to: ", err.Error())
				continue
			} else {
				panic(err)
			}
		}
		f.UserFetchStarredRepositoriesStargazers(respData.User.StarredRepositories.Nodes, userInfo)
		running = respData.User.StarredRepositories.PageInfo.HasNextPage
		if running {
			fetchMore = respData.User.StarredRepositories.PageInfo.EndCursor
		}
	}
}
