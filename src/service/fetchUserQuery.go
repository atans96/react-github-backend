package service

import (
	userQuery "backend/src/service/userQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"strings"
)

func (f *Fetch) FetchUserQuery(userInfo UserInfo, StargazersNode []types.StargazersNode) {
	for _, node := range StargazersNode {
		query := userQuery.QueryUser
		req := graphql.NewRequest(query)
		// set any variables
		req.Var("login", node.Login)
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+userInfo.Token)

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryUserResponse
		_, err := f.FetchRateLimit(userInfo.Token)
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

		go f.UserFetchRepositoryDiscussionReactionsMore(node, respData, userInfo)
		go f.UserFetchRepositoryDiscussionMore(node, userInfo)

		go f.UserFetchRepositoryDiscussionCommentsReactionsMore(node, respData, userInfo)
		go f.UserFetchRepositoryDiscussionComments(node, userInfo)

		go f.UserFetchStarredRepositoriesStargazers(node.StarredRepositories.Nodes, userInfo)
		go f.UserFetchStarredRepositoriesMore(node, userInfo)
		if len(node.ContributionsCollection.PopularIssueContribution.Issue.Repository.Owner.Id) > 0 {
			go f.UserFetchPopularIssueContributionParticipantsMore(node, userInfo)
		}

		if len(node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository.Owner.Id) > 0 {
			go f.UserFetchPopularPullRequestContributionParticipantsMore(node, userInfo)
		}
	}
}
