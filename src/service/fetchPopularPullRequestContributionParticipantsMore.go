package service

import (
	"backend/src/service/stargazersQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
	"sync"
)

func (f *Fetch) FetchPopularPullRequestContributionParticipantsMore(node types.StargazersNode, idx int, wg *sync.WaitGroup, userInfo UserInfo) {
	running := node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.HasNextPage
	fetchMore := node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.EndCursor
	go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestContributionsByRepository, "pullRequestByRepo", userInfo)
	go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestReviewContributionsByRepository, "pullRequestReviewByRepo", userInfo)
	go f.FetchFileName(node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository, userInfo)
	go f.FetchUserQuery(userInfo, node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.Nodes)

	for running {
		req := graphql.NewRequest(stargazersQuery.PopularPullRequestContributionParticipantsMore)
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
		var respData types.GraphqlQueryResponse
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
		var contributionCollectionScore ContributionCollectionScore
		contributionCollectionScore.Score = CalculateScoreContributionCollection(node.ContributionsCollection)

		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": node.Login}
		data, err := bson.Marshal(contributionCollectionScore)
		if err != nil {
			panic(err)
		}

		err = bson.Unmarshal(data, &contributionCollectionScore)
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "contributionCollectionScore", Value: contributionCollectionScore},
					{Key: "linkage", Value: userInfo.UserName},
				},
			},
			{Key: "$set",
				Value: bson.D{
					{Key: "bio", Value: node.Bio},
					{Key: "id", Value: node.Id},
					{Key: "location", Value: node.Location},
					{Key: "isGithubStar", Value: node.IsGithubStar},
				},
			},
		}
		err = f.Mongo.DBSuggested.FindOneAndUpdate(nil, query, update, opts).Err()
		if err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == mongoDB.ErrNoDocuments {
				fmt.Println(err)
			}
			fmt.Println(err)
		}
		go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestContributionsByRepository, "pullRequestByRepo", userInfo)
		go f.processContributionByRepository(node, node.ContributionsCollection.PullRequestReviewContributionsByRepository, "pullRequestReviewByRepo", userInfo)
		go f.FetchFileName(respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository, userInfo)
		go f.FetchUserQuery(userInfo, respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.Nodes)

		running = respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.HasNextPage
		if running {
			fetchMore = respData.Repository.Stargazers.Nodes[idx].ContributionsCollection.PopularPullRequestContribution.PullRequest.Participants.PageInfo.EndCursor
		}
	}
	wg.Done()
}
