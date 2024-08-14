package service

import (
	"backend/src/service/stargazersQuery"
	"backend/src/types"
	"context"
	"github.com/machinebox/graphql"
	"os"
	"strings"
	"sync"
)

func (f *Fetch) FetchStarredRepositoriesStargazers(StarredRepositoriesNodes []types.StarredRepositoriesNodes, idx int, wg *sync.WaitGroup, userInfo UserInfo) {
	for ix, starredRepoNode := range StarredRepositoriesNodes {
		running := starredRepoNode.Stargazers.PageInfo.HasNextPage
		fetchMore := starredRepoNode.Stargazers.PageInfo.EndCursor
		f.FetchUserQuery(userInfo, starredRepoNode.Stargazers.Nodes)
		for running {
			req := graphql.NewRequest(stargazersQuery.StarredRepositoriesStargazers)
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
			var respData types.GraphqlQueryResponse
			_, err := f.FetchRateLimit(os.Getenv("TOKEN_C"))
			if err != nil {
				panic(err)
			}
			if err := GQLClient.Run(ctx, req, &respData); err != nil {
				panic(err)
			}
			go f.FetchUserQuery(userInfo, respData.Repository.Stargazers.Nodes[idx].StarredRepositories.Nodes[ix].Stargazers.Nodes)
			running = respData.Repository.Stargazers.Nodes[idx].StarredRepositories.Nodes[ix].Stargazers.PageInfo.HasNextPage
			if running {
				fetchMore = respData.Repository.Stargazers.Nodes[idx].StarredRepositories.Nodes[ix].Stargazers.PageInfo.EndCursor
			}
		}
	}
	wg.Done()
}
