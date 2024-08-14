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

func (f *Fetch) FetchRepositoryDiscussionMore(node types.StargazersNode, idx int, wg *sync.WaitGroup, userInfo UserInfo) {
	running := node.RepositoryDiscussions.PageInfo.HasNextPage
	fetchMore := node.RepositoryDiscussions.PageInfo.EndCursor
	var discussionScore DiscussionScore
	var score int
	var upvoteScore int
	for ix, discussionCommentsNode := range node.RepositoryDiscussions.Nodes {
		if node.RepositoryDiscussions.Nodes[ix].ViewerHasUpvoted {
			upvoteScore = node.RepositoryDiscussions.Nodes[ix].UpvoteCount - 1
		} else {
			if node.RepositoryDiscussions.Nodes[ix].UpvoteCount > 1 {
				upvoteScore = node.RepositoryDiscussions.Nodes[ix].UpvoteCount
			}
		}
		temp := CalculateScoreReactions(node.RepositoryDiscussions.Nodes[ix].Reactions.Nodes, strings.Split(userInfo.FullName, "/")[0])
		score += temp + upvoteScore
		discussionScore.Languages = ProcessLanguages(discussionCommentsNode.Repository.Languages.Edges)
		discussionScore.Score = score
		discussionScore.FullName = userInfo.FullName

		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": node.Login}
		data, err := bson.Marshal(discussionScore)
		if err != nil {
			panic(err)
		}

		err = bson.Unmarshal(data, &discussionScore)
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "discussionCommentsScore", Value: discussionScore},
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
	}

	for running {
		wg.Add(1)
		req := graphql.NewRequest(stargazersQuery.RepositoryDiscussionsMore)
		// set any variables
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_A"))

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryResponse
		_, err := f.FetchRateLimit(os.Getenv("TOKEN_A"))
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
		f.FetchRepositoryDiscussionReactionsMore(respData.Repository.Stargazers.Nodes[idx], idx, wg, respData, userInfo)
		running = respData.Repository.Stargazers.Nodes[idx].RepositoryDiscussions.PageInfo.HasNextPage
		if running {
			fetchMore = respData.Repository.Stargazers.Nodes[idx].RepositoryDiscussions.PageInfo.EndCursor
		}
	}
	wg.Done()
}
