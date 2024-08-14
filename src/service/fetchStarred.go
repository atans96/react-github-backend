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
	"strconv"
	"time"
)

func (f *Fetch) fetchStarred(page int, userInfo *UserInfo) {
	_, err := f.FetchRateLimit(userInfo.Token)
	if err != nil {
		panic(err)
	}
	var repoInfo []types.RepoInfo
	var starred []types.Starred
	makeNew, err := http.NewRequest("GET", "https://api.github.com/users/"+userInfo.UserName+"/starred?per_page=100&page="+strconv.Itoa(page), nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Set("accept", "application/vnd.github.v3+json,application/vnd.github.mercy-preview+json,application/vnd.github.nebula-preview+json")
	makeNew.Header.Set("Authorization", "Bearer "+userInfo.Token)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 50 * time.Second
	standardClient := retryClient.StandardClient() // *http.GQLClient
	resp, err := standardClient.Do(makeNew)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 403 {
		f.fetchStarred(page, userInfo)
	} else if resp.StatusCode != 200 {
		panic(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	go resp.Body.Close()
	page += 1
	if len(result) > 2 {
		if err := json.Unmarshal(result, &repoInfo); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(result, &starred); err != nil {
			panic(err)
		}
		for _, value := range repoInfo {
			go f.fetchReadme(value, userInfo)
			go f.fetchContributors(value, userInfo)
			//go f.FetchFileName(value) //TODO:
		}
		f.fetchStarred(page, userInfo)
		if len(result) > 2 {
			go func(starred []types.Starred) {
				opts := options.FindOneAndUpdate().SetUpsert(true)
				query := bson.M{"userName": userInfo.UserName}
				update := bson.D{
					{Key: "$set",
						Value: bson.D{
							{Key: "starred", Value: starred},
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
			}(starred)
		}
	}
}
