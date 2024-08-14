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

func (f *Fetch) UserFetchPopularPullRequestContributionParticipantsMore(node types.StargazersNode, userInfo UserInfo) {
	running := node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.HasNextPage
	fetchMore := node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.EndCursor
	go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestContributionsByRepository, "pullRequestByRepo", userInfo)
	go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestReviewContributionsByRepository, "pullRequestReviewByRepo", userInfo)
	go f.FetchFileName(node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository, userInfo)
	go f.FetchUserQuery(userInfo, node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.Nodes)

	for running {
		req := graphql.NewRequest(userQuery.PopularPullRequestContributionParticipantsMore)
		// set any variables
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_E"))

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryUserResponse
		_, err := f.FetchRateLimit(os.Getenv("TOKEN_E"))
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
		go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestContributionsByRepository, "pullRequestByRepo", userInfo)
		go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestReviewContributionsByRepository, "pullRequestReviewByRepo", userInfo)
		go f.FetchFileName(respData.User.ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository, userInfo)
		go f.FetchUserQuery(userInfo, respData.User.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.Nodes)
		running = respData.User.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.HasNextPage
		if running {
			fetchMore = respData.User.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.EndCursor
		}
	}
}
