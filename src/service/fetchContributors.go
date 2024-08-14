package service

import (
	"backend/src/types"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"time"
)

func (f *Fetch) fetchContributors(repo types.RepoInfo, userInfo *UserInfo) {
	_, err := f.FetchRateLimit(userInfo.Token)
	if err != nil {
		panic(err)
	}
	makeNew, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo.FullName+"/contributors?q=contributions&order=desc", nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Set("Authorization", "Bearer "+userInfo.Token)
	makeNew.Header.Set("accept", "application/vnd.github.v3+json,application/vnd.github.mercy-preview+json,application/vnd.github.nebula-preview+json")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 50 * time.Second
	standardClient := retryClient.StandardClient() // *http.GQLClient
	response, err := standardClient.Do(makeNew)
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	go response.Body.Close()
	var contributorResponsesServer []types.ContributorResponse
	var contributorOutput types.ReturnContributors
	err = json.Unmarshal(result, &contributorResponsesServer)
	if err != nil {
		panic(err)
	}
	contributorOutput.FullName = repo.FullName
	contributorOutput.Contributors = &contributorResponsesServer
	opts := options.FindOneAndUpdate().SetUpsert(true)
	query := bson.M{"userName": userInfo.UserName}
	data, err := bson.Marshal(contributorOutput)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &contributorOutput)
	update := bson.D{
		{Key: "$addToSet",
			Value: bson.D{
				{Key: "repoContributions", Value: contributorOutput},
			},
		},
	}
	err = f.Mongo.DB.FindOneAndUpdate(nil, query, update, opts).Err()
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongoDB.ErrNoDocuments {
			fmt.Println(err)
		}
		panic(err)
	}
}
