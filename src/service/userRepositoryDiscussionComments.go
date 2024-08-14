package service

import (
	userQuery "backend/src/service/userQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

func (f *Fetch) UserFetchRepositoryDiscussionComments(node types.StargazersNode, userInfo UserInfo) {
	running := node.RepositoryDiscussionComments.PageInfo.HasNextPage
	fetchMore := node.RepositoryDiscussionComments.PageInfo.EndCursor
	var discussionScore DiscussionScore
	var score int
	var upvoteScore int
	for ix, discussionCommentsNode := range node.RepositoryDiscussionComments.Nodes {
		if node.RepositoryDiscussionComments.Nodes[ix].ViewerHasUpvoted {
			upvoteScore = node.RepositoryDiscussionComments.Nodes[ix].UpvoteCount - 1
		} else {
			if node.RepositoryDiscussionComments.Nodes[ix].UpvoteCount > 1 {
				upvoteScore = node.RepositoryDiscussionComments.Nodes[ix].UpvoteCount
			}
		}
		temp := CalculateScoreReactions(node.RepositoryDiscussionComments.Nodes[ix].Reactions.Nodes, strings.Split(userInfo.FullName, "/")[0])
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
		req := graphql.NewRequest(userQuery.RepositoryDiscussionComments)
		// set any variables
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN_B"))

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryUserResponse
		_, err := f.FetchRateLimit(os.Getenv("TOKEN_B"))
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
		f.UserFetchRepositoryDiscussionCommentsReactionsMore(respData.User, respData, userInfo)
		running = respData.User.RepositoryDiscussionComments.PageInfo.HasNextPage
		if running {
			fetchMore = respData.User.RepositoryDiscussionComments.PageInfo.EndCursor
		}
	}
}
