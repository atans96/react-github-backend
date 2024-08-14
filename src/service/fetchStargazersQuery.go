package service

import (
	"backend/src/service/stargazersQuery"
	"backend/src/types"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"sort"
	"strings"
	"sync"
)

func internalizeReaction(reaction string) int {
	switch reaction {
	case "THUMBS_UP":
		return 1
	case "THUMBS_DOWN":
		return -1
	case "LAUGH":
		return 0
	case "HOORAY":
		return 1
	case "CONFUSED":
		return 0
	case "HEART":
		return 1
	case "ROCKET":
		return 1
	case "EYES":
		return 0
	default:
		panic("unreachable reactions emoji")
	}
}

type DiscussionScore struct {
	Score     int      `json:"score" bson:"score"`
	FullName  string   `json:"full_name" bson:"full_name"`
	Languages []string `json:"languages" bson:"languages"`
}
type ContributionByRepo struct {
	FullName  string   `json:"full_name" bson:"full_name"`
	Languages []string `json:"languages" bson:"languages"`
}
type ContributionCollectionScore struct {
	Score int `json:"score" bson:"score"`
}

func CalculateScoreReactions(Reactions []types.Reactions, owner string) int {
	var total int
	for _, reaction := range Reactions {
		if reaction.User.Login != owner {
			total += internalizeReaction(reaction.Content)
		}
	}
	return total
}
func CalculateScoreContributionCollection(C types.ContributionsCollection) int {
	var total int
	total += C.TotalIssueContributions + C.TotalPullRequestContributions + C.TotalPullRequestReviewContributions
	return total
}
func ProcessLanguages(edges []types.LanguagesEdges) []string {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Size > edges[j].Size
	})
	var languages []string
	for _, edge := range edges {
		languages = append(languages, edge.Node.Name)
	}
	return languages
}

type CustomError struct {
	Message string
}

func (f *Fetch) fetchStargazersQuery(userInfo UserInfo) {
	running := true
	query := stargazersQuery.QuerySuggestedFetchUsers
	fetchMore := ""
	var WG sync.WaitGroup
	for running {
		WG.Add(1)
		req := graphql.NewRequest(query)
		req.Var("owner", strings.Split(userInfo.FullName, "/")[0])
		req.Var("name", strings.Split(userInfo.FullName, "/")[1])
		if len(fetchMore) > 0 {
			req.Var("after", fetchMore)
		}
		// set header fields
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Authorization", "Bearer "+userInfo.Token)

		// define a Context for the request
		ctx := context.TODO()

		// run it and capture the response
		var respData types.GraphqlQueryResponse
		_, err := f.FetchRateLimit(userInfo.Token)
		if err != nil {
			panic(err)
		}
		if err := GQLClient.Run(ctx, req, &respData); err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "wrong") {
				fmt.Println("Retrying due to: ", err.Error())
				WG.Done()
				continue
			} else {
				panic(err)
			}
		}
		if len(respData.Repository.Stargazers.Nodes) > 0 {
			var wg sync.WaitGroup
			for idx, node := range respData.Repository.Stargazers.Nodes {
				wg.Add(1) //need to be based on how many goroutine running in this loop
				//go f.FetchRepositoryDiscussionReactionsMore(node, idx, &wg, respData, userInfo)
				//go f.FetchRepositoryDiscussionMore(node, idx, &wg, userInfo)
				//
				//go f.FetchRepositoryDiscussionCommentsReactionsMore(node, idx, &wg, respData, userInfo)
				//go f.FetchRepositoryDiscussionComments(node, idx, &wg, userInfo)
				//
				//go f.FetchStarredRepositoriesStargazers(node.StarredRepositories.Nodes, idx, &wg, userInfo)
				//go f.FetchStarredRepositoriesMore(node, idx, &wg, userInfo)
				if len(node.ContributionsCollection.PopularIssueContribution.Issue.Repository.Owner.Id) > 0 {
					go f.FetchPopularIssueContributionParticipantsMore(node, idx, &wg, userInfo)
				} else {
					wg.Done()
				}
				//if len(node.ContributionsCollection.PopularPullRequestContribution.PullRequest.Repository.Name) > 0 {
				//	go f.FetchPopularPullRequestContributionParticipantsMore(node, idx, &wg, userInfo)
				//} else {
				//	wg.Done()
				//}
			}
			wg.Wait()
			WG.Done()
		} else {
			WG.Done()
		}
		WG.Wait()
		running = respData.Repository.Stargazers.PageInfo.HasNextPage
		if running {
			query = stargazersQuery.QuerySuggestedFetchUsersMore
			fetchMore = respData.Repository.Stargazers.PageInfo.EndCursor
		}
	}
	fmt.Println("done")
}
