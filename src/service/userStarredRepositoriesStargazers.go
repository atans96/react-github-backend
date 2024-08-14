package service

import (
	userQuery "backend/src/service/userQuery"
	"backend/src/types"
	"context"
	"github.com/machinebox/graphql"
	"os"
	"strings"
)

func (f *Fetch) UserFetchStarredRepositoriesStargazers(StarredRepositoriesNodes []types.StarredRepositoriesNodes, userInfo UserInfo) {
	for ix, starredRepoNode := range StarredRepositoriesNodes {
		running := starredRepoNode.Stargazers.PageInfo.HasNextPage
		fetchMore := starredRepoNode.Stargazers.PageInfo.EndCursor
		for running {
			req := graphql.NewRequest(userQuery.StarredRepositoriesStargazers)
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
				panic(err)
			}
			go f.FetchUserQuery(userInfo, respData.User.StarredRepositories.Nodes[ix].Stargazers.Nodes)
			running = respData.User.StarredRepositories.Nodes[ix].Stargazers.PageInfo.HasNextPage
			if running {
				fetchMore = respData.User.StarredRepositories.Nodes[ix].Stargazers.PageInfo.EndCursor
			}
		}
	}
}
