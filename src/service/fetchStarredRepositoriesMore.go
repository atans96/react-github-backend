package service

import (
	"backend/src/service/stargazersQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
	"strings"
	"sync"
)

func (f *Fetch) FetchStarredRepositoriesMore(node types.StargazersNode, idx int, wg *sync.WaitGroup, userInfo UserInfo) {
	running := node.StarredRepositories.PageInfo.HasNextPage
	fetchMore := node.StarredRepositories.PageInfo.EndCursor
	f.FetchStarredRepositoriesStargazers(node.StarredRepositories.Nodes, idx, wg, userInfo)
	for running {
		req := graphql.NewRequest(stargazersQuery.StarredRepositoriesMore)
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
			if strings.Contains(strings.ToLower(err.Error()), "wrong") {
				fmt.Println("Retrying due to: ", err.Error())
				continue
			} else {
				panic(err)
			}
		}
		f.FetchStarredRepositoriesStargazers(respData.Repository.Stargazers.Nodes[idx].StarredRepositories.Nodes, idx, wg, userInfo)
		running = respData.Repository.Stargazers.Nodes[idx].StarredRepositories.PageInfo.HasNextPage
		if running {
			fetchMore = respData.Repository.Stargazers.Nodes[idx].StarredRepositories.PageInfo.EndCursor
		}
	}
	wg.Done()
}
