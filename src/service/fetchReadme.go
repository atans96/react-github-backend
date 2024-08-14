package service

import (
	"backend/src/types"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/hashicorp/go-retryablehttp"
	stripmd "github.com/writeas/go-strip-markdown"
	"go.mongodb.org/mongo-driver/bson"
	mongoDB "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (f *Fetch) fetchReadme(repo types.RepoInfo, userInfo *UserInfo) {
	_, err := f.FetchRateLimit(userInfo.Token)
	if err != nil {
		panic(err)
	}
	makeNew, err := http.NewRequest("GET", "https://api.github.com/repos/"+repo.FullName+"/readme", nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Set("Authorization", "Bearer "+userInfo.Token)
	makeNew.Header.Set("accept", "application/vnd.github.VERSION.raw")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 50 * time.Second
	standardClient := retryClient.StandardClient() // *http.GQLClient
	response, err := standardClient.Do(makeNew)

	if err != nil {
		return
	}
	if (response.StatusCode < 200 || response.StatusCode >= 300) && response.StatusCode != 404 {
		fmt.Println("Negative status code: " + strconv.Itoa(response.StatusCode) + ". For url: " + response.Request.RequestURI)
		fmt.Println(response.Body)
		panic("failed")
	} else if response.StatusCode == 404 {
		language := repo.Language
		if len(language) == 0 {
			language = "No Language"
		}
		output := types.ReturnRepoInfo{
			FullName:      repo.FullName,
			Description:   repo.Description,
			Stars:         repo.StargazersCount,
			Forks:         repo.Forks,
			UpdatedAt:     repo.UpdatedAt,
			Language:      language,
			Topics:        repo.Topics,
			DefaultBranch: repo.DefaultBranch,
			HtmlUrl:       repo.HTMLURL,
			Readme:        "",
		}
		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": userInfo.UserName}
		data, err := bson.Marshal(output)
		if err != nil {
			return
		}

		err = bson.Unmarshal(data, &output)
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "repoInfo", Value: output},
				},
			},
			{Key: "$push",
				Value: bson.D{
					{Key: "languages", Value: language},
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
	} else {
		language := repo.Language
		if len(language) == 0 {
			language = "No Language"
		}
		bo, err := ioutil.ReadAll(response.Body)
		go response.Body.Close()
		if err != nil {
			return
		}
		output := types.ReturnRepoInfo{
			FullName:      repo.FullName,
			Description:   repo.Description,
			Stars:         repo.StargazersCount,
			Forks:         repo.Forks,
			UpdatedAt:     repo.UpdatedAt,
			Language:      repo.Language,
			Topics:        repo.Topics,
			DefaultBranch: repo.DefaultBranch,
			HtmlUrl:       repo.HTMLURL,
			Readme:        strip.StripTags(stripmd.Strip(strings.ToValidUTF8(string(bo), ""))),
		}
		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": userInfo.UserName}
		data, err := bson.Marshal(output)
		if err != nil {
			return
		}
		err = bson.Unmarshal(data, &output)
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "repoInfo", Value: output},
				},
			},
			{Key: "$push",
				Value: bson.D{
					{Key: "languages", Value: language},
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
}
