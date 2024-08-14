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

func (f *Fetch) FetchPopularIssueContributionParticipantsMore(node types.StargazersNode, idx int, wg *sync.WaitGroup, userInfo UserInfo) {
	running := node.ContributionsCollection.PopularIssueContribution.Issue.Participants.PageInfo.HasNextPage
	fetchMore := node.ContributionsCollection.PopularIssueContribution.Issue.Participants.PageInfo.EndCursor
	go f.FetchFileName(node.ContributionsCollection.PopularIssueContribution.Issue.Repository, userInfo)
	//go f.FetchUserQuery(userInfo, node.ContributionsCollection.PopularIssueContribution.Issue.Participants.Nodes)
	for running {
		req := graphql.NewRequest(stargazersQuery.PopularIssueContributionParticipantsMore)
		// set any variables
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_D"))

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryResponse
		_, err := f.FetchRateLimit(os.Getenv("TOKEN_D"))
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
		go f.FetchFileName(respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularIssueContribution.Issue.Repository, userInfo)
		//go f.FetchUserQuery(userInfo, respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularIssueContribution.Issue.Participants.Nodes)
		running = respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularIssueContribution.Issue.Participants.PageInfo.HasNextPage
		if running {
			fetchMore = respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularIssueContribution.Issue.Participants.PageInfo.EndCursor
		}
	}
	wg.Done()
}
